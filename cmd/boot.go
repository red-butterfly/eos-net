package cmd

import (
	"fmt"

	"../eosapi"
	"../script"
)

func StartNodeos(info string) {
	script.StartEOSNode()
}

func Boot()  {
	api := eosapi.GetEOS([]string{
		"5K463ynhZoCDDa4RDcr63cUwWLTnKqmdcoTKTHBjqoKfv4u5V7p",
		"5KYvKrAr3a3Qv17pkfLtqERTsowubspcArRzzaVvf6gDxsZmxnX",
	})

	infoResp, _ := api.Api.GetInfo()
	fmt.Println("Get info:", infoResp.VirtualBlockCPULimit)
	accountResp, _ := api.Api.GetAccount("eosio")
	fmt.Println("Permission for initn:", accountResp)

	api.InitEosioNode()
	api.CreateVoters()
	api.RegProducer()
	api.ListProducers()
}

func Vote() {
	api := eosapi.GetEOS([]string{})
	api.VoteProducer()
}

func Resign() {
	api := eosapi.GetEOS([]string{
		"5K463ynhZoCDDa4RDcr63cUwWLTnKqmdcoTKTHBjqoKfv4u5V7p",
	})
	api.StepResign()
}

func Test() {
	api := eosapi.GetEOS([]string{
		"5JGhSHzNDCL8oSLht1N5Tx3Hg1mBHjVy8FXwgqG5pd6Etoi9tS1",
		"5KhZiTkG8S4dQQVe9Wn7VYBpABYEAV7mKh3wu4AMPWFycsGnmfh",
		"5Hz6miKviYhpThMTL6XYHMdFDjXZaVJZTBVwAksQDWrTczaefVu",
		"5JxZUseYxcBmx4zc4TAMfXbHCC258pUV14RAJVW6zbZxme5Dr2n",
	})
	api.TestPower()
}