```
 █████████████    ██████  █████ █████
░░███░░███░░███  ███░░███░░███ ░░███ 
 ░███ ░███ ░███ ░███ ░███ ░░░█████░  
 ░███ ░███ ░███ ░███ ░███  ███░░░███ 
 █████░███ █████░░██████  █████ █████
░░░░░ ░░░ ░░░░░  ░░░░░░  ░░░░░ ░░░░░ 
```

**mox** is a tool to stub external dependencies.

# About

It is written in [Go](https://github.com/golang/go) with mappings defined in JSON.

Responses can be generated using [Go templates](https://pkg.go.dev/text/template). It supports [sprig](https://masterminds.github.io/sprig/) template functions. 

Try it together with [bro](https://github.com/lameaux/bro) - a load testing tool.

Check out [nft repo](https://github.com/lameaux/nft) to learn more about **bro** & **mox** for non-functional testing.

# Installation

Make sure you have [Go](https://go.dev/doc/install) installed and `GOPATH` is set up correctly.

Clone this repository and run:

```shell
make install
```

# Usage

## mox

```shell
mox [flags]

--debug
--logJson 
--accessLog
--skipBanner
--port=8080
--adminPort=9090
--configPath=./config
```

### Flags

#### --debug

Enables debug mode. Results in more detailed logging.

#### --logJson

Changes log format to JSON.

#### --accessLog

Requires debug mode (`--debug`).

For all incoming requests, it logs whether they matched any mapping.

#### --skipBanner

Skips printing banner to stdout.

#### --port=8080

Defines a port for mocks handler.

#### --adminPort=9090

Defines a port for admin API (including /metrics endpoint).

#### --configPath=./config

Path to config location with mappings, files and templates.

# API Endpoints

## mocks handler

### user-defined mappings

- [GET /<mapping_url>](http://0.0.0.0:8080/user-defined-mapping)
- POST /<mapping_url>
- ...

### predefined functions

- sleep for N ms [/mox/sleep?seconds=1](http://0.0.0.0:8080/mox/sleep?ms=1)
- sleep for random 0-N ms [/mox/sleep?seconds=1](http://0.0.0.0:8080/mox/sleep/random?ms=5)

## admin handler

### api

- [http://mox:9090/api](http://0.0.0.0:9090/api)

### metrics

- [http://mox:9090/metrics](http://0.0.0.0:9090/metrics)