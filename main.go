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
	cfgFile        string = "./.config.yaml"
	serverAPIURL   string = ""
	serverAPIToken string = ""
	headerToken    string = ""
	listen         string = ""
	ruleEnv        string = ""
	ruleMap        headerRules
)

// LookupEnvOrString returns the value from env variable key is exists or defaultVal as string
func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

// LookupEnvOrBool returns the value from env variable key is exists or defaultVal as boolean
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
	flag.StringVar(&ruleEnv, "rules", LookupEnvOrString("RAP_RULES", ruleEnv), "json format rules.")
	flag.StringVar(&listen, "listen", LookupEnvOrString("RAP_LISTEN", listen), "Listenning socket.")
	flag.StringVar(&headerToken, "header-key", LookupEnvOrString("RAP_HEADER_KEY", headerToken), "Header key used to authentication.")
	flag.StringVar(&serverAPIURL, "url", LookupEnvOrString("RAP_API_URL", serverAPIURL), "Remote server API URL.")
	flag.StringVar(&serverAPIToken, "server-api-token", LookupEnvOrString("RAP_API_TOKEN", serverAPIURL), "Remote server API TOKEN.")
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
	if headerToken == "" {
		if config.HeaderToken != "" {
			headerToken = config.HeaderToken
		} else {
			headerToken = "X-API-Key"
		}
	}
	if serverAPIURL == "" {
		if config.ServerAPIURL == "" {
			klog.Fatalf("Missing Remote Server API URL.")
		}
		serverAPIURL = config.ServerAPIURL
	}
	if serverAPIToken == "" {
		if config.ServerAPIToken == "" {
			klog.Fatalf("Missing remote API token")
		}
		serverAPIToken = config.ServerAPIToken
	}
	if err := json.Unmarshal([]byte(ruleEnv), &ruleMap); err != nil {
		if len(config.Rules) == 0 {
			klog.Warning("No rules defined via environment variables or config file.")
		}
		ruleMap = config.Rules
	}

	if klog.V(5) {
		klog.Info(fmt.Sprintf("Finished process info will be send to %s.", serverAPIURL))
	}

	http.HandleFunc("/", ruleMap.proxyHandler)

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
