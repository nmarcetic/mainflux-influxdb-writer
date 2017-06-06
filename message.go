/**
 * Copyright (c) Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/cisco/senml"
	ic "github.com/influxdata/influxdb/client/v2"
)

// writeMessage function
// Writtes message into DB.
func writeMessage(nm NatsMsg) error {

	var s senml.SenML
	var err error

	if s, err = senml.Decode(nm.Payload, senml.JSON); err != nil {
		println("ERROR")
		return err
	}

	// Normalize (i.e. resolve) SenMLRecord
	sn := senml.Normalize(s)
	if len(sn.Records) == 0 {
		err := errors.New("ERROR: Wrong SenML")
		log.Print(err)
		return err
	}

	// Timestamp
	t := time.Now().UTC().Format(time.RFC3339)

	// New InfluxDB point batch
	bp, err := ic.NewBatchPoints(ic.BatchPointsConfig{
		Database:  InfluxDBName,
		Precision: InfluxTimePrecision,
	})

	for _, r := range sn.Records {

		// InfluxDB tags
		tags := map[string]string{
			"name":         r.Name,
			"content_type": nm.ContentType,
		}

		// InfluxDB fields
		fields := make(map[string]interface{})

		if len(r.Unit) != 0 {
			fields["unit"] = r.Unit
		}
		if r.UpdateTime != 0 {
			fields["update_time"] = r.UpdateTime
		}
		if r.Value != nil {
			fields["value"] = *(r.Value)
		}
		if len(r.StringValue) != 0 {
			fields["string_value"] = r.StringValue
		}
		if len(r.DataValue) != 0 {
			fields["data_value"] = r.DataValue
		}
		if r.BoolValue != nil {
			fields["bool_value"] = *(r.BoolValue)
		}
		if r.Sum != nil {
			fields["sum"] = *(r.Sum)
		}

		fields["channel"] = nm.Channel
		fields["content_type"] = nm.ContentType
		fields["publisher"] = nm.Publisher
		fields["protocol"] = nm.Protocol
		fields["created"] = t

		pt, err := ic.NewPoint(nm.Channel, tags, fields, time.Unix(int64(r.Time), 0))
		if err != nil {
			log.Print(err)
			return err
		}
		bp.AddPoint(pt)
	}

	// Write the batch
	if err := InfluxClient.Write(bp); err != nil {
		log.Print(err)
		return err
	}

	fmt.Println("Msg written")
	return nil
}
