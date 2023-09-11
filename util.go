package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func callURL(aurl string) ([]byte, error) {
	client := &http.Client{Timeout: 15 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	res, err := client.Get(aurl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	return data, err

}

func downloadFile(ctx context.Context, w http.ResponseWriter, url string) (err error) {

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.WithContext(ctx)
	res, err := client.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	// Check server response
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	// Writer the body to file
	w.Header().Add("Content-Length", res.Header.Get("Content-Length"))
	w.Header().Add("Content-type", res.Header.Get("Content-Type"))
	w.Header().Add("Content-Disposition", res.Header.Get("Content-Disposition"))
	w.Header().Add("Content-Filename", strings.ReplaceAll(res.Header.Get("Content-Disposition"), "attachment; filename=", ""))
	_, err = io.Copy(w, res.Body)
	if err != nil {
		return
	}

	return nil
}
