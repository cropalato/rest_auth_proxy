FROM golang:alpine
RUN mkdir -p /app/tmp
COPY . /app/tmp
WORKDIR /app/tmp
RUN go build -o ../rest_auth_proxy .
WORKDIR /app
RUN rm -rf /app/tmp
ENTRYPOINT ["/app/rest_auth_proxy"]
