Master [![Circle CI](https://circleci.com/gh/eris-ltd/eris-db/tree/master.svg?style=svg)](https://circleci.com/gh/eris-ltd/eris-db)
 Develop [![Circle CI (develop)](https://circleci.com/gh/eris-ltd/eris-db/tree/develop.svg?style=svg)](https://circleci.com/gh/eris-ltd/eris-db)

# Eris-DB (Alpha)

`eris-db` is Eris' blockchain-client. It consists of a [Tendermint](http://tendermint.com/) node wrapped by a simple server. The server allows requests to be made over HTTP - either using JSON-RPC 2.0 or a RESTlike web-api - and websocket (JSON-RPC 2.0). The web-APIs are documented in `api.md`. We also have javascript bindings for the RPC methods in [eris-db.js](https://github.com/eris-ltd/eris-db.js).


# TMSP-TODO

- second genesis file
- get tests working
- eris-cli wrapping



## Installation

There are no pre-built releases other then the docker images.

The recommended way of working with eris-db is through [eris-cli](https://github.com/eris-ltd/eris-cli) (develop branch as of now).

### Building from source

#### Ubuntu 14.04 (OSX ?)

Make sure you have the proper [Go](https://golang.org/) distribution for your OS and architecture. The recommended version is `1.4.2`. Follow the instructions on how to set up GOROOT and GOPATH.

You will also need the following libraries: `git, libgmp3-dev`

On Ubuntu: `sudo apt-get install git libgmp3-dev`

On Mac: `brew install git gmp`

Next you pull in the code:

`go get github.com/eris-ltd/eris-db/cmd/erisdb`

This will build and install the `erisdb` executable and put it in `$GOPATH/bin`, which should be on your PATH. If not, then add it.

To run `erisdb`, just type `$ erisdb /path/to/working/folder`

This will start the node using the provided folder as working dir. If the path is omitted it defaults to `~/.erisdb` 

#### Docker

It is best to use [eris-cli](https://github.com/eris-ltd/eris-cli) which will help setting up and running eris-db (and individual chains) through docker.

##### Others

Tendermint officially supports only 64 bit Ubuntu. 

### Usage

####Native

The simplest way to get started is by simply running `$ erisdb`. That will start a fresh node with `~/.erisdb` as the working directory, and the default settings. You will be asked to type in a hostname, which could be anything. `anonymous` is a commonly used one.

Once the server has started, it will begin syncing up with the network. At that point you may begin using it. The preferred way is through our [javascript api](https://github.com/eris-ltd/erisdb-js), but it is possible to connect directly via HTTP or websocket. The JSON-RPC and web-api reference can be found [here](https://github.com/eris-ltd/erisdb/blob/master/api.md).

### Configuration

There will be more info on how to set up a private net when this is added to Tendermint. That would include information about the various different fields in `config.toml`, `genesis.json`, and `priv_validator.json`.

#### server_conf.toml

The server configuration file looks like this:

```
[bind]
  address= <string>
  port= <number>
[TLS]
  tls= <boolean>
  cert_path= <string>
  key_path= <string>
[CORS]
  enable            <boolean>
  allow_origins     <[]string>
  allow_credentials <boolean>
  allow_methods     <[]string>
  allow_headers     <[]string>
  expose_headers    <[]string>
  max_age           <number>
[HTTP]
  json_rpc_endpoint= <string>
[web_socket]
  websocket_endpoint= <string>
  max_websocket_sessions= <number>
  read_buffer_size = <number>
  write_buffer_size = <number>
[logging]
  console_log_level= <string>
  file_log_level= <string>
  log_file= <string>
```

**NOTE**: **CORS** and **TLS** are not yet fully implemented, and cannot be used. CORS is implemented through [gin middleware](https://github.com/tommy351/gin-cors), and TLS through the standard Go http package via the [graceful library](https://github.com/tylerb/graceful).

##### Bind

- `address` is the address.
- `port` is the port number

##### TLS

- `tls` is used to enable/disable TLS
- `cert_path` is the absolute path to the certificate file.
- `key_path` is the absolute path to the key file.

##### CORS

- `enable` is whether or not the CORS middleware should be added at all. 

Details about the other fields and how this is implemented can be found [here](https://github.com/tommy351/gin-cors).

##### HTTP

- `json_rpc_endpoint` is the name of the endpoint used for JSON-RPC (2.0) over HTTP.

##### web_socket

- `websocket_endpoint` is the name of the endpoint that is used to establish a websocket connection.
- `max_websocket_connections` is the maximum number of websocket connections that is allowed at the same time.
- `read_buffer_size` is the size of the read buffer for each socket in bytes.
- `read_buffer_size` is the size of the write buffer for each socket in bytes.

##### logging

- `console_log_level` is the logging level used for the console.
- `file_log_level` is the logging level used for the log file.
- `log_file` is the path to the log file. Leaving this empty means file logging will not be used.

The possible log levels are these: `crit`, `error`, `warn`, `info`, `debug`.

The server log level will override the log level set in the Tendermint `config.toml`.

##### example server_conf.toml file

```
[bind]
address="0.0.0.0"
port=1337
[TLS]
tls=false
cert_path=""
key_path=""
[CORS]
enable=false
allow_origins=[]
allow_credentials=false
allow_methods=[]
allow_headers=[]
expose_headers=[]
max_age=0
[HTTP]
json_rpc_endpoint="/rpc"
[web_socket]
websocket_endpoint="/socketrpc"
max_websocket_sessions=50
read_buffer_size = 4096
write_buffer_size = 4096
[logging]
console_log_level="info"
file_log_level="warn"
log_file=""
```

### Server-server

**NOTE: This feature is being deprecated in favor of `eris-cli` generation of configurable throw-away chains.**

The library includes a "server-server". This server accepts POST requests with some chain data (such as priv_validator.json and genesis.json), and will use that to create a new working directory in the temp folder, write the data, deploy a new node in that folder, generate a port, use it to serve that node and then pass the url back in the response. It will also manage all the servers and shut them down as they become inactive. 

NOTE: This is not safe in production, as it requires private keys to be passed over a network, but it is useful when doing tests. If the same chain data is used, then each node is  guaranteed to give the same output (for the same input) when calling the methods.

To start one up, just run `go install` in the `erisdb/cmd/erisdbss` directory, then run `erisdbss`. It takes no parameters. There are many examples on how to call it in the javascript library, and if people find it useful there will be a  tutorial.

### Testing

In root: `go test ./...`

### Benchmarking

As of `0.11.0`, there are no benchmarks. We aim to have a framework built before `1.0`.
