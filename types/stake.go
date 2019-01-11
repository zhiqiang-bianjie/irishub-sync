package types

type BeginRedelegate struct {
	DelegatorAddr    string `json:"delegator_addr"`
	ValidatorSrcAddr string `json:"validator_src_addr"`
	ValidatorDstAddr string `json:"validator_dst_addr"`
	SharesAmount     string `json:"shares_amount"`
}

func NewBeginRedelegate(msg MsgBeginRedelegate) BeginRedelegate {
	shares := msg.SharesAmount.String()
	return BeginRedelegate{
		DelegatorAddr:    msg.DelegatorAddr.String(),
		ValidatorSrcAddr: msg.ValidatorSrcAddr.String(),
		ValidatorDstAddr: msg.ValidatorDstAddr.String(),
		SharesAmount:     shares,
	}
}

type CreateValidator struct {
	Description
	Commission    Commission
	DelegatorAddr string `json:"delegator_address"`
	ValidatorAddr string `json:"validator_address"`
	PubKey        string `json:"pubkey"`
	Delegation    string `json:"delegation"`
}

type Description struct {
	Moniker  string `json:"moniker"`  // name
	Identity string `json:"identity"` // optional identity signature (ex. UPort or Keybase)
	Website  string `json:"website"`  // optional website link
	Details  string `json:"details"`  // optional details
}

type Commission struct {
	Rate          string `json:"rate"`            // the commission rate charged to delegators
	MaxRate       string `json:"max_rate"`        // maximum commission rate which validator can ever charge
	MaxChangeRate string `json:"max_change_rate"` // maximum daily increase of the validator commission
}

func NewCreateValidator(msg MsgCreateValidator) CreateValidator {
	pubkey, _ := Bech32ifyValPub(msg.PubKey)
	return CreateValidator{
		Description: Description{
			Moniker:  msg.Description.Moniker,
			Identity: msg.Description.Identity,
			Website:  msg.Description.Website,
			Details:  msg.Description.Details,
		},
		Commission: Commission{
			Rate:          msg.Commission.Rate.String(),
			MaxRate:       msg.Commission.MaxRate.String(),
			MaxChangeRate: msg.Commission.MaxChangeRate.String(),
		},
		DelegatorAddr: msg.DelegatorAddr.String(),
		ValidatorAddr: msg.ValidatorAddr.String(),
		PubKey:        pubkey,
		Delegation:    msg.Delegation.String(),
	}
}

type EditValidator struct {
	Description
	ValidatorAddr  string `json:"address"`
	CommissionRate string `json:"commission_rate"`
}

func NewEditValidator(msg MsgEditValidator) EditValidator {
	return EditValidator{
		Description: Description{
			Moniker:  msg.Description.Moniker,
			Identity: msg.Description.Identity,
			Website:  msg.Description.Website,
			Details:  msg.Description.Details,
		},
		ValidatorAddr:  msg.ValidatorAddr.String(),
		CommissionRate: msg.CommissionRate.String(),
	}
}
