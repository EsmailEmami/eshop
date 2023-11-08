FROM golang:1.20

RUN mkdir -p /build
RUN mkdir -p /app

ADD . /build/

WORKDIR /build
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /app/main .

EXPOSE 6060

WORKDIR /app

COPY ./oauthkeys .
COPY ./config.toml .

ENTRYPOINT ["/app/main"]
CMD ["serve"]