/**
 * Copyright (c) 2017 Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */
package main

import (
	"github.com/influxdata/influxdb/client/v2"
)

var (
	InfluxClient        client.Client
	InfluxDBName        string
	InfluxTimePrecision string
)

func InfluxInit(host string, port string, db string, username string, password string, precision string) error {
	/** Connect to NATS broker */
	var err error

	InfluxDBName = db
	InfluxTimePrecision = precision

	// Create a new HTTPClient
	InfluxClient, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://" + host + ":" + port,
		Username: username,
		Password: password,
	})

	return err
}

// queryDB convenience function to query the database
func InfluxQueryDB(cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: InfluxDBName,
	}
	if response, err := InfluxClient.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}
