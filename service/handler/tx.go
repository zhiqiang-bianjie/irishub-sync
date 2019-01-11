package handler

import (
	"encoding/hex"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/rpc"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	itypes "github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util"
	"strconv"
	"strings"
	"sync"
)

// save Tx document into collection
func SaveTx(docTx document.CommonTx, mutex sync.Mutex) {
	var (
		methodName = "SaveTx"
	)
	logger.Debug("Start", logger.String("method", methodName))

	// save common docTx document
	saveCommonTx := func(commonTx document.CommonTx) {
		//save tx
		err := store.Save(commonTx)
		if err != nil {
			logger.Error("Save commonTx failed", logger.Any("Tx", commonTx), logger.String("err", err.Error()))
		}
		handleProposal(commonTx)
	}

	saveCommonTx(docTx)
	logger.Debug("End", logger.String("method", methodName))
}

func ParseTx(txBytes itypes.Tx, block *itypes.Block) document.CommonTx {
	var (
		authTx     itypes.StdTx
		methodName = "ParseTx"
		docTx      document.CommonTx
		gasPrice   float64
		actualFee  store.ActualFee
	)

	cdc := itypes.GetCodec()

	err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &authTx)
	if err != nil {
		logger.Error(err.Error())
		return docTx
	}

	height := block.Height
	time := block.Time
	txHash := GetTxHash(txBytes.Hash())
	fee := buildFee(authTx.Fee)
	memo := authTx.Memo

	// get tx status, gasUsed, gasPrice and actualFee from tx result
	status, result, err := queryTxResult(txBytes.Hash())
	if err != nil {
		logger.Error("get txResult err", logger.String("method", methodName), logger.String("err", err.Error()))
	}
	log := result.Log
	gasUsed := result.GasUsed
	if len(fee.Amount) > 0 {
		gasPrice = fee.Amount[0].Amount / float64(fee.Gas)
		actualFee = store.ActualFee{
			Denom:  fee.Amount[0].Denom,
			Amount: float64(gasUsed) * gasPrice,
		}
	} else {
		gasPrice = 0
		actualFee = store.ActualFee{}
	}

	msgs := authTx.GetMsgs()
	if len(msgs) <= 0 {
		logger.Error("can't get msgs", logger.String("method", methodName))
		return docTx
	}
	msg := msgs[0]

	docTx = document.CommonTx{
		Height:    height,
		Time:      time,
		TxHash:    txHash,
		Fee:       fee,
		Memo:      memo,
		Status:    status,
		Code:      result.Code,
		Log:       log,
		GasUsed:   gasUsed,
		GasPrice:  gasPrice,
		ActualFee: actualFee,
		Tags:      parseTxTags(result),
		Msg:       util.Struct2Map(msg),
	}

	switch msg.(type) {
	case itypes.MsgSend:
		msg := msg.(itypes.MsgSend)

		docTx.From = msg.Inputs[0].Address.String()
		docTx.To = msg.Outputs[0].Address.String()
		docTx.Amount = store.ParseCoins(msg.Inputs[0].Coins.String())
		docTx.Type = itypes.TxTypeTransfer
		return docTx
	case itypes.MsgCreateValidator:
		msg := msg.(itypes.MsgCreateValidator)

		docTx.From = msg.DelegatorAddr.String()
		docTx.To = msg.ValidatorAddr.String()
		docTx.Amount = []store.Coin{store.ParseCoin(msg.Delegation.String())}
		docTx.Type = itypes.TxTypeStakeCreateValidator

		// struct of createValidator
		valDes := document.ValDescription{
			Moniker:  msg.Moniker,
			Identity: msg.Identity,
			Website:  msg.Website,
			Details:  msg.Details,
		}
		pubKey, err := itypes.Bech32ifyValPub(msg.PubKey)
		if err != nil {
			logger.Error("Can't get pubKey", logger.String("txHash", txHash))
			pubKey = ""
		}
		docTx.StakeCreateValidator = document.StakeCreateValidator{
			PubKey:      pubKey,
			Description: valDes,
		}

		docTx.Msg = util.Struct2Map(itypes.NewCreateValidator(msg))

		return docTx
	case itypes.MsgEditValidator:
		msg := msg.(itypes.MsgEditValidator)

		docTx.From = msg.ValidatorAddr.String()
		docTx.To = ""
		docTx.Amount = []store.Coin{}
		docTx.Type = itypes.TxTypeStakeEditValidator

		// struct of editValidator
		valDes := document.ValDescription{
			Moniker:  msg.Moniker,
			Identity: msg.Identity,
			Website:  msg.Website,
			Details:  msg.Details,
		}
		docTx.StakeEditValidator = document.StakeEditValidator{
			Description: valDes,
		}
		docTx.Msg = util.Struct2Map(itypes.NewEditValidator(msg))

		return docTx
	case itypes.MsgDelegate:
		msg := msg.(itypes.MsgDelegate)

		docTx.From = msg.DelegatorAddr.String()
		docTx.To = msg.ValidatorAddr.String()
		docTx.Amount = []store.Coin{store.ParseCoin(msg.Delegation.String())}
		docTx.Type = itypes.TxTypeStakeDelegate

		return docTx
	case itypes.MsgBeginUnbonding:
		msg := msg.(itypes.MsgBeginUnbonding)

		shares := util.ParseFloat(msg.SharesAmount.String())
		docTx.From = msg.DelegatorAddr.String()
		docTx.To = msg.ValidatorAddr.String()

		coin := store.Coin{
			Amount: shares,
		}
		docTx.Amount = []store.Coin{coin}
		docTx.Type = itypes.TxTypeStakeBeginUnbonding
		return docTx
	case itypes.MsgBeginRedelegate:
		msg := msg.(itypes.MsgBeginRedelegate)

		shares := util.ParseFloat(msg.SharesAmount.String())
		docTx.From = msg.DelegatorAddr.String()
		docTx.To = msg.ValidatorDstAddr.String()
		coin := store.Coin{
			Amount: shares,
		}
		docTx.Amount = []store.Coin{coin}
		docTx.Type = itypes.TxTypeBeginRedelegate
		docTx.Msg = util.Struct2Map(itypes.NewBeginRedelegate(msg))
		return docTx
	case itypes.MsgUnjail:
		msg := msg.(itypes.MsgUnjail)

		docTx.From = msg.ValidatorAddr.String()
		docTx.Type = itypes.TxTypeUnjail

	case itypes.MsgWithdrawDelegatorReward:
		msg := msg.(itypes.MsgWithdrawDelegatorReward)

		docTx.From = msg.DelegatorAddr.String()
		docTx.To = msg.ValidatorAddr.String()
		docTx.Type = itypes.TxTypeWithdrawDelegatorReward
		docTx.Msg = util.Struct2Map(itypes.NewWithdrawDelegatorRewardMsg(msg))

		for _, tag := range result.Tags {
			key := string(tag.Key)
			if key == itypes.TagDistributionReward {
				reward := string(tag.Value)
				docTx.Amount = store.ParseCoins(reward)
				break
			}
		}
	case itypes.MsgWithdrawDelegatorRewardsAll:
		msg := msg.(itypes.MsgWithdrawDelegatorRewardsAll)

		docTx.From = msg.DelegatorAddr.String()
		docTx.Type = itypes.TxTypeWithdrawDelegatorRewardsAll
		docTx.Msg = util.Struct2Map(itypes.NewWithdrawDelegatorRewardsAllMsg(msg))
		for _, tag := range result.Tags {
			key := string(tag.Key)
			if key == itypes.TagDistributionReward {
				reward := string(tag.Value)
				docTx.Amount = store.ParseCoins(reward)
				break
			}
		}
	case itypes.MsgWithdrawValidatorRewardsAll:
		msg := msg.(itypes.MsgWithdrawValidatorRewardsAll)

		docTx.From = msg.ValidatorAddr.String()
		docTx.Type = itypes.TxTypeWithdrawValidatorRewardsAll
		docTx.Msg = util.Struct2Map(itypes.NewWithdrawValidatorRewardsAllMsg(msg))
		for _, tag := range result.Tags {
			key := string(tag.Key)
			if key == itypes.TagDistributionReward {
				reward := string(tag.Value)
				docTx.Amount = store.ParseCoins(reward)
				break
			}
		}
	case itypes.MsgSubmitProposal:
		msg := msg.(itypes.MsgSubmitProposal)

		docTx.From = msg.Proposer.String()
		docTx.To = ""
		docTx.Amount = store.ParseCoins(msg.InitialDeposit.String())
		docTx.Type = itypes.TxTypeSubmitProposal
		docTx.Msg = util.Struct2Map(itypes.NewSubmitProposal(msg))

		//query proposal_id
		for _, tag := range result.Tags {
			key := string(tag.Key)
			if key == itypes.TagGovProposalID {
				proposalId, err := strconv.ParseInt(string(tag.Value), 10, 0)
				if err == nil {
					docTx.ProposalId = uint64(proposalId)
					break
				}
			}
		}
		return docTx
	case itypes.MsgSubmitSoftwareUpgradeProposal:
		msg := msg.(itypes.MsgSubmitSoftwareUpgradeProposal)

		docTx.From = msg.Proposer.String()
		docTx.To = ""
		docTx.Amount = store.ParseCoins(msg.InitialDeposit.String())
		docTx.Type = itypes.TxTypeSubmitProposal
		docTx.Msg = util.Struct2Map(itypes.NewSubmitSoftwareUpgradeProposal(msg))

		//query proposal_id
		for _, tag := range result.Tags {
			key := string(tag.Key)
			if key == itypes.TagGovProposalID {
				proposalId, err := strconv.ParseInt(string(tag.Value), 10, 0)
				if err == nil {
					docTx.ProposalId = uint64(proposalId)
					break
				}
			}
		}
		return docTx
	case itypes.MsgDeposit:
		msg := msg.(itypes.MsgDeposit)

		docTx.From = msg.Depositor.String()
		docTx.Amount = store.ParseCoins(msg.Amount.String())
		docTx.Type = itypes.TxTypeDeposit
		docTx.Msg = util.Struct2Map(itypes.NewDeposit(msg))
		docTx.ProposalId = msg.ProposalID
		return docTx
	case itypes.MsgVote:
		msg := msg.(itypes.MsgVote)

		docTx.From = msg.Voter.String()
		docTx.Amount = []store.Coin{}
		docTx.Type = itypes.TxTypeVote
		docTx.Msg = util.Struct2Map(itypes.NewVote(msg))
		docTx.ProposalId = msg.ProposalID
		return docTx

	default:
		logger.Warn("unknown msg type")
	}

	return docTx
}

func GetTxHash(bytes []byte) string {
	return strings.ToUpper(hex.EncodeToString(bytes))
}

func buildFee(fee itypes.StdFee) store.Fee {
	return store.Fee{
		Amount: store.ParseCoins(fee.Amount.String()),
		Gas:    int64(fee.Gas),
	}
}

func queryTxResult(txHash []byte) (string, itypes.ResponseDeliverTx, error) {
	var resDeliverTx itypes.ResponseDeliverTx
	status := document.TxStatusSuccess
	res, err := rpc.GetTxResult(txHash)
	if err != nil {
		return "unknown", resDeliverTx, err
	}
	if res.Code != 0 {
		status = document.TxStatusFail
	}

	return status, res, nil
}

func parseTxTags(result itypes.ResponseDeliverTx) map[string]string {
	tags := make(map[string]string, 0)
	for _, tag := range result.Tags {
		key := string(tag.Key)
		value := string(tag.Value)
		tags[key] = value
	}
	return tags
}
