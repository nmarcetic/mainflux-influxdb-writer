/**
 * Copyright (c) Mainflux
 *
 * Mainflux server is licensed under an Apache license, version 2.0.
 * All rights not explicitly granted in the Apache license, version 2.0 are reserved.
 * See the included LICENSE file for more details.
 */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cisco/senml"
	ic "github.com/influxdata/influxdb/client/v2"
	"github.com/mainflux/mainflux-influxdb-writer/models"
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

	// Timestamp
	t := time.Now().UTC().Format(time.RFC3339)

	// New InfluxDB point batch
	bp, err := ic.NewBatchPoints(ic.BatchPointsConfig{
		Database:  InfluxDBName,
		Precision: InfluxTimePrecision,
	})

	for _, r := range sn.Records {

		m := models.Message{}

		// Copy SenMLRecord struct to Message struct
		b, err := json.Marshal(r)
		if err != nil {
			log.Print(err)
			return err
		}
		if err := json.Unmarshal(b, &m); err != nil {
			log.Print(err)
			return err
		}

		// Fill-in Mainflux stuff
		m.Channel = nm.Channel
		m.Publisher = nm.Publisher
		m.Protocol = nm.Protocol
		m.Timestamp = t

		// Create an InfluxDB point and add to batch
		if m.Value != nil {
			tags := map[string]string{"name": m.Name}

			fields := map[string]interface{}{
				"value": *(m.Value),
			}

			pt, err := ic.NewPoint(nm.Channel, tags, fields, time.Unix(int64(m.Time), 0))
			if err != nil {
				log.Print(err)
				return err
			}
			bp.AddPoint(pt)
		}
	}

	// Write the batch
	if err := InfluxClient.Write(bp); err != nil {
		log.Print(err)
		return err
	}

	fmt.Println("Msg written")
	return nil
}
