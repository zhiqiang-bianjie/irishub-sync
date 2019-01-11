package types

import (
	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/client/utils"
	"github.com/irisnet/irishub/codec"
	"github.com/irisnet/irishub/modules/auth"
	"github.com/irisnet/irishub/modules/bank"
	"github.com/irisnet/irishub/modules/distribution"
	dtags "github.com/irisnet/irishub/modules/distribution/tags"
	"github.com/irisnet/irishub/modules/gov"
	"github.com/irisnet/irishub/modules/gov/tags"
	"github.com/irisnet/irishub/modules/slashing"
	"github.com/irisnet/irishub/modules/stake"
	stags "github.com/irisnet/irishub/modules/stake/tags"
	staketypes "github.com/irisnet/irishub/modules/stake/types"
	"github.com/irisnet/irishub/types"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tm "github.com/tendermint/tendermint/types"
)

type (
	MsgSend = bank.MsgSend

	MsgCreateValidator             = stake.MsgCreateValidator
	MsgEditValidator               = stake.MsgEditValidator
	MsgDelegate                    = stake.MsgDelegate
	MsgBeginUnbonding              = stake.MsgBeginUnbonding
	MsgBeginRedelegate             = stake.MsgBeginRedelegate
	MsgUnjail                      = slashing.MsgUnjail
	MsgWithdrawDelegatorReward     = distribution.MsgWithdrawDelegatorReward
	MsgWithdrawDelegatorRewardsAll = distribution.MsgWithdrawDelegatorRewardsAll
	MsgWithdrawValidatorRewardsAll = distribution.MsgWithdrawValidatorRewardsAll
	StakeValidator                 = stake.Validator
	Delegation                     = stake.Delegation
	UnbondingDelegation            = stake.UnbondingDelegation

	MsgDeposit                       = gov.MsgDeposit
	MsgSubmitProposal                = gov.MsgSubmitProposal
	MsgSubmitSoftwareUpgradeProposal = gov.MsgSubmitSoftwareUpgradeProposal
	MsgVote                          = gov.MsgVote
	Proposal                         = gov.Proposal
	SdkVote                          = gov.Vote

	ResponseDeliverTx = abci.ResponseDeliverTx

	StdTx      = auth.StdTx
	StdFee     = auth.StdFee
	SdkCoins   = types.Coins
	KVPair     = types.KVPair
	AccAddress = types.AccAddress
	ValAddress = types.ValAddress
	Dec        = types.Dec
	Validator  = tm.Validator
	Tx         = tm.Tx
	Block      = tm.Block
	BlockMeta  = tm.BlockMeta
	HexBytes   = cmn.HexBytes
	TmKVPair   = cmn.KVPair

	ABCIQueryOptions = rpcclient.ABCIQueryOptions
	Client           = rpcclient.Client
	HTTP             = rpcclient.HTTP
	ResultStatus     = ctypes.ResultStatus
)

var (
	ValidatorsKey        = stake.ValidatorsKey
	GetValidatorKey      = stake.GetValidatorKey
	GetDelegationKey     = stake.GetDelegationKey
	GetDelegationsKey    = stake.GetDelegationsKey
	GetUBDKey            = stake.GetUBDKey
	GetUBDsKey           = stake.GetUBDsKey
	ValAddressFromBech32 = types.ValAddressFromBech32

	UnmarshalValidator      = staketypes.UnmarshalValidator
	MustUnmarshalValidator  = staketypes.MustUnmarshalValidator
	UnmarshalDelegation     = staketypes.UnmarshalDelegation
	MustUnmarshalDelegation = staketypes.MustUnmarshalDelegation
	MustUnmarshalUBD        = staketypes.MustUnmarshalUBD

	Bech32ifyValPub      = types.Bech32ifyValPub
	RegisterCodec        = types.RegisterCodec
	AccAddressFromBech32 = types.AccAddressFromBech32
	BondStatusToString   = types.BondStatusToString

	NewDecFromStr = types.NewDecFromStr

	AddressStoreKey   = auth.AddressStoreKey
	GetAccountDecoder = utils.GetAccountDecoder

	KeyProposal      = gov.KeyProposal
	KeyVotesSubspace = gov.KeyVotesSubspace

	NewHTTP = rpcclient.NewHTTP

	//tags
	TagGovProposalID                   = tags.ProposalID
	TagDistributionReward              = dtags.Reward
	TagStakeActionCompleteRedelegation = stags.ActionCompleteRedelegation
	TagStakeDelegator                  = stags.Delegator
	TagStakeSrcValidator               = stags.SrcValidator
	TagAction                          = types.TagAction

	cdc *codec.Codec
)

// 初始化账户地址前缀
func init() {
	//TODO
	//config := types.GetConfig()
	//config.SetBech32PrefixForAccount(server.Bech32.PrefixAccAddr, server.Bech32.PrefixAccPub)
	//config.SetBech32PrefixForValidator(server.Bech32.PrefixValAddr, server.Bech32.PrefixValPub)
	//config.SetBech32PrefixForConsensusNode(server.Bech32.PrefixAccAddr, server.Bech32.PrefixConsPub)
	//config.Seal()

	cdc = app.MakeLatestCodec()
}

func GetCodec() *codec.Codec {
	return cdc
}
