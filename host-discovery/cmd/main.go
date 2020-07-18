package main

import (
	"fmt"
	"log"

	"github.com/td0m/smarthome/host-discovery/pkg/api"
)

func main() {
	token, err := api.GetToken()
	if err != nil {
		log.Panic(err)
	}
	clients, err := api.GetClientList(token)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("%+v\n", clients)
}
