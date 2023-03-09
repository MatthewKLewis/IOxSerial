# IOx Application that reads from Serial Port
Simple Python and Go applications that scan from a COM port, formats the data, and posts to an API. This applications are intended to be containerized in a Docker image, and converted via ioxclient to a Cisco AP runnable program.

The activation json connects the linux container running the program from the inside-out - the package yaml from the outside-in

---

# pull base docker image for binary builds
FROM --platform=$TARGETARCH alpine:latest AS builder
RUN apk update
RUN apk upgrade
RUN apk add --no-cache build-base curl wget openssl rabbitmq-c jq

# download source for lrzsz binaries
RUN wget https://ohse.de/uwe/releases/lrzsz-0.12.20.tar.gz

$extract the source
RUN tar xvf lrzsz-0.12.20.tar.gz

# compile and install lrzsz binaries in /usr/local/bin/
RUN cd lrzsz-0.12.20 && ./configure
RUN cd lrzsz-0.12.20 && make
RUN cd lrzsz-0.12.20 && make install


# start from base rootfs again for the final docker image
FROM --platform=$TARGETARCH alpine:latest
RUN apk update
RUN apk upgrade
RUN apk add --no-cache curl wget openssl rabbitmq-c jq

# copy lrzsz compiled above in builder work space
COPY --from=builder /usr/local/bin/* /usr/bin/
COPY --from=builder /usr/local/lib/* /usr/lib/