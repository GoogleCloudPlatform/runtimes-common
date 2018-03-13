# About

This package contains examples for ContainerToolCommand Usage.

# How to Run
* See Available Commands
``` shell
go run main.go --help

Usage:
  echo [flags]
  echo [command]

Available Commands:
  help        Help about any command
  version     Print the version of echo

Flags:
  -h, --help              help for echo
  -m, --message string    Message to Echo (default "text")
  -t, --template string   Output format (default "{{.}}")

Use "echo [command] --help" for more information about a command.

```
* Run echo with command args
``` shell
cd ctc_lib/examples
go run main.go --message ping
{ping}
```
* Run echo with --template argument
``` shell
cd ctc_lib/examples
go run main.go --message ping --template {{.Message}}
ping

```
* Run Version Command.
```shell
cd ctc_lib/examples
go run main.go version --template={{.Version}}
0.0.1
```


