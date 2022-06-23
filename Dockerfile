FROM golang:1 as builder

RUN apt-get update && apt-get install -y --no-install-recommends protobuf-compiler
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /ports
COPY go.mod go.sum ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux make build

FROM alpine
COPY --from=builder /ports/bin/ports .
ENTRYPOINT ["./ports"]
