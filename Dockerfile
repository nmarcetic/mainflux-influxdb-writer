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

# Dockerize
ENV DOCKERIZE_VERSION v0.2.0
ADD https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz dockerize.tar.gz
RUN tar -C /usr/local/bin -xzf dockerize.tar.gz && rm -f dockerize.tar.gz


###
# Run main command with dockerize
###
CMD dockerize -wait tcp://$INFLUX_HOST:$INFLUX_PORT \
				-wait tcp://$NATS_HOST:$NATS_PORT \
				-timeout 10s /go/bin/mainflux-influxdb-writer -i INFLUX_HOST -n NATS_HOST

