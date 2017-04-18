# Mainflux InfluxDB Service

[![License](https://img.shields.io/badge/license-Apache%20v2.0-blue.svg)](LICENSE)
[![Build Status](https://travis-ci.org/mainflux/mainflux-influxdb-writer.svg?branch=master)](https://travis-ci.org/mainflux/mainflux-influxdb-writer)
[![Go Report Card](https://goreportcard.com/badge/github.com/mainflux/mainflux-influxdb-writer)](https://goreportcard.com/report/github.com/mainflux/mainflux-influxdb-writer)
[![Join the chat at https://gitter.im/Mainflux/mainflux](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/Mainflux/mainflux?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

InfluxDB (https://www.influxdata.com/) service of Mainflux IoT Cloud.

Connects to NATS and consumes all messages sent via all Mainflux protocols (HTTP, MQTT, WS, CoAP). Writes these messages into InfluxDB.

### Installation
#### Prerequisite
If not set already, please set your `GOPATH` and `GOBIN` environment variables. For example:
```bash
mkdir -p ~/go
export GOPATH=~/go
export GOBIN=$GOPATH/bin
# It's often useful to add $GOBIN to $PATH
export PATH=$PATH:$GOBIN
```

#### Get the code
Use [`go`](https://golang.org/cmd/go/) tool to "get" (i.e. fetch and build) `mainflux-influxdb-writer` package:
```bash
go get github.com/mainflux/mainflux-influxdb-writer
```

This will download the code to `$GOPATH/src/github.com/mainflux/mainflux-influxdb-writer` directory,
and then compile it and install the binary in `$GOBIN` directory.

Now you can run the program with:
```
mainflux-influxdb-writer
```
if `$GOBIN` is in `$PATH` (otherwise use `$GOBIN/mainflux-influxdb-writer`)

### Documentation
Development documentation can be found [here](http://mainflux.io/).

### Community
#### Mailing lists
[mainflux](https://groups.google.com/forum/#!forum/mainflux) Google group.

#### IRC
[Mainflux Gitter](https://gitter.im/Mainflux/mainflux?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

#### Twitter
[@mainflux](https://twitter.com/mainflux)

### License
[Apache License, version 2.0](LICENSE)
