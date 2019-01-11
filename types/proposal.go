package types

type SubmitProposal struct {
	Title          string `json:"title"`          //  Title of the proposal
	Description    string `json:"description"`    //  Description of the proposal
	Proposer       string `json:"proposer"`       //  Address of the proposer
	InitialDeposit Coins  `json:"initialDeposit"` //  Initial deposit paid by sender. Must be strictly positive.
	ProposalType   string `json:"proposalType"`   //  Initial deposit paid by sender. Must be strictly positive.
	Params         Params `json:"params"`
}

type Param struct {
	Subspace string `json:"subspace"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}

type Params []Param

func NewSubmitProposal(msg MsgSubmitProposal) SubmitProposal {
	var params Params
	for _, p := range msg.Params {
		params = append(params, Param{
			Subspace: p.Subspace,
			Key:      p.Key,
			Value:    p.Value,
		})
	}
	return SubmitProposal{
		Title:          msg.Title,
		Description:    msg.Description,
		ProposalType:   msg.ProposalType.String(),
		Proposer:       msg.Proposer.String(),
		InitialDeposit: ParseCoins(msg.InitialDeposit.String()),
		Params:         params,
	}
}

type Vote struct {
	ProposalID uint64 `json:"proposal_id"`
	Voter      string `json:"voter"`
	Option     string `json:"option"`
}

func NewVote(v MsgVote) Vote {
	return Vote{
		ProposalID: v.ProposalID,
		Voter:      v.Voter.String(),
		Option:     v.Option.String(),
	}
}

type Deposit struct {
	ProposalID uint64 `json:"proposal_id"` // ID of the proposal
	Depositer  string `json:"depositer"`   // Address of the depositer
	Amount     Coins  `json:"amount"`      // Coins to add to the proposal's deposit
}

func NewDeposit(deposit MsgDeposit) Deposit {
	return Deposit{
		ProposalID: deposit.ProposalID,
		Depositer:  deposit.Depositor.String(),
		Amount:     ParseCoins(deposit.Amount.String()),
	}
}

type SubmitSoftwareUpgradeProposal struct {
	SubmitProposal
	Version      uint64 `json:"version"`
	Software     string `json:"software"`
	SwitchHeight uint64 `json:"switch_height"`
}

func NewSubmitSoftwareUpgradeProposal(msg MsgSubmitSoftwareUpgradeProposal) SubmitSoftwareUpgradeProposal {
	return SubmitSoftwareUpgradeProposal{
		SubmitProposal: NewSubmitProposal(msg.MsgSubmitProposal),
		Version:        msg.Version,
		Software:       msg.Software,
		SwitchHeight:   msg.SwitchHeight,
	}
}
