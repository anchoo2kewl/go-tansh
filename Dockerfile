FROM golang:1.18-bullseye

RUN go install github.com/beego/bee/v2@latest

ENV GO111MODULE=on

ENV APP_HOME /go/src/tansh
RUN mkdir -p "$APP_HOME"

WORKDIR "$APP_HOME"

ADD . .
RUN go mod download github.com/go-chi/chi/v5

EXPOSE 3000
CMD ["bee", "run"]