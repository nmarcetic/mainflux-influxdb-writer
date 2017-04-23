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
	"log"
	"os"
	"runtime"

	"github.com/cenkalti/backoff"
	"github.com/fatih/color"
	"github.com/nats-io/go-nats"
)

const (
	help string = `
Usage: mainflux-influxdb [options]
Options:
	-i, --ihost	InfluxDB host
	-q, --iport	InfluxDB port
	-s, --db	InfluxDB database
	-u, --user	InfluxDB username
	-p, --pass	InfluxDB password
	-p, --precision	InfluxDB time precision
	-n, --nhost	NATS host
	-m, --nport	NATS port
	-h, --help	Show help
`
)

type (
	Opts struct {
		InfluxHost      string
		InfluxPort      string
		InfluxDatabase  string
		InfluxUser      string
		InfluxPass      string
		InfluxPrecision string

		NatsHost string
		NatsPort string

		Help bool
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
	opts     Opts
)

func tryInfluxInit() error {
	var err error

	log.Print("Connecting to InfluxDB... ")
	if err = InfluxInit(opts.InfluxHost, opts.InfluxPort, opts.InfluxDatabase,
		opts.InfluxUser, opts.InfluxPass, opts.InfluxPrecision); err != nil {
		log.Println(err)
		return err
	}

	if _, err = InfluxQueryDB(fmt.Sprintf("CREATE DATABASE %s", opts.InfluxDatabase)); err != nil {
		log.Println(err)
	}

	return err
}

func tryNatsConnect() error {
	var err error

	log.Print("Connecting to NATS... ")
	NatsConn, err = nats.Connect("nats://" + opts.NatsHost + ":" + opts.NatsPort)
	return err
}

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
	writeMessage(m)
	fmt.Println(m.Publisher, m.Protocol, m.Channel, m.Payload)
}

func main() {

	flag.StringVar(&opts.InfluxHost, "i", "localhost", "InfluxDB host.")
	flag.StringVar(&opts.InfluxPort, "q", "8086", "InfluxDB port.")
	flag.StringVar(&opts.InfluxDatabase, "d", "mainflux", "InfluxDB databse name.")
	flag.StringVar(&opts.InfluxUser, "u", "mainflux", "InfluxDB username.")
	flag.StringVar(&opts.InfluxPass, "s", "", "InfluxDB password.")
	flag.StringVar(&opts.InfluxPrecision, "p", "s", "InfluxDB time precision.")
	flag.StringVar(&opts.NatsHost, "n", "localhost", "NATS broker address.")
	flag.StringVar(&opts.NatsPort, "m", "4222", "NATS broker port.")
	flag.BoolVar(&opts.Help, "h", false, "Show help.")
	flag.BoolVar(&opts.Help, "help", false, "Show help.")

	flag.Parse()

	if opts.Help {
		fmt.Printf("%s\n", help)
		os.Exit(0)
	}

	// Connect to InfluxDB
	if err := backoff.Retry(tryInfluxInit, backoff.NewExponentialBackOff()); err != nil {
		log.Fatalf("InfluxDB: Can't connect: %v\n", err)
	} else {
		log.Println("OK")
	}

	/*

	 */
	// Connect to NATS broker
	if err := backoff.Retry(tryNatsConnect, backoff.NewExponentialBackOff()); err != nil {
		log.Fatalf("NATS: Can't connect: %v\n", err)
	} else {
		log.Println("OK")
	}

	// Subscribe to NATS
	NatsConn.Subscribe("msg.*", influxdbHandler)

	// Print banner
	color.Cyan(banner)

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
