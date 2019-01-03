FROM alpine as build
RUN apk update && apk add alpine-sdk go

RUN mkdir -p /go/src
COPY . /go/src/hyperkit-measure-time-drift
WORKDIR /go/src/hyperkit-measure-time-drift/cmd/server
ENV GOPATH=/go
RUN go build

FROM alpine
COPY --from=build /go/src/hyperkit-measure-time-drift/cmd/server/server /server
ENTRYPOINT [ "/server" ]
