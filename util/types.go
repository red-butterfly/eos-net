package util

import (


	"github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
)


type ProducerInfo struct{
	Owner 					eos.AccountName		`json:"owner"`
	TotalVotes				string 		`json:"total_votes"`
	ProducerKey 			ecc.PublicKey	`json:"producer_key"`
	Url 					string 			`json:"url"`
	UnpaidBlocks			uint32			`json:"unpaid_blocks"`
	LastClaimTime			uint64		`json:"last_claim_time"`
	Location				uint16			`json:"location"`
	TimeBecameActive		uint32			`json:"time_became_active"`
	LastProducedBlockTime	uint32		`json:"last_produced_block_time"`
}

type ProducersResp struct {
	// TODO: fill this in !
	Producer	[]ProducerInfo			`json:"rows"`
	More 		bool 			`json:"more"`
}
