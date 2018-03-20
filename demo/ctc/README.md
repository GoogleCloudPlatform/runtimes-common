# About

This package contains examples for ContainerToolCommand Usage.

# How to Run
* See Available Commands
``` shell
bazel run demo/ctc:ctc_demo --  --help
INFO: Running command line: bazel-bin/demo/ctc/linux_amd64_stripped/ctc_demo --help

Usage:
  echo [flags]
  echo [command]

Available Commands:
  help        Help about any command
  panic       Raises Error
  version     Print the version

Flags:
      --alsoLogToStdOut           Also Log to Std Out
  -h, --help                      help for echo
      --logDir string             LogDir (default "/tmp/")
  -l, --logLevel types.LogLevel   LogLevel (default info)
  -m, --message string            Message to Echo (default "YOUR TEXT TO ECHO")
  -t, --template string           Output format (default "{{.}}")
  -u, --updateCheck               Run Update Check (default true)

Use "echo [command] --help" for more information about a command.


```
* Run echo with command args
``` shell
bazel run demo/ctc:ctc_demo --  --message=ping
INFO: Running command line: bazel-bin/demo/ctc/linux_amd64_stripped/ctc_demo '--message=ping'
ping
```

You can check the logs at
``` shell
{"level":"info","msg":"You are running echo command with message ping","time":"2018-03-20T15:14:11-07:00"}
```

* Run Version Command.
```shell
bazel run demo/ctc:ctc_demo --  version
INFO: Running command line: bazel-bin/demo/ctc/linux_amd64_stripped/ctc_demo version
1.0.1
```

* Run panic Subcommand with --alsoLogToStdOut
```shell
bazel run demo/ctc:ctc_demo --  panic --alsoLogToStdOut
INFO: Running command line: bazel-bin/demo/ctc/linux_amd64_stripped/ctc_demo panic --alsoLogToStdOut
time="2018-03-20T15:15:53-07:00" level=error msg="Oh you called Panic"
ERROR: Non-zero return code '1' from command: Process exited with status 1
```


