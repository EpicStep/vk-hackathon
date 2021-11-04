FROM golang:alpine AS builder

WORKDIR $GOPATH/src/github.com/EpicStep/vk-hackathon/

COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o /go/bin/back ./cmd/back/main.go

FROM alpine:latest
COPY --from=builder /go/bin/back /go/bin/
CMD ["/go/bin/back"]