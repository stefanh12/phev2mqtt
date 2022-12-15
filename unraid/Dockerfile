FROM golang:alpine

WORKDIR /usr/src/app
COPY .  /usr/src/app
RUN apk update && apk add --no-cache \
	g++ \
	gcc \
	git \
	bash \
	libpcap-dev
RUN git clone https://github.com/stefanh12/phev2mqtt.git
COPY --from=golang:alpine /usr/local/go/ /usr/local/go/

COPY . .
RUN cd /usr/src/app/phev2mqtt && \
    /usr/local/go/bin/go build
RUN apk update && apk add --no-cache \
	libpcap-dev

# COPY ./entrypoint.sh /
RUN chmod +x ./entrypoint.sh

ENV  mqtt_server 192.168.1.1:1883
ENV  mqtt_user user
ENV  mqtt_password password
ENV  phev_register false
ENV  debug true


ENTRYPOINT ["bash", "./entrypoint.sh"]
