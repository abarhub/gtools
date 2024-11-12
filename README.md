# gtools

![last build](https://github.com/abarhub/gtools/actions/workflows/tests.yml/badge.svg)

![Last release](https://img.shields.io/github/v/release/abarhub/gtools)


Tools for files

to build executable for windows
```shell
set GOOS=windows
set GOARCH=amd64
go build ./cmd/gtools.go
```

to build for linux
```shell
set GOOS=linux
set GOARCH=amd64
go build ./cmd/gtools.go
```

to build for raspberry
```shell
set GOOS=linux
set GOARCH=arm
set GOARM=5
go build ./cmd/gtools.go
```

to run tests
```shell
go test ./...
```

man
```
gtools is a super simple CLI tools
   
simple in CLI

Usage:
  gtools [flags]
  gtools [command]

Available Commands:
  base64      encode/decode in base64
  completion  Generate the autocompletion script for the specified shell
  copy        copy files
  du          disk usage
  help        Help about any command
  ls          list files
  merge       merge files
  mv          move files
  password    generate password
  rename      rename files
  rm          remove files
  split       split file
  time        time execution of command
  unzip       unzip directory
  zip         zip directory

Flags:
  -h, --help      help for gtools
  -v, --version   version for gtools

Use "gtools [command] --help" for more information about a command.
```
