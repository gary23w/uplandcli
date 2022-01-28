# UPLAND CLI

[![GoDev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://dothissomeday.com)

**A small go application to collect and view upland related blockchain information and raw data**

## Prerequisites

You must have the following versions of Node.js, npm, and Go installed:

- ~~Node.js: [Node.js](https://nodejs.org/en/) v14.17.6~~ currently not required
- ~~NPM: [npm](https://www.npmjs.com/) > v6.14.x~~ currently not required

- Go: [Go](https://golang.org/doc/install) > v1.16.x

## Database

Currently setup to deploy a postgresql database onto heroku.

_Make sure you are signed into heroku frst._

```
go run main.go database --deploy
```

see configuration for more details.

## Basic Commands

```shell
go run main.go upland --help
```

collects and displays the data within a TermUI interface.

```
go run main.go upland --collect
```

collects and tails the data in your shell.

```
go run main.go upland --live

```

## Configuration

the current collection methods rely on a **POSTGRES** database.
Update utils/database.json if you choose to use your own data center.

```
{
    "Url": "null",
    "PSQLurl": "null",
    "User": "<DB NAME>",
    "Password": "<DB PASSWORD>",
    "Host": "<DB HOST>",
    "Port": "<PORT>",
    "Database": "<DATABASE>"
}
```

## CLI Goals

- A small tool to live load all blockchain information related to upland.me
- Create a pipeline that collects from eos.hyperion.eosrio.io
- Transform data to selected outputs
- Create a local database for collected data
- Implement graphing and various analytics systems.
  - Analyze longterm data to predict possible trends.

## API Docs

A small beego/swagger CRUD api is implemented

```
GET /upland/properties       | Get properties from database.
```

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

- [BEEGO]("https://github.com/beego/beego")

## Acknowledgements

- Shout out to UPLAND.ME for an amazing metaverse experience.

  - https://api.upland.me

- EOS Hyperion
  - https://eos.hyperion.eosrio.io
