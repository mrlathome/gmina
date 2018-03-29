// package mina forked from github.com/sariina/mina

package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
)

func writeBodyToWR(wr http.ResponseWriter, resp *http.Response) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("\033[0;31mError: %s\033[0m", err)
		return
	}
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	wr.Write(body)
}

func cacheWrite(path string, filename string, body []byte) {
	err := os.MkdirAll(path, 0755)
	if err != nil {
		log.Printf("Error while mkdir: %s", err)
		return
	}

	err = ioutil.WriteFile(filename, body, 0644)
	if err != nil {
		log.Printf("Error while writing: %s", err)
		return
	}
}

func isFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func requestMD5(req *http.Request) (string, []byte) {
	h := md5.New()
	body, _ := httputil.DumpRequest(req, true)
	io.WriteString(h, fmt.Sprintf("%+v", string(body)))

	return fmt.Sprintf("%x", h.Sum(nil)), body
}

func (g *GMina) request(req *http.Request) (*http.Response, error) {
	md5, reqDump := requestMD5(req)
	reqFilename := filepath.Join(g.CacheDir, fmt.Sprintf("%s.req", md5))
	resFilename := filepath.Join(g.CacheDir, fmt.Sprintf("%s.res", md5))

	if isFileExist(resFilename) {
		resDump, err := ioutil.ReadFile(resFilename)
		if err != nil {
			return nil, err
		}
		dumpIO := bufio.NewReader(bytes.NewBuffer(resDump))
		resp, err := http.ReadResponse(dumpIO, req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	} else {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}

		resDump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return nil, err
		}

		go cacheWrite(g.CacheDir, resFilename, resDump)
		go cacheWrite(g.CacheDir, reqFilename, reqDump)

		return resp, nil
	}
}
