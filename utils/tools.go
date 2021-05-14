package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Get(url string) (*BwCounter, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var ret BwCounter
	err = json.Unmarshal(body, &ret)
	return &ret, err
}
