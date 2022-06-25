package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

type proxyResp struct {
	headers    http.Header
	body       []byte
	statusCode int
}

type requesAuthz struct {
	Method string   `json:method`
	URLs   []string `json:urls`
}

func requestAuthz(method string, url string, token string) error {
	authz_map := make(map[string][]requesAuthz)
	authz_map["6WoSWhBpZahwZXjP53gu5zkrWEYbivMTT"] = append(authz_map["6WoSWhBpZahwZXjP53gu5zkrWEYbivMTT"], requesAuthz{Method: "GET", URLs: []string{"/api/v1/servers/localhost/zones/ludia.me"}})
	authz, ok := authz_map[token]
	if !ok {
		return errors.New("authz: Invalid token.")
	}
	for _, a := range authz {
		if match, _ := regexp.MatchString(a.Method, method); !match {
			continue
		}
		for _, u := range a.URLs {
			if u == url {
				return nil
			}

		}
	}
	return errors.New("authz: Permission deny.")
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	err := requestAuthz(r.Method, r.URL.Path, r.Header.Get("X-API-Key"))
	if err != nil {
		fmt.Printf("Ops. %v\n", err)
		w.WriteHeader(403)
		w.Write([]byte("Access denied.\n"))
		return
	}
	new_url := fmt.Sprintf("%s%s", pdns_api_url, r.URL)
	pr, err := forwardRequest(new_url, r)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed forwaring request. %v", err))
	}
	w.WriteHeader(pr.statusCode)
	w.Write(pr.body)
}

func forwardRequest(url string, r *http.Request) (proxyResp, error) {
	pr := new(proxyResp)
	method := r.Method
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return *pr, fmt.Errorf("Got error %s", err.Error())
	}
	for k, v := range r.Header {
		if len(v) > 1 {
			for _, value := range v {
				req.Header.Add(k, value)
			}
		} else {
			req.Header.Set(k, v[0])
		}
	}
	req.Header.Set("X-API-Key", pdns_api_token)
	response, err := client.Do(req)
	if err != nil {
		return *pr, fmt.Errorf("Got error %s", err.Error())
	}
	body, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	pr.statusCode = response.StatusCode
	pr.headers = response.Header
	pr.body = body

	return *pr, nil
}
