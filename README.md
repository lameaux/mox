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

It is written in Go with mappings defined in JSON.

Responses can be generated using [Go templates](https://pkg.go.dev/text/template). We support [sprig](https://masterminds.github.io/sprig/) template functions. 

Try it together with [bro](https://github.com/lameaux/bro) - a load testing tool.

Check out [nft](https://github.com/lameaux/nft) repo to learn more about bro & mox for non-functional testing.

# Installation

Make sure you have `GOPATH` set up correctly.

```shell
make install
```

# Usage

```shell
mox [flags]

--debug 
--accessLog
--skipBanner
--port=8080
--adminPort=8081
--metricsPort=9090
```

## Flags

### --debug

Enables debug mode. Results in more detailed logging.

### --accessLog

Requires debug mode (`--debug`).

For all incoming requests, it logs whether they matched any mapping.

### --skipBanner

Skips printing banner to stdout

### --port=

Defines a port for mock handler

### --adminPort

Defines a port for admin API

### --metricsPort

Defines a port for metrics endpoint
