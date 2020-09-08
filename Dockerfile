# FROM golang:alpine as builder
# RUN mkdir /build 
# ADD . /build/
# WORKDIR /build
# COPY go.mod .
# COPY go.sum .

# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o server api/main.go

# FROM scratch
# COPY --from=builder /build/server /app/
# WORKDIR /app

# ENTRYPOINT ["./server"]


FROM alpine
ADD ./bin/users-api-linux .

ENTRYPOINT ["./users-api-linux"]