FROM golang:1.15-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh
        
WORKDIR /go/src/term-frequency
COPY . .
        
COPY go.mod go.sum ./
        
RUN go mod download
        
CMD ["/bin/bash"]