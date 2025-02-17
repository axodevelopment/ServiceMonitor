FROM golang:1.23.0 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o service .

FROM golang
WORKDIR /
COPY --from=builder /app/service /service
ENTRYPOINT ["/service"]