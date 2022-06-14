FROM golang:1 as builder
WORKDIR /ports
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ports ./cmd

FROM alpine
COPY --from=builder /ports/ports .
ENTRYPOINT ["./ports"]
