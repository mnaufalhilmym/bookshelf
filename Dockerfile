FROM golang:1.23 as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go vet -v ./cmd
RUN go test -v ./cmd

RUN CGO_ENABLED=0 go build -o /go/bin/app ./cmd

FROM gcr.io/distroless/static-debian12

COPY --from=build /go/bin/app /
CMD ["/app"]