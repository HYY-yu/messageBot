FROM golang:1.20 AS builder

WORKDIR /build

ENV GOPROXY https://goproxy.cn,direct
COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o app .

FROM  ccr.ccs.tencentyun.com/hyy-yu/alpine-timezone:3.15 AS final
WORKDIR /app

# Copy the controller binary from the context into our
# container image
COPY --from=builder /build/app /app

ENTRYPOINT ["./app"]

# docker build -f deploy/Dockerfile -t bot:v1 .
# docker run -d -p 8080:8080 bot:v1