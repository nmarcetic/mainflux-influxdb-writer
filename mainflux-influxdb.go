/**
 * Copyright (c) 2017 Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/nats-io/go-nats"
	"log"
	"os"
	"runtime"
)

const (
	help string = `
Usage: mainflux-influxdb [options]
Options:
	-a, --host	Host address
	-p, --port	Port
	-h, --help	Prints this message end exits`
)

type (
	Opts struct {
		NatsHost string
		NatsPort string
		Help     bool
	}

	NatsMsg struct {
		Channel   string `json:"channel"`
		Publisher string `json:"publisher"`
		Protocol  string `json:"protocol"`
		Payload   []byte `json:"payload"`
	}
)

var (
	NatsConn *nats.Conn
)

func influxdbHandler(nm *nats.Msg) {
	fmt.Printf("Received a message: %s\n", string(nm.Data))

	m := NatsMsg{}
	if len(nm.Data) > 0 {
		if err := json.Unmarshal(nm.Data, &m); err != nil {
			println("Can not decode MQTT msg")
			return
		}
	}

	println("Calling writeMessage()")
	fmt.Println(m.Publisher, m.Protocol, m.Channel, m.Payload)
}

func main() {
	opts := Opts{}

	flag.StringVar(&opts.NatsHost, "a", "nats", "NATS broker address.")
	flag.StringVar(&opts.NatsPort, "p", "4222", "NATS broker port.")
	flag.BoolVar(&opts.Help, "h", false, "Show help.")
	flag.BoolVar(&opts.Help, "help", false, "Show help.")

	flag.Parse()

	if opts.Help {
		fmt.Printf("%s\n", help)
		os.Exit(0)
	}

	// Connect to NATS broker
	var err error
	NatsConn, err = nats.Connect("nats://" + opts.NatsHost + ":" + opts.NatsPort)
	if err != nil {
		log.Fatalf("NATS: Can't connect: %v\n", err)
	}

	// Print banner
	color.Cyan(banner)

	// Subscribe to NATS
	NatsConn.Subscribe("mainflux/core/out", influxdbHandler)

	// Prevent program to exit
	runtime.Goexit()
}

var banner = `
╔╦╗┌─┐┬┌┐┌┌─┐┬  ┬ ┬─┐ ┬  ╦┌┐┌┌─┐┬  ┬ ┬─┐ ┬╔╦╗╔╗ 
║║║├─┤││││├┤ │  │ │┌┴┬┘  ║│││├┤ │  │ │┌┴┬┘ ║║╠╩╗
╩ ╩┴ ┴┴┘└┘└  ┴─┘└─┘┴ └─  ╩┘└┘└  ┴─┘└─┘┴ └─═╩╝╚═╝

          == Industrial IoT System ==
       
         Made with <3 by Mainflux Team
[w] http://mainflux.io
[t] @mainflux

`
