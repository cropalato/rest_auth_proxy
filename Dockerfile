FROM golang:alpine
RUN mkdir /app
RUN mkdir /app/tmp
COPY . /app/tmp
WORKDIR /app
RUN go build -o rest_auth_proxy ./tmp/.
RUN rm -rf /app/tmp
CMD ["/app/rest_auth_proxy"]
