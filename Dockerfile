# sea-server-go/Dockerfile

FROM golang:alpine as builder

RUN apk --no-cache add git

WORKDIR /app/sea-server-go

# Copy the current code into our workdir
COPY . .

RUN go mod download

RUN sh ./script/build_image.sh

FROM alpine:latest

# Security related package, good to have.
RUN apk --no-cache add ca-certificates

# Same as before, create a directory for our app.
RUN mkdir /app
WORKDIR /app
RUN mkdir bin
RUN mkdir conf
RUN mkdir log

COPY --from=builder /app/sea-server-go/sea-server-go bin/
COPY seago-test.toml conf/

CMD ["./bin/sea-server-go"]