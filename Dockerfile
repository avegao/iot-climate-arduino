FROM golang:1.9.2-alpine AS build

WORKDIR /go/src/github.com/avegao/iot-climate-arduino

RUN apk add --no-cache --update \
    git \
    glide

COPY glide.yaml glide.yaml
COPY glide.lock glide.lock

RUN glide install

COPY ./ ./

ARG VCS_REF="unknown"
ARG BUILD_DATE="unknown"

RUN go test ./... -cover &&\
    go install \
        -ldflags "-X main.buildDate=$BUILD_DATE -X main.commitHash=$VCS_REF"

########################################################################################################################

FROM alpine:3.7

MAINTAINER "Álvaro de la Vega Olmedilla <alvarodlvo@gmail.com>"

ENV GRPC_VERBOSITY ERROR

RUN addgroup iot-climate-arduino && \
    adduser -D -G iot-climate-arduino iot-climate-arduino

USER iot-climate-arduino

WORKDIR /app

COPY --from=build /go/bin/iot-climate-arduino /app/iot-climate-arduino

EXPOSE 50000/tcp

LABEL com.avegao.iot.temp.vcs_ref=$VCS_REF \
      com.avegao.iot.temp.build_date=$BUILD_DATE \
      maintainer="Álvaro de la Vega Olmedilla <alvarodlvo@gmail.com>"

ENTRYPOINT ["./iot-climate-arduino"]
