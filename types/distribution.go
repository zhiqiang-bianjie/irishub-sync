package types

type WithdrawDelegatorRewardsAllMsg struct {
	DelegatorAddr string `json:"delegator_addr"`
}

func NewWithdrawDelegatorRewardsAllMsg(msg MsgWithdrawDelegatorRewardsAll) WithdrawDelegatorRewardsAllMsg {
	return WithdrawDelegatorRewardsAllMsg{
		DelegatorAddr: msg.DelegatorAddr.String(),
	}
}

type WithdrawDelegatorRewardMsg struct {
	DelegatorAddr string `json:"delegator_addr"`
	ValidatorAddr string `json:"validator_addr"`
}

func NewWithdrawDelegatorRewardMsg(msg MsgWithdrawDelegatorReward) WithdrawDelegatorRewardMsg {
	return WithdrawDelegatorRewardMsg{
		DelegatorAddr: msg.DelegatorAddr.String(),
		ValidatorAddr: msg.ValidatorAddr.String(),
	}
}

type WithdrawValidatorRewardsAllMsg struct {
	ValidatorAddr string `json:"validator_addr"`
}

func NewWithdrawValidatorRewardsAllMsg(msg MsgWithdrawValidatorRewardsAll) WithdrawValidatorRewardsAllMsg {
	return WithdrawValidatorRewardsAllMsg{
		ValidatorAddr: msg.ValidatorAddr.String(),
	}
}
