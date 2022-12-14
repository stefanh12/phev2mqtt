FROM golang:alpine

WORKDIR /usr/src/app
COPY .  /usr/src/app
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
# COPY go.mod go.sum ./
# RUN go mod download && go mod verify
RUN apk update && apk add --no-cache \
	g++ \
	gcc \
	git \
	bash \
	libpcap-dev
# WORKDIR /tmp
RUN git clone https://github.com/stefanh12/phev2mqtt.git
COPY --from=golang:alpine /usr/local/go/ /usr/local/go/

COPY . .
RUN cd /usr/src/app/phev2mqtt && \
    /usr/local/go/bin/go build
RUN apk update && apk add --no-cache \
	libpcap-dev

# COPY --from=builder /tmp/phev2mqtt/phev2mqtt /opt/phev2mqtt
# RUN go mod download && go mod verify
# RUN go build -v -o /usr/src/app/phev2mqtt ./...

# COPY ./entrypoint.sh /
RUN chmod +x ./entrypoint.sh

ENV  mqtt_server 192.168.1.1:1883
ENV  mqtt_user user
ENV  mqtt_password password
ENV  phev_register false
ENV  debug true

# mqtt_server 1
# mqtt_user 2
# mqtt_password 3
# phev_register 4
# debug 5

# ENTRYPOINT ["/bin/sh", "-c", "/entrypoint.sh  ${mqtt_server} ${mqtt_user} ${mqtt_password} ${phev_register} ${debug}"]
#SHELL ["/bin/sh", "-c", "./entrypoint.sh  ${mqtt_server} ${mqtt_user} ${mqtt_password} ${phev_register} ${debug}"]
ENTRYPOINT ["bash", "./entrypoint.sh"]
