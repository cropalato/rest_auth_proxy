package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"k8s.io/klog"
)

var (
	cfgFile          string = "./.config.yaml"
	server_api_url   string = ""
	server_api_token string = ""
	header_token     string = ""
	listen           string = ""
	rule_env         string = ""
	rule_map         headerRules
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
			klog.Fatal(fmt.Sprintf("LookupEnvOrBool[%s]: %v", key, err))
		}
		return v
	}
	return defaultVal
}

func main() {
	klog.InitFlags(nil)
	defer klog.Flush()

	flag.StringVar(&cfgFile, "config-file", LookupEnvOrString("RAP_CONFIG_FILE", cfgFile), "Config File.")
	flag.StringVar(&rule_env, "rules", LookupEnvOrString("RAP_RULES", rule_env), "json format rules.")
	flag.StringVar(&listen, "listen", LookupEnvOrString("RAP_LISTEN", listen), "Listenning socket.")
	flag.StringVar(&header_token, "header-key", LookupEnvOrString("RAP_HEADER_KEY", header_token), "Header key used to authentication.")
	flag.StringVar(&server_api_url, "url", LookupEnvOrString("RAP_API_URL", server_api_url), "Remote server API URL.")
	flag.StringVar(&server_api_token, "server-api-token", LookupEnvOrString("RAP_API_TOKEN", server_api_url), "Remote server API TOKEN.")
	flag.Parse()

	if err := config.loadConfig(cfgFile); err != nil {
		klog.Fatal(fmt.Sprintf("Error loading config file %v. %v\n", cfgFile, err))
	}

	if listen == "" {
		if config.Listen != "" {
			listen = config.Listen
		} else {
			listen = ":9000"
		}
	}
	if header_token == "" {
		if config.Header_token != "" {
			header_token = config.Header_token
		} else {
			header_token = "X-API-Key"
		}
	}
	if server_api_url == "" {
		if config.Server_api_url == "" {
			klog.Fatalf("Missing Remote Server API URL.")
		}
		server_api_url = config.Server_api_url
	}
	if server_api_token == "" {
		if config.Server_api_token == "" {
			klog.Fatalf("Missing remote API token")
		}
		server_api_token = config.Server_api_token
	}
	if err := json.Unmarshal([]byte(rule_env), &rule_map); err != nil {
		if len(config.Rules) == 0 {
			klog.Warning("No rules defined via environment variables or config file.")
		}
		rule_map = config.Rules
	}

	if klog.V(5) {
		klog.Info(fmt.Sprintf("Finished process info will be send to %s.", server_api_url))
	}

	http.HandleFunc("/", rule_map.proxyHandler)

	if klog.V(1) {
		klog.Info(fmt.Sprintf("Listening http requests from %s.", listen))
	}
	err := http.ListenAndServe(listen, nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}

}
