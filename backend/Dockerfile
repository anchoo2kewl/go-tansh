# Dockerfile.distroless
FROM golang:1.18-bullseye as base

ENV APP_HOME /go/src/tansh
RUN mkdir -p "$APP_HOME"

WORKDIR "$APP_HOME"

COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /main .

FROM gcr.io/distroless/static-debian11 as production

COPY --from=base /main .

CMD ["./main"]