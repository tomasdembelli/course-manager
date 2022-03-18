# syntax=docker/dockerfile:1

# Alpine is chosen for its small footprint
# compared to Ubuntu
FROM golang:1.18-alpine

RUN apk --no-cache add git

WORKDIR /go/src/github.com/tomasdembelli/course-manager

# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go install -v ./cmd/...

EXPOSE 8000

CMD ["/go/bin/course-manager"]