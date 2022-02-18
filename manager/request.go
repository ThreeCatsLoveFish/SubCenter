package manager

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Get(rawUrl string) ([]byte, error) {
	resp, err := http.Get(rawUrl)
	if err != nil || resp == nil || resp.StatusCode != 200 {
		// FIXME: handle error
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func GetWithParams(rawUrl string, params url.Values) ([]byte, error) {
	structUrl, err := url.Parse(rawUrl)
	if err != nil {
		// FIXME: handle error
		return nil, err
	}
	structUrl.RawQuery = params.Encode()
	resp, err := http.Get(structUrl.String())
	if err != nil || resp == nil || resp.StatusCode != 200 {
		// FIXME: handle error
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func PostForm(url string, data url.Values) error {
	resp, err := http.PostForm(url, data)
	if err != nil || resp == nil || resp.StatusCode != 200 {
		// FIXME: handle error
		return err
	}
	resp.Body.Close()
	return nil
}

func PostFormWithCookie(url string, data url.Values, cookie string) error {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		// FIXME: handle error
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Connection", "keep-alive")
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp == nil || resp.StatusCode != 200 {
		// FIXME: handle error
		return err
	}
	resp.Body.Close()
	// TODO: cover body content
	return nil
}
