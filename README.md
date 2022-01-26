# UPLAND CLI

[![GoDev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://dothissomeday.com)

**A small go application to collect and view upland related blockchain information and raw data**

## Prerequisites

You must have the following versions of Node.js, npm, and Go installed:

- Node.js: [Node.js](https://nodejs.org/en/) v14.17.6
- NPM: [npm](https://www.npmjs.com/) > v6.14.x

- Go: [Go](https://golang.org/doc/install) > v1.16.x

## Basic CLI cmds

```shell
go run main.go upland --help
go run main.go upland --collect
go run main.go upland --live

```

## Database

Currently setup to deploy a postgresql database onto heroku.
Make sure you are signed into heroku and then execute.

```
go run main.go database --deploy

```

the current collection methods rely on the database so be sure to have one installed.

## Goals

- A small tool to live load all blockchain information related to upland.me
- Create a pipeline that collects from eos.hyperion.eosrio.io
- Transform data to selected outputs
- Create a local database for collected data
- Implement graphing and various analytics systems.
  - Analyze longterm data to predict possible trends.

## Generated Documentation

To generate the documentation, run the following command:

```shell
godoc-static -destination=docs ./
cd docs && python -m http.server 8000
```

Ensure that both godoc and godoc-static are installed.

- `go get -d code.rocketnine.space/tslocum/godoc-static`
- `go get -d golang.org/x/tools/cmd/godoc`

## Built With

- [Golang]("https://go.dev/")
- [chromedp]("https://github.com/chromedp/chromedp")
- [TermUI]("https://github.com/gizak/termui")
- [Cobra]("https://github.com/spf13/cobra")
- [ZAP]("go.uber.org/zap")
