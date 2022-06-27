package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"k8s.io/klog"
)

var (
	cfgFile        string = "./config.yaml"
	pdns_api_url   string = ""
	pdns_api_token string = ""
	listen         string = ":9090"
	proxyRules     map[string][]Authorization
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
	klog.InitFlags(nil)
	defer klog.Flush()

	flag.StringVar(&cfgFile, "config-file", LookupEnvOrString("PAP_CONFIG_FILE", cfgFile), "Config File.")
	flag.StringVar(&listen, "listen", LookupEnvOrString("PAP_LISTEN", listen), "Listenning socket.")
	flag.StringVar(&pdns_api_url, "url", LookupEnvOrString("PAP_PDNS_API_URL", pdns_api_url), "PowerDNS server API URL.")
	flag.StringVar(&pdns_api_token, "pdns-api-token", LookupEnvOrString("PAP_PDNS_API_TOKEN", pdns_api_url), "PowerDNS server API TOKEN.")
	flag.Parse()

	if err := config.loadConfig(cfgFile); err != nil {
		log.Fatalf("Error loading config file %v. %v\n", cfgFile, err)
	}

	if pdns_api_url == "" {
		pdns_api_url = config.Pdns_api_url
	}
	if pdns_api_token == "" {
		pdns_api_token = config.Pdns_api_token
	}
	if len(config.Rules) == 0 {
		log.Fatalf("Missing Rules from config file.")
	}
	proxyRules = config.Rules

	klog.V(5).Info(fmt.Sprintf("Finished process info will be send to %s.", pdns_api_url))

	http.HandleFunc("/", proxyHandler)

	klog.V(1).Info(fmt.Sprintf("Listening http requests from %s.", listen))
	err := http.ListenAndServe(listen, nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}
