package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
	for k := range r.Header {
		err = h.requestAuthz(r.Method, r.URL.Path, k, r.Header.Get(k))
		if err != nil {
			klog.V(5).Info(fmt.Sprintf("No rule for %v %v using header %v:%v.", r.Method, r.URL.Path, k, r.Header.Get(k)))
			continue
		}
	}
	if override = err == nil; override {
		klog.V(3).Info(fmt.Sprintf("Forwarding request without changes."))
	}
	new_url := fmt.Sprintf("%s%s", pdns_api_url, r.URL)
	pr, err := forwardRequest(new_url, r, override)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed forwaring request. %v", err))
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
		req.Header.Set(header_key, pdns_api_token)
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
