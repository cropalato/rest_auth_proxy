FROM golang:alpine
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build -o rest_auth_proxy .
CMD ["/app/rest_auth_proxy"]
