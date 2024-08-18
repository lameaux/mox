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
--skipBanner
--port=8080
--adminPort=8081
--metricsPort=9090
```
