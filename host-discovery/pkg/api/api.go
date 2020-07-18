package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func GetToken() (string, error) {
	var loginBody LoginBody
	err := post("?fname=system&opt=login&function=set&usrid=f6fdffe48c908deb0f4c3bd36c032e72", &loginBody)
	if err != nil {
		return "", err
	}
	return loginBody.Token, nil
}

func GetClientList(token string) ([]Client, error) {
	var clientListBody ClientListBody
	err := post(fmt.Sprintf("?token=%s&fname=system&opt=main&function=get&math=%f", token, rand.Float32()), &clientListBody)
	if err != nil {
		return []Client{}, nil
	}
	return clientListBody.Terminals, nil
}

func post(params string, bodyPtr interface{}) error {
	res, err := http.Post("http://192.168.10.1/protocol.csp"+params, "application/json", nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(bodyBytes, bodyPtr)
	return nil
}
