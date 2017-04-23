###
# Mainflux InfluxDB Writer Dockerfile
###

FROM golang:alpine
MAINTAINER Mainflux

ENV INFLUX_HOST influx
ENV INFLUX_PORT 8086

ENV NATS_HOST nats
ENV NATS_PORT 4222

###
# Install
###
# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/mainflux/mainflux-influxdb-writer
RUN cd /go/src/github.com/mainflux/mainflux-influxdb-writer && go install

###
# Run main command with dockerize
###
CMD mainflux-influxdb-writer -i INFLUX_HOST -n NATS_HOST

