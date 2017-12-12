FROM golang

# Fetch common go package dependencies.
# RUN go get github.com/eclipse/paho.mqtt.golang
RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get github.com/onsi/gomega

RUN go get github.com/c9s/goprocinfo/linux

# The resulting contianer image associated with this Dockerfile can be
# used to create cross-compiled runnables for the Raspberry Pi.

WORKDIR /go/src/github.com/andrew-bodine/informix

COPY . .

# Build next level testing candidate.
RUN go build -o /go/bin/informix main.go

# Run unit and integration tests.
RUN ginkgo -r --race --cover

EXPOSE 80
