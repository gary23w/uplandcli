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

~upland~

```
++UPLD-CLI UPLAND++
========================
The UPLD-PIPELINE will query from the blockchain and collect data related
to Upland properties. This data will be used to populate the CLI based user interface.

Example:
upldcli upland --collect
upldcli upland --live
upldcli upland --live -a  // run API in async mode
upldcli upland --live -a -b  // run API in async mode and bypass database connections

The UPLD-PIPELINE will also scrape the Upland website and collect data via a headless browser.
using chromedp and chromedp-go for headless browsing. This is a future implementation, and should be available soon.

Usage:
uplandcli upland [flags]

Flags:
-a, --api       run API in async mode
-b, --bypass    bypass database connections and inserts
-d, --collect   will get all of the recent properties listed for sale.
-h, --help      help for upland
-q, --live      live mode which tails collected data in your shell.
```

~database~

```
++UPLD-DB UPLAND++
========================
The DB command is used to setup and initialize a postgresql database on heroku.

Usage:
uplandcli database [flags]

Flags:
-q, --check     checks to see if a database is already active
-d, --deploy    will setup a postgresql database on heroku
-u, --destroy   will attempt to destroy the database
-h, --help      help for database
```

~api~

```
++UPLD-DB UPLAND++
========================
The api command is used to deploy a crud api to interact with the database.

Usage:
uplandcli api [flags]

Flags:
-d, --deploy   will initialize a crud api to interact with the database
-h, --help     help for api
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

---

You can also view the data in an html table format.
Simply fire up the API and go to:

`http://127.0.0.1:1337/upland/properties/analysis`

in your browser. Data here is pulled in DESC order. Meaning that whichever is at the top of the list will be the most recently collected data.

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
