# --- Build stage ---
FROM golang:1.18 AS builder
WORKDIR /app
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io
# Copy the dependency definition
COPY ./go.mod .
COPY ./go.sum .
# Download dependencies
RUN go mod download
# Copy the source code
COPY . .
# Build for release
RUN CGO_ENABLED=0 GOOS=linux go build

# --- Final stage ---
FROM alpine:3

# Install AWS CLI
RUN apk add --no-cache \
        python3 \
        py3-pip \
        jq \
    && pip3 install --upgrade pip \
    && pip3 install --no-cache-dir \
        awscli \
    && rm -rf /var/cache/apk/*

RUN aws --version

COPY --from=builder /app/auth-service /
COPY ./deploy/scripts/app/init.sh ./deploy/scripts/app/keys.sh ./deploy/scripts/app/entrypoint.sh /

RUN chmod +x /init.sh && chmod +x /keys.sh && chmod +x entrypoint.sh
ENTRYPOINT ["/init.sh", "/keys.sh", "/entrypoint.sh"]

EXPOSE 3000
LABEL Name=auth-service

CMD /auth-service
