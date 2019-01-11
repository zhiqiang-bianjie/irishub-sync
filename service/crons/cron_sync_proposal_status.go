package crons

import (
	conf "github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/rpc"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
)

func syncProposalStatus() {
	var status = []string{types.StatusDepositPeriod, types.StatusVotingPeriod}
	if proposals, err := document.QueryByStatus(status); err == nil {
		for _, proposal := range proposals {
			propo, err := rpc.GetProposal(proposal.ProposalId)
			if err != nil {
				store.Delete(proposal)
				return
			}
			if propo.Status != proposal.Status {
				propo.SubmitTime = proposal.SubmitTime
				propo.Votes = proposal.Votes
				store.SaveOrUpdate(propo)
			}
		}
	}
}

func MakeSyncProposalStatusTask() Task {
	return NewLockTaskFromEnv(conf.SyncProposalStatus, "sync_proposal_status_lock", func() {
		logger.Info("start cron", logger.String("cronNm", "syncProposalStatus"))
		syncProposalStatus()
	})
}
