FROM golang:latest

# Set destination for COPY
WORKDIR /app

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
#COPY ./main .
COPY . .

# 设置go GOPROXY
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

ENV REDIS_HOST=10.168.1.8
ENV REDIS_PORT=6379

#安装brew
RUN apt-get update && apt-get install -y curl && \
    curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh | bash


# Build
RUN brew install FiloSottile/musl-cross/musl-cross
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8011

# Run
CMD ["./main"]