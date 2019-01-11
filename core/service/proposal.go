package service

import (
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/rpc"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	"github.com/irisnet/irishub-sync/util"
)

func handleProposal(docTx document.CommonTx) {
	switch docTx.Type {
	case types.TxTypeSubmitProposal:
		if proposal, err := rpc.GetProposal(docTx.ProposalId); err == nil {
			propo := ConvertProp(proposal)
			store.SaveOrUpdate(propo)
		}
	case types.TxTypeDeposit:
		if proposal, err := document.QueryProposal(docTx.ProposalId); err == nil {
			prop, err := rpc.GetProposal(docTx.ProposalId)
			if err != nil {
				logger.Error("ProposalId not existed", logger.Uint64("ProposalId", docTx.ProposalId))
				return
			}
			proposal.TotalDeposit = store.ParseCoins(prop.TotalDeposit)
			proposal.Status = prop.Status
			proposal.VotingStartTime = prop.VotingStartTime
			proposal.VotingEndTime = prop.VotingEndTime
			store.SaveOrUpdate(proposal)
		}
	case types.TxTypeVote:
		//失败的投票不计入统计
		if docTx.Status == document.TxStatusFail {
			return
		}
		if proposal, err := document.QueryProposal(docTx.ProposalId); err == nil {
			var msg types.Vote
			util.Map2Struct(docTx.Msg, &msg)
			vote := document.PVote{
				Voter:  msg.Voter,
				Option: msg.Option,
				Time:   docTx.Time,
			}
			var i int
			var hasVote = false
			for i = range proposal.Votes {
				if proposal.Votes[i].Voter == vote.Voter {
					hasVote = true
					break
				}
			}
			if hasVote {
				proposal.Votes[i] = vote
			} else {
				proposal.Votes = append(proposal.Votes, vote)
			}
			store.SaveOrUpdate(proposal)
		}
	}
}

func ConvertProp(prop rpc.Proposal) document.Proposal {
	var votes []document.PVote
	for _, v := range prop.Votes {
		votes = append(votes, document.PVote{
			Voter:  v.Voter,
			Option: v.Option,
			Time:   v.Time,
		})
	}
	return document.Proposal{
		ProposalId:      prop.ProposalId,
		Title:           prop.Title,
		Type:            prop.Type,
		Description:     prop.Description,
		Status:          prop.Status,
		SubmitTime:      prop.SubmitTime,
		DepositEndTime:  prop.DepositEndTime,
		VotingStartTime: prop.VotingStartTime,
		VotingEndTime:   prop.VotingEndTime,
		TotalDeposit:    store.ParseCoins(prop.TotalDeposit),
		Votes:           votes,
	}
}
