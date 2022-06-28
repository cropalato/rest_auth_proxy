FROM golang:alpine
RUN mkdir /app
RUN mkdir /app/tmp
COPY . /app/tmp
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o rest_auth_proxy ./tmp/.
RUN rm -rf /app/tmp
CMD ["/app/rest_auth_proxy"]
