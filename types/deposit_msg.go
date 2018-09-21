package types

import (
	"encoding/json"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/util/constant"
)

type Deposit struct {
	ProposalID int64       `json:"proposal_id"` // ID of the proposal
	Depositer  string      `json:"depositer"`   // Address of the depositer
	Amount     store.Coins `json:"amount"`      // Coins to add to the proposal's deposit
}

func NewDeposit(deposit MsgDeposit) Deposit {
	return Deposit{
		ProposalID: deposit.ProposalID,
		Depositer:  deposit.Depositer.String(),
		Amount:     BuildCoins(deposit.Amount),
	}
}

func (s Deposit) Type() string {
	return constant.TxTypeDeposit
}

func (s Deposit) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

func UnmarshalDeposit(str string) (deposit Deposit) {
	json.Unmarshal([]byte(str), &deposit)
	return
}