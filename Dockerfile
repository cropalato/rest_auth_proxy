FROM golang:alpine
RUN apk upgrade --update-cache --available && \
    apk add openssl && \
    rm -rf /var/cache/apk/*
RUN mkdir -p /app/tmp
COPY . /app/tmp
WORKDIR /app/tmp
RUN go build -o ../rest_auth_proxy .
WORKDIR /app
RUN rm -rf /app/tmp
ENTRYPOINT ["/app/rest_auth_proxy"]
