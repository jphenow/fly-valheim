FROM golang as proxy-build

COPY proxy /code
WORKDIR /code
RUN go build

FROM ghcr.io/lloesche/valheim-server

COPY --from=proxy-build /code/proxy /usr/local/bin/udp-proxy
COPY fly-bootstrap.sh .

CMD ["./fly-bootstrap.sh"]
