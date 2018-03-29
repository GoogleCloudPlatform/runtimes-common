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
bazel run demo/ctc:ctc_demo --  --message=ping --alsoLogToStdOut
INFO: You are running echo command with message ping
ping
```

You can check the logs at
``` shell
{"level":"info","msg":"You are running echo command with message ping","time":"2018-03-20T15:14:11-07:00"}
```

* Run Version Command.
```shell
bazel run demo/ctc:ctc_demo --  version
1.0.1
```

* Run panic Subcommand
```shell
bazel run demo/ctc:ctc_demo --  panic
ERROR: Oh you called Panic
ERROR: Non-zero return code '1' from command: Process exited with status 1
```

* Run config command to see command configurations.
```shell
go run  demo/ctc/main.go  config
logdir : /tmp
message : echo
updatecheck : false

```
* Set Config Variable to a new value
```shell
$go run  demo/ctc/main.go  config set message hi
Config key Changed and written to file demo/ctc/testConfig.json

$cat demo/ctc/testConfig.json
{
  "logdir": "/tmp",
  "message": "hi",
  "update_check": "false"
}
```


