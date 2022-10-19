FROM golang:1.19 as build
ENV CGO_ENABLED=0

WORKDIR /go/src/app

COPY . .

RUN go mod download
RUN go test ./...
RUN go build -o app

FROM gcr.io/distroless/static

COPY --from=build --chown=nonroot:nonroot /go/src/app/app /app
COPY --from=build --chown=nonroot:nonroot /go/src/app/.env.docker /.env.docker
 
ENTRYPOINT ["/app"]
