package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"strconv"
)

var (
	lg, lg_err            = syslog.New(syslog.LOG_INFO, "pdnsAPI_auth_proxy")
	pdns_api_url   string = "http://localhost:8100"
	pdns_api_token string = ""
	listen         string = ":9090"
	debugMode      bool   = false
)

func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func LookupEnvOrBool(key string, defaultVal bool) bool {
	if val, ok := os.LookupEnv(key); ok {
		v, err := strconv.ParseBool(val)
		if err != nil {
			log.Fatalf("LookupEnvOrBool[%s]: %v", key, err)
		}
		return v
	}
	return defaultVal
}

func main() {
	if lg_err != nil {
		log.Fatalln(lg_err)
	}

	flag.StringVar(&listen, "listen", LookupEnvOrString("PAP_LISTEN", listen), "Listenning socket.")
	flag.StringVar(&pdns_api_url, "url", LookupEnvOrString("PAP_PDNS_API_URL", pdns_api_url), "PowerDNS server API URL.")
	flag.StringVar(&pdns_api_token, "pdns-api-token", LookupEnvOrString("PAP_PDNS_API_TOKEN", pdns_api_url), "PowerDNS server API TOKEN.")
	flag.BoolVar(&debugMode, "debug", LookupEnvOrBool("PAP_DEBUG", debugMode), "Enable debug mode.")
	flag.Parse()
	lg.Debug(fmt.Sprintf("Finished process info will be send to %s.", pdns_api_url))

	if debugMode {
		lg, lg_err = syslog.New(syslog.LOG_DEBUG, "pdnsAPI_auth_proxy")
	}

	if lg_err != nil {
		log.Fatalln(lg_err)
	}

	if debugMode {
		fmt.Printf("%s=%v\n", "listen", listen)
		fmt.Printf("%s=%v\n", "pdns_api_url", pdns_api_url)
		fmt.Printf("%s=%v\n", "debugMode", debugMode)
	}

	http.HandleFunc("/", proxyHandler)

	lg.Info(fmt.Sprintf("Listening http requests from %s.", listen))
	err := http.ListenAndServe(listen, nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}
