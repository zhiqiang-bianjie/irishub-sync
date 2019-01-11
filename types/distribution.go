package types

import (
	"encoding/json"
	"github.com/irisnet/explorer/backend/types"
)

type WithdrawDelegatorRewardsAllMsg struct {
	DelegatorAddr string `json:"delegator_addr"`
}

func NewWithdrawDelegatorRewardsAllMsg(msg MsgWithdrawDelegatorRewardsAll) WithdrawDelegatorRewardsAllMsg {
	return WithdrawDelegatorRewardsAllMsg{
		DelegatorAddr: msg.DelegatorAddr.String(),
	}
}

func (s WithdrawDelegatorRewardsAllMsg) Type() string {
	return types.TxTypeWithdrawDelegatorRewardsAll
}

func (s WithdrawDelegatorRewardsAllMsg) String() string {
	str, _ := json.Marshal(s)
	return string(str)
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

func (s WithdrawDelegatorRewardMsg) Type() string {
	return TxTypeWithdrawDelegatorReward
}

func (s WithdrawDelegatorRewardMsg) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}

type WithdrawValidatorRewardsAllMsg struct {
	ValidatorAddr string `json:"validator_addr"`
}

func NewWithdrawValidatorRewardsAllMsg(msg MsgWithdrawValidatorRewardsAll) WithdrawValidatorRewardsAllMsg {
	return WithdrawValidatorRewardsAllMsg{
		ValidatorAddr: msg.ValidatorAddr.String(),
	}
}

func (s WithdrawValidatorRewardsAllMsg) Type() string {
	return types.TxTypeWithdrawValidatorRewardsAll
}

func (s WithdrawValidatorRewardsAllMsg) String() string {
	str, _ := json.Marshal(s)
	return string(str)
}
