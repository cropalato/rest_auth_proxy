# rest_auth_proxy

This HTTP header reverse proxy can be used to extend autorization REST API call based on request header.

The authorization rules will validate a request header and match the HTTP method with URL path. In case it matches, the proxy will override http request header by remote server api token and forward the request to the remote server.

In case there is no rule or not match the rule, it will forward the request without override.


## Supported env variables to configure the service:

| Variable | Description |
| -------- | ----------- |
| RAP_CONFIG_FILE | Config file path. |
| RAP_RULES | Rules in json format. |
| RAP_LISTEN | ip_address:port used to listen client requests. |
| RAP_HEADER_KEY | HTTP request header used to authorize/override the request. |
| RAP_API_URL | Remote API server URL. |
| RAP_API_TOKEN | Token used in forward request when match rule. |


## Config file

In case you dont use the env variables, the app will try to collect the info from a config file. By default it is ./.config.yaml.

``` yaml
---
listen: "127.0.0.1:9998"
server_api_url: "http://api.foo.acme:8081"
server_api_token: "Naiheocohkese6saephaiziquaineeHa"
header_token: "X-API-Key"
rules:
  ye4icu7xeengeenio9aiSheev7raiv8A:
    - method: "GET"
      pathregex:
        - "/api/v1/servers/localhost/zones/foo.acme"
```


## arguments

If you wish you can use the following arguments to configure your service.

| Argument | env variable |
| -------- | ------------ |
| config-file | RAP_CONFIG_FILE |
| rules | RAP_RULES |
| listen | RAP_LISTEN |
| header-key | RAP_HEADER_KEY |
| url | RAP_API_URL |
| server-api-token | RAP_API_TOKEN |


## docker image

You can find a docker image from [rest_auth_proxy](https://hub.docker.com/r/cropalato/rest_auth_proxy).

### Running docker with config file

```bash
docker run --name rest_auth_proxy -p 9998:9000 --rm -v "$(pwd)/.config.yaml:/app/.config.yaml" cropalato/rest_auth_proxy:latest
```

### Running docker with env variable

```bash
docker run --name rest_auth_proxy \
-e 'RAP_RULES={"ye4icu7xeengeenio9aiSheev7raiv8A":[{"method":"GET","pathregex":["/api/v1/servers/localhost/zones/foo.acme"]}]}' \
-e 'RAP_API_TOKEN=Naiheocohkese6saephaiziquaineeHa' \
-e 'RAP_API_URL=http://dns.acme.foo:8081' \
-p 9998:9000 --rm cropalato/rest_auth_proxy:latest
```
