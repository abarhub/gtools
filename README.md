# gtools

![last build](https://github.com/abarhub/gtools/actions/workflows/tests.yml/badge.svg)

![Last release](https://img.shields.io/github/v/release/abarhub/gtools)


Tools for files

to build executable
```shell
go build ./cmd/gtools.go
```

to build for linux
```shell
set GOOS=linux
set GOARCH=amd64
go build ./cmd/gtools.go
```

to run tests
```shell
go test ./...
```

man
```
Usage:
  gtools [flags]
  gtools [command]

Available Commands:
  base64      encode/decode in base64
  completion  Generate the autocompletion script for the specified shell
  copy        copy files
  du          disk usage
  help        Help about any command

Flags:
  -h, --help   help for gtools

Use "gtools [command] --help" for more information about a command.
```
