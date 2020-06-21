FROM golang:1.13-buster as build

# install protobuf from source
RUN apt-get update && \
    apt-get -y install git unzip build-essential autoconf libtool
RUN git clone https://github.com/google/protobuf.git && \
    cd protobuf && \
    ./autogen.sh && \
    ./configure && \
    make && \
    make install && \
    ldconfig && \
    make clean && \
    cd .. && \
    rm -r protobuf

RUN go get google.golang.org/grpc
RUN go get github.com/golang/protobuf/protoc-gen-go
RUN go get -u github.com/gorilla/mux
RUN go get github.com/jesseokeya/go-httplogger
RUN go get github.com/joho/godotenv
RUN ls -la

WORKDIR /go/src/app
COPY . /go/src/app

RUN go build -o /go/bin/app

FROM gcr.io/distroless/base-debian10
COPY --from=build /go/bin/app /
ARG URL_GRPC="localhost"
EXPOSE 4000
CMD ["/app"]
