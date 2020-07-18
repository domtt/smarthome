package api

type bodyMeta struct {
	Opt string `json:"opt"`
}

type LoginBody struct {
	bodyMeta
	Token string `json:"token"`
}

type ClientListBody struct {
	bodyMeta
	Terminals []Client `json:"terminals"`
}

type Client struct {
	Mac       string     `json:"mac"`
	IP        string     `json:"ip"`
	Speed     uint       `json:"speed"`
	UpSpeed   uint       `json:"up_speed"`
	UpBytes   uint       `json:"up_bytes"`
	DownBytes uint       `json:"down_bytes"`
	Name      string     `json:"name"`
	Vendor    uint       `json:"vendor"`
	OSType    uint       `json:"ostype"`
	Uptime    uint       `json:"uptime"`
	OnTime    uint       `json:"ontime"`
	Mode      uint       `json:"mode"`
	Pic       uint       `json:"pic"`
	Flag      string     `json:"flag"`
	Sig       string     `json:"sig"`
	Apps      []struct{} `json:"apps"` // TODO: apps
}
