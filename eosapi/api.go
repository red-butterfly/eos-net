package eosapi

import (
	"fmt"
	"log"
	"strings"
	"time"
	"encoding/json"

	"../util"

	"github.com/eoscanada/eos-go/ecc"
	"github.com/eoscanada/eos-go/token"
	eos_go "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/system"
)

type Eos struct {
	Api *eos_go.API
}

var SYSSymbol = eos_go.Symbol{4,"SIS"}

func NewSYSAsset(amount int64) eos_go.Asset {
	return eos_go.Asset{Amount: amount, Symbol: SYSSymbol}
}

func GetEOS(PriKey []string) *Eos {
	fmt.Println("connect to ", util.NetInfo.Hosturl)
	geteos := &Eos{
		Api: eos_go.New(util.NetInfo.Hosturl),
	}

	geteos.Api.Debug = false
	eos_go.Debug = false

	keyBag := eos_go.NewKeyBag()
	for _, key := range PriKey {
		if err := keyBag.Add(key); err != nil {
			log.Fatalln("Couldn't load private key:", err)
		}
	}

	geteos.Api.SetSigner(keyBag)

	return geteos
}

func (eos *Eos) setSystemCode(account string, contract string) error {
	// set contract eosio.token
	setCodeTx, err := system.NewSetCodeTx(
		eos_go.AN(account),
		strings.Join([]string{".", "contracts", contract, contract+".wasm"}, "/"),
		strings.Join([]string{".", "contracts", contract, contract+".abi"}, "/"),
	)

	if err != nil {
		log.Fatal(err)
		return err
	}

	resp, err := eos.Api.SignPushTransaction(setCodeTx, &eos_go.TxOptions{})
	if err != nil {
		fmt.Println("ERROR calling SetCode:", err)
		return err
	} else {
		fmt.Printf("RESP the tranID (set contract %s): %s\n",contract, resp.TransactionID)
	}
	return nil
}

func (eos *Eos) SystemCreateAccount(create string, account string, pubkey string, ram uint32, net, cpu int64, ) error {
	if ram == 0 {
		ram = 8192
	}
	if net == 0 {
		net = 1000
	}
	if cpu == 0 {
		cpu = 1000
	}
	actionResp, err := eos.Api.SignPushActions(
		system.NewNewAccount(eos_go.AN(create), eos_go.AN(account), ecc.MustNewPublicKey(pubkey)),
		system.NewBuyRAMBytes(
			eos_go.AN(create),
			eos_go.AN(account),
			ram,
		),
		system.NewDelegateBW(
			eos_go.AN(create),
			eos_go.AN(account),
			NewSYSAsset(cpu),
			NewSYSAsset(net),
			true,
		),
	)
	if err != nil {
		fmt.Println("ERROR calling :", err)
		return err
	} else {
		fmt.Printf("RESP the tranID (system create account %s): %s\n",account, actionResp.TransactionID)
	}
	return nil
}

func (eos *Eos) SysTransfer(from string, to string, amount int64) error {
	actionResp, err := eos.Api.SignPushActions(
		token.NewTransfer(eos_go.AN(from), eos_go.AN(to), NewSYSAsset(amount), "transfer SYS"),
	)
	if err != nil {
		fmt.Println("ERROR calling :", err)
		return err
	} else {
		fmt.Printf("RESP the tranID (transfer SYS: %s-(%f)->%s): %s \n",from, float64(amount/10000), to, actionResp.TransactionID)
	}
	return nil
}

func (eos *Eos) InitEosioNode() {
	// create eosio.token

	for _,account := range util.ActInfo.Systemaccount {
		actionResp, err := eos.Api.SignPushActions(
			system.NewNewAccount(eos_go.AN("eosio"), eos_go.AN(account), ecc.MustNewPublicKey(util.ActInfo.Sys_publickey)),
		)
		if err != nil {
			fmt.Println("ERROR calling :", err)
		} else {
			fmt.Printf("RESP the tranID (create account %s): %s\n",account, actionResp.TransactionID)
		}

		time.Sleep(time.Millisecond  * 100)
	}

	// set contract eosio.token
	if err := eos.setSystemCode("eosio.token", "eosio.token"); err != nil {
		return
	}

	time.Sleep(time.Millisecond  * 100)

	// set contract eosio.msig
	if err := eos.setSystemCode("eosio.msig", "eosio.msig"); err != nil {
		return
	}

	time.Sleep(time.Millisecond  * 100)

	// create SYS
	actionResp, err := eos.Api.SignPushActions(
		token.NewCreate(eos_go.AN("eosio"), NewSYSAsset(10000000000*10000)),
	)
	if err != nil {
		fmt.Println("ERROR calling :", err)
	} else {
		fmt.Println("RESP the tranID (create SYS):", actionResp.TransactionID)
	}
	time.Sleep(time.Millisecond  * 100)

	// issue SYS to eosio
	actionResp, err = eos.Api.SignPushActions(
		token.NewIssue(eos_go.AN("eosio"), NewSYSAsset(1000000000*10000), "issue SYS"),
	)
	if err != nil {
		fmt.Println("ERROR calling :", err)
	} else {
		fmt.Println("RESP the tranID (issue SYS):", actionResp.TransactionID)
	}
	time.Sleep(time.Millisecond  * 100)

	// set contract eosio.system
	for true {
		if err := eos.setSystemCode("eosio", "eosio.system"); err == nil {
			break
		}
	}

	// set priv to eosio.msig
	actionResp, err = eos.Api.SignPushActions(
		system.NewSetPriv(eos_go.AN("eosio.msig")),
	)
	if err != nil {
		fmt.Println("ERROR calling :", err)
	} else {
		fmt.Printf("RESP the tranID (setpriv to eosio.msig): %s\n", actionResp.TransactionID)
	}
}

func (eos *Eos) CreateVoters()  {
	var allSYSInVoter int64 = 300000000
	var VoterOut int64 = 100000
	err := eos.SystemCreateAccount("eosio", util.ActInfo.Voteaccount[0], util.ActInfo.Voteaccount[2], 0,100*10000,500*10000)
	if err != nil {
		return
	}
	time.Sleep(time.Millisecond  * 100)

	eos.Api.Signer.ImportPrivateKey(util.ActInfo.Voteaccount[1])
	if err = eos.SysTransfer("eosio", util.ActInfo.Voteaccount[0], allSYSInVoter*10000 + VoterOut*10000); err != nil {
		return
	}

	everyamount := (allSYSInVoter*10000)/int64(len(util.ActInfo.Voters))
	unstaked := everyamount/100
	staked := (everyamount-unstaked)/2
	for _,voter := range util.ActInfo.Voters {
		time.Sleep(time.Millisecond  * 1000)
		err := eos.SystemCreateAccount(util.ActInfo.Voteaccount[0], voter[0], voter[2], 0,staked,staked)
		if err != nil {
			return
		}

		eos.Api.Signer.ImportPrivateKey(voter[1])
		if err = eos.SysTransfer(util.ActInfo.Voteaccount[0], voter[0], unstaked); err != nil {
			return
		}
	}

}

func (eos *Eos) RegProducer()  {
	var allSYSInProd int64 = 10000000
	var prodOut int64 = 100000

	err := eos.SystemCreateAccount("eosio", util.ActInfo.Produceraccount[0], util.ActInfo.Produceraccount[2], 0,100*10000,500*10000)
	if err != nil {
		return
	}
	time.Sleep(time.Millisecond  * 100)

	eos.Api.Signer.ImportPrivateKey(util.ActInfo.Produceraccount[1])
	if err = eos.SysTransfer("eosio", util.ActInfo.Produceraccount[0], allSYSInProd*10000 + prodOut*10000); err != nil {
		return
	}

	everyamount := (allSYSInProd*10000)/int64(len(util.ActInfo.Producers))
	unstaked := everyamount/100
	staked := (everyamount-unstaked)/2
	for _,proder := range util.ActInfo.Producers {
		time.Sleep(time.Millisecond  * 1000)
		err := eos.SystemCreateAccount(util.ActInfo.Produceraccount[0], proder[0], proder[2], 0,staked,staked)
		if err != nil {
			return
		}

		eos.Api.Signer.ImportPrivateKey(proder[1])
		if err = eos.SysTransfer(util.ActInfo.Produceraccount[0], proder[0], unstaked); err != nil {
			return
		}

		actionResp, err := eos.Api.SignPushActions(
			system.NewRegProducer(eos_go.AN(proder[0]), ecc.MustNewPublicKey(proder[2]), "http://"+proder[0]+".com"),
		)
		if err != nil {
			fmt.Println("ERROR calling :", err)
			return
		} else {
			fmt.Printf("RESP the tranID (regproducer %s): %s \n", proder[0], actionResp.TransactionID)
		}
	}

}

func (eos *Eos) ListProducers()  {
	prodResp, err := eos.Api.GetTableRows(
		eos_go.GetTableRowsRequest{
			Code: "eosio",
			Scope: "eosio",
			Table: "producers",
			JSON: true,
			Limit: 10,
		},
	)

	data, err := json.MarshalIndent(prodResp, "", "  ")
	if err != nil {
		fmt.Printf("Error: json conversion , %s\n", err.Error())
		return
	}
	var undata util.ProducersResp
	if err := json.Unmarshal(data, &undata); err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	for _, prod := range undata.Producer {
		fmt.Printf("prod_name: %s\n", prod.Owner)
		fmt.Printf("  pub_key: %s\n", prod.ProducerKey)
		fmt.Printf("  url: %s\n", prod.Url)
	}
}

func (eos *Eos) VoteProducer() {
	prodlen := len(util.ActInfo.Producers)
	prods := make([]eos_go.AccountName, prodlen)
	for i,proder := range util.ActInfo.Producers {
		prods[i] = eos_go.AN(proder[0])
	}

	for _,voter := range util.ActInfo.Voters {
		eos.Api.Signer.ImportPrivateKey(voter[1])
		actionResp, err := eos.Api.SignPushActions(
			system.NewVoteProducer(eos_go.AN(voter[0]), eos_go.AN(""),prods...),
		)
		if err != nil {
			fmt.Println("ERROR calling :", err)
			return
		} else {
			fmt.Printf("RESP the tranID (%s %s): %s \n", voter[0], util.ActInfo.Producers[0][0], actionResp.TransactionID)
		}
	}
}

func (eos *Eos) updateAuth(seter, change_perm, set_perm, parent, setAccount, setPerm string)  {
	actionResp, err := eos.Api.SignPushActions(
		system.NewUpdateAuth(eos_go.AN(seter),
			eos_go.PN(change_perm),
			eos_go.PN(parent),
			eos_go.Authority{
				1,
				[]eos_go.KeyWeight{},
				[]eos_go.PermissionLevelWeight{
					{
						eos_go.PermissionLevel{
							eos_go.AN(setAccount),
							eos_go.PN(setPerm),
						},
						1,
					},
				},
				[]eos_go.WaitWeight{},
			},
			eos_go.PN(set_perm),
		),
	)
	if err != nil {
		fmt.Println("ERROR calling :", err)
		return
	} else {
		fmt.Printf("RESP the tranID (resign %s(%s) to %s(%s)): %s \n", seter, set_perm,setAccount, setPerm,actionResp.TransactionID)
	}
}

func (eos *Eos) StepResign() {
	for _,proder := range util.ActInfo.Producers {
		eos.Api.Signer.ImportPrivateKey(proder[1])
	}
	eos.updateAuth("eosio", "owner", "owner", "", "eosio.prods","active")
	eos.updateAuth("eosio", "active", "owner", "owner", "eosio.prods","active")

	for _,account := range util.ActInfo.Systemaccount {
		eos.updateAuth(account, "owner", "owner", "", "eosio","active")
		eos.updateAuth(account, "active", "owner", "owner", "eosio","active")
	}
}

func (eos *Eos) TestPower() {
	//err := eos.SystemCreateAccount("yan", "testman12342", "EOS8fCPpS2cMKekEQkevHJuAxm7XaNwgZjoMtAgQh9jEWpuCpSGqm", 0,1,1)
	//
	//if err != nil {
	//	return
	//}
	//if err := eos.SysTransfer("yan", "testman12344",3000*10000 ); err != nil {
	//	return
	//}
	//
	//for a := 0; a < 1; a++ {
	//	if err := eos.SysTransfer("testman12344", "yan",1*10000+int64(a) ); err != nil {
	//		return
	//	}
	//	time.Sleep(time.Millisecond  * 100)
	//}


}