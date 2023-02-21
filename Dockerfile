FROM golang:1.18-alpine3.15 as builder
RUN apk --no-cache add build-base

WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go vet -v
RUN CGO_ENABLED=0 go build -o /go/bin/app -v

FROM scratch
COPY --from=builder /go/bin/app /
EXPOSE 8080
CMD ["/app"]
