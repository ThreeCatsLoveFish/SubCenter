package manager

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func Post(url string, data string) {
	client := &http.Client{}

	fmt.Println(url)
	fmt.Println(data)
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		// FIXME: handle error
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		// FIXME: handle error
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// FIXME: handle error
	}

	fmt.Println(string(body))
}
