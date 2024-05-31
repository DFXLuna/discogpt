FROM golang:1.22.3-bookworm AS builder

ENV CGO_ENABLED 0

WORKDIR /go/src
COPY main.go go.mod go.sum LICENSE /go/src/
COPY pkg /go/src/pkg
RUN go mod download

RUN go build -o /go/bin/app .

FROM alpine:latest AS runner
LABEL org.opencontainers.image.source https://github.com/dfxluna/discogpt

RUN apk add --upgrade apk-tools & \
    apk upgrade --available

WORKDIR /discogpt

COPY --from=builder /go/bin/app /discogpt/

CMD [ "./app" ]
