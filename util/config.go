package util

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
)

var (
	ActInfo AccountInfo
	NetInfo EosNetInfo
)

type AccountInfo struct {
	Systemaccount []string `json:"systemaccount"`
	Sys_publickey string `json:"system_publickey"`
	Sys_privatekey string `json:"system_privatekey"`
	Voteaccount []string `json:"voteaccount"`
	Produceraccount []string `json:"produceraccount"`
	Voters [][]string `json:"voters"`
	Producers [][]string `json:"producers"`
	NodeProducer []string `json:"nodeproducer"`
}

type EosNetInfo struct {
	Hosturl string `json:"hosturl"`
} 

func (ai *AccountInfo) LoadFile(filepath string) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, ai)
	if err != nil {
		return err
	}

	return nil
}

func (net *EosNetInfo) LoadFile(filepath string) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, net)
	if err != nil {
		return err
	}

	return nil
}

func init()  {
	aifilepath := "./config/accountinfo.json"
	nifilepath := "./config/netinfo.json"
	err := ActInfo.LoadFile(aifilepath)
	if err != nil {
		fmt.Println("get error: ", err)
	} else {
		fmt.Println("Config init OK")
	}

	err = NetInfo.LoadFile(nifilepath)
	if err != nil {
		fmt.Println("get error: ", err)
	} else {
		fmt.Println("Config init OK")
	}
}