FROM alpine as build
RUN apk update && apk add alpine-sdk go

COPY . /server
WORKDIR /server
RUN go build

FROM alpine
COPY --from=build /server/server /server
ENTRYPOINT [ "/server" ]
