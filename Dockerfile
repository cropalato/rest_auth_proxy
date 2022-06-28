FROM golang:alpine
RUN mkdir /app
COPY rest_auth_proxy /app
WORKDIR /app
CMD ["/app/rest_auth_proxy"]
