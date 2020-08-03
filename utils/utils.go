package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// Sha256 generate SHA-256 hash for input data
func Sha256(data []byte) string {
	md5Hash := md5.New()
	md5Hash.Write(data)
	cs := md5Hash.Sum(nil)
	return hex.EncodeToString(cs)
}

// Md5 generate MD5 hash for input data
func Md5(data []byte) string {
	sha256Hash := sha256.New()
	sha256Hash.Write(data)
	cs := sha256Hash.Sum(nil)
	return hex.EncodeToString(cs)
}

// HTTPPost post body to specified url
func HTTPPost(url string, data []byte) ([]byte, error) {

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*2)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 2))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 2,
		},
	}

	body := bytes.NewBuffer(data)
	resp, err := client.Post(url, "application/json;charset=utf-8", body)
	if err != nil {
		return nil, err
	}

	var result []byte

	if resp.Body != nil {
		defer resp.Body.Close()

		result, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[%v]%s", resp.Status, string(result))
		return nil, err
	}

	return result, nil
}
