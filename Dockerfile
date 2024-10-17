FROM golang:1.23 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go vet -v ./cmd
RUN go test -v ./...

RUN CGO_ENABLED=1 go build -o /go/bin/app ./cmd

FROM gcr.io/distroless/base-debian12

COPY --from=build /go/bin/app /
CMD ["/app"]