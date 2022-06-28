FROM golang:alpine
RUN mkdir /app
COPY . /app
WORKDIR /app
CMD ["/app/rest_auth_proxy"]
