package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"k8s.io/klog"
)

type proxyResp struct {
	headers    http.Header
	body       []byte
	statusCode int
}

func (h *headerRules) requestAuthz(method string, url string, header_key string, header_value string) error {

	for k, v := range *h {
		if k != header_value {
			continue
		}
		for _, a := range v {
			if a.Method != method {
				continue
			}
			for _, u := range a.PathRegEx {
				if match, _ := regexp.MatchString(u, url); match {
					if klog.V(5) {
						klog.Info(fmt.Sprintf("Matched"))
					}
					return nil
				}

			}
		}
	}
	return errors.New("authz: Permission deny.")
}

func (h *headerRules) proxyHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var override bool
	err = nil
	header_token_req := r.Header.Get(header_token)
	if header_token_req != "" {
		err = h.requestAuthz(r.Method, r.URL.Path, header_token, header_token_req)
		if err != nil {
			if klog.V(5) {
				klog.Info(fmt.Sprintf("No rule for %v %v using header %v:%v.", r.Method, r.URL.Path, header_token, header_token_req))
			}
		}
	}
	if override = err == nil; !override {
		if klog.V(3) {
			klog.Info(fmt.Sprintf("Forwarding request without changes."))
		}
	}
	new_url := fmt.Sprintf("%s%s", server_api_url, r.URL)
	pr, err := forwardRequest(new_url, r, override)
	if err != nil {
		klog.Fatal(fmt.Sprintf("Failed forwaring request. %v", err))
	}
	w.WriteHeader(pr.statusCode)
	w.Write(pr.body)
}

func forwardRequest(url string, r *http.Request, override bool) (proxyResp, error) {
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
	if override {
		req.Header.Set(header_token, server_api_token)
	}
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
