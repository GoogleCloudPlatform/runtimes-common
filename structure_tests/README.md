GCP Structure Tests
====================

The GCP Structure Tests provide a powerful framework to validate the structure of a Docker image. These tests can be used to check the output of commands in an image, as well as verify metadata and contents of the filesystem.

To run the test framework, simply download the binary for your OS here:
- TODO: add links to binaries

## Example Run
An example run of the test framework:
```shell
./structure-test -test.v -image gcr.io/google-appengine/python python_test_config.yaml
```
This command will run the tests on the GAE Python image, with verbose logging, using the python_test_config.yaml test config.

Tests within this framework are specified through a YAML or JSON config file, which is provided to the test driver as the last positional argument of the command. Multiple config files may be specified in a single test run. The config file will be loaded in by the test driver, which will execute the tests in order. Within this config file, three distinct types of tests can be written:

- Command Tests (testing output/error of a specific command issued)
- File Existence Tests (making sure a file is, or isn't, present in the file system of the image)
- File Content Tests (making sure files in the file system of the image contain, or do not contain, specific contents)

## Command Tests
Command tests ensure that certain commands run properly in the target image. Regexes can be used to check for expected or excluded strings in both stdout and stderr. Additionally, any number of flags can be passed to the argument as normal.

#### Supported Fields:

This is the current schema version (v2.0.0).

- Name (string, **required**): The name of the test
- Setup ([][]string, *optional*): A list of commands (each with optional flags) to run before the actual command under test.
- Teardown ([][]string, *optional*): A list of commands (each with optional flags) to run after the actual command under test.
- Command (string, **required**): The command to run in the test.
- Args ([]string, *optional*): The arguments to pass to the command.
- EnvVars ([]EnvVar, *optional*): A list of environment variables to set for the individual test. See the **Environment Variables** section for more info.
- Expected Output ([]string, *optional*): List of regexes that should match the stdout from running the command.
- Excluded Output ([]string, *optional*): List of regexes that should **not** match the stdout from running the command.
- Expected Error ([]string, *optional*): List of regexes that should match the stderr from running the command.
- Excluded Error ([]string, *optional*): List of regexes that should **not** match the stderr from running the command.
- Exit Code (int, *optional*): Exit code that the command should exit with.

Example:
```yaml
commandTests:
  - name: "gunicorn flask"
    setup: [["virtualenv", "/env"], ["pip", "install", "gunicorn", "flask"]]
    command: "which"
    args: ["gunicorn"]
    expectedOutput: ["/env/bin/gunicorn"]
- name:  "apt-get upgrade"
  command: "apt-get"
  args: ["-qqs", "upgrade"]
  excludedOutput: [".*Inst.*Security.* | .*Security.*Inst.*"]
  excludedError: [".*Inst.*Security.* | .*Security.*Inst.*"]  
```


## File Existence Tests
File existence tests check to make sure a specific file (or directory) exist within the file system of the image. No contents of the files or directories are checked. These tests can also be used to ensure a file or directory is **not** present in the file system.

#### Supported Fields:

- Name (string, **required**): The name of the test
- Path (string, **required**): Path to the file or directory under test
- IsDirectory (boolean, **required**): Whether or not the specified path is a directory (as opposed to a file)
- ShouldExist (boolean, **required**): Whether or not the specified file or directory should exist in the file system
- Permissions (string, *optional*): The expected Unix permission string (e.g.
  drwxrwxrwx) of the files or directory.

Example:
```yaml
fileExistenceTests:
- name: 'Root'
  path: '/'
  isDirectory: true
  shouldExist: true
  permissions: '-rw-r--r--'
```

## File Content Tests
File content tests open a file on the file system and check its contents. These tests assume the specified file **is a file**, and that it **exists** (if unsure about either or these criteria, see the above **File Existence Tests** section). Regexes can again be used to check for expected or excluded content in the specified file.

#### Supported Fields:

- Name (string, **required**): The name of the test
- Path (string, **required**): Path to the file under test
- ExpectedContents (string[], *optional*): List of regexes that should match the contents of the file
- ExcludedContents (string[], *optional*): List of regexes that should **not** match the contents of the file

Example:
```yaml
fileContentTests:
- name: 'Debian Sources'
  path: '/etc/apt/sources.list'
  expectedContents: ['.*httpredir\\.debian\\.org.*']
  excludedContents: ['.*gce_debian_mirror.*']
```

## License Tests
License tests check a list of copyright files and makes sure all licenses are
allowed at Google. By default it will look at where Debian lists all copyright
files, but can also look at an arbitrary list of files.

#### Supported Fields:

- Debian (bool, **required**): If the image is based on Debian, check where
  Debian lists all licenses.
- Files (string[], *optional*): A list of other files to check.

Example:
```yaml
licenseTests:
- debian: true
  files: ["/foo/bar", "/baz/bat"]
```

### Environment Variables
A list of environment variables can optionally be specified as part of the test setup. They can either be set up globally (for all test runs), or test-local as part of individual command test runs (see the **Command Tests** section above). Each environment variable is specified as a key-value pair. Unix-style environment variable substitution is supported.

To specify, add a section like this to your config:

```yaml
globalEnvVars:
  - key: "VIRTUAL_ENV"
    value: "/env"
  - key: "PATH"
    value: "/env/bin:$PATH"
```


### Running Structure Tests Through Bazel
Structure tests can also be run through bazel. To do so, include the rule definitions in your BUILD file:

```BUILD
load("@runtimes_common//structure_tests:tests.bzl", "structure_test")
```

and create a `structure_test` rule, passing in your image and config file as parameters:

```BUILD
docker_build(
    name = "hello",
    base = "//java:java8",
    cmd = ["/HelloJava_deploy.jar"],
    files = [":HelloJava_deploy.jar"],
)

load("@runtimes_common//structure_tests:tests.bzl", "structure_test")

structure_test(
    name = "hello_test",
    config = "testdata/hello.yaml",
    image = ":hello",
)
```