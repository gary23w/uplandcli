# UPLANDCLI

UPLANDCLI is a compact, highly capable Go application that streamlines the collection and visualization of Upland blockchain information and raw data. The tool facilitates the creation of a local database to store the collected data and employs graphing and analytics systems to offer insights into the long-term trends of the Upland ecosystem.

### Prerequisites

To build and run UPLANDCLI, you need:

- Go: Go >= 1.17 (see `go.mod`)

Notes:
- Node.js/npm are not required for the current CLI/API flow. References to browser scraping via `chromedp` are future/planned.

### Basic Commands

The `upland` command initiates a pipeline to collect data related to Upland properties. The collected data is then used to populate the CLI-based user interface.

#### Usage:

```
uplandcli upland [flags]

uplandcli upland --collect
uplandcli upland --live
uplandcli upland --live -a      # start API in the background
uplandcli upland --live -a -b   # bypass database inserts
```

#### Flags:

```
-a, --api: Run API in async mode.
-b, --bypass: Bypass database connections and inserts.
-d, --collect: Retrieve all recently listed properties for sale.
-h, --help: Display help for the upland command.
-q, --live: Display live mode data in the shell.
```

---

The database command is used to set up and initialize a PostgreSQL database on Heroku.

##### Usage:

```
uplandcli database [flags]

-q, --check: Check if a database is already active.
-d, --deploy: Set up a PostgreSQL database on Heroku.
-u, --destroy: Attempt to destroy the database.
-h, --help: Display help for the database command.
```

---

The api command deploys a CRUD API to interact with the database.

##### Usage:

```
uplandcli api [flags]
Flags:

-d, --deploy: Initialize a CRUD API to interact with the database.
-h, --help: Display help for the api command.
```

#### Configuration

UPLANDCLI relies on a PostgreSQL database to collect data. To configure the collection system, a database.json file must be available within the conf/ directory. Here's an example:

Create the directory and files if they don’t exist:

```
mkdir conf
```

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

UPLANDCLI also uses a .conf file to configure the API. This file should be available within the conf/ directory. Here's an example:

```
    appname = MYAPPNAME
    httpport = 1337
    runmode = dev
    autorender = false
    copyrequestbody = true
    EnableDocs = true
    sqlconn = postgres://<username>:<password>@<ip>:<port>/<database>?sslmode=require
```

#### Example Logging Configuration

UPLANDCLI reads logging config from `conf/logging.json`.
Below is an example of what this file might look like:

```
{
  "level": "info",
  "encoding": "console",
  "outputPaths": ["stdout", "tmp/logs"],
  "errorOutputPaths": ["tmp/errorlogs"],
  "initialFields": { "initFieldKey": "fieldValue" },
  "encoderConfig": {
    "messageKey": "message",
    "levelKey": "level",
    "nameKey": "logger",
    "timeKey": "time",
    "callerKey": "logger",
    "stacktraceKey": "stacktrace",
    "callstackKey": "callstack",
    "errorKey": "error",
    "timeEncoder": "iso8601",
    "fileKey": "file",
    "levelEncoder": "capitalColor",
    "durationEncoder": "second",
    "callerEncoder": "full",
    "nameEncoder": "full",
    "sampling": {
      "initial": "3",
      "thereafter": "10"
    }
  }
}
```

##### CLI Goals

The Upland CLI has several goals:

- Load all blockchain information related to upland.me
- Collect data from eos.hyperion.eosrio.io
- Transform data to selected outputs
- Create a local database for collected data
- Implement graphing and various analytics systems
- Analyze long-term data to predict possible trends

### API Documentation

The Upland CLI also includes a small beego/swagger CRUD API. The following endpoint is currently available:

```
GET /upland/properties       | Get properties from database.
```

**You can also view the collected data from your web browser by starting the API and visiting:**

http://127.0.0.1:1337/upland/properties/analysis

Data is displayed in descending order, with the most recently collected data at the top of the list.

## Built With:

The Upland CLI was built using the following technologies:

- [Golang]("https://go.dev/")

- chromedp (planned/future)

- [TermUI]("https://github.com/gizak/termui")

- [Cobra]("https://github.com/spf13/cobra")

- [ZAP]("go.uber.org/zap")

- [BEEGO]("https://github.com/beego/beego")

##### Acknowledgements

We would like to give a special shoutout to Upland.me for providing an amazing metaverse experience. We would also like to thank EOS Hyperion for providing their services at https://eos.hyperion.eosrio.io.

### Performance Notes

Recent performance-focused changes include:

- HTTP requests reuse a pooled client with timeouts and keep-alives.
- JSON responses are decoded using streaming decoders to reduce allocations.
- Database inserts are batched in a single transaction with a prepared statement.
- UI/live loops avoid unbounded goroutine buildup.
