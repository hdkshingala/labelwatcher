FROM golang:1.18.0
WORKDIR /labelwatcher
ADD . .
RUN go mod download && CGO_ENABLED=0 go build

FROM scratch
WORKDIR /labelwatcher
COPY --from=0 labelwatcher .
ENTRYPOINT [ "./labelwatcher" ]
