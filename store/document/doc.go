package document

import (
	"github.com/irisnet/irishub-sync/store"
)

func init() {
	store.RegisterDocs(new(Account))
	store.RegisterDocs(new(StakeTx))
	store.RegisterDocs(new(StakeTxDeclareCandidacy))
	store.RegisterDocs(new(StakeTxEditCandidacy))
	store.RegisterDocs(new(Candidate))
	store.RegisterDocs(new(Delegator))
	store.RegisterDocs(new(Block))
	store.RegisterDocs(new(SyncTask))
}