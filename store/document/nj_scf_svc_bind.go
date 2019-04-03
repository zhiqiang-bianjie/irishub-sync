package document

import (
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2/bson"
)

const CollectionNmSvcBind = "nj_scf_svc_bind"

type SvcBind struct {
	ID          string      `bson:"_id"`
	Hash        string      `bson:"hash"`
	DefName     string      `bson:"def_name"`
	DefChainID  string      `bson:"def_chain_id"`
	BindChainID string      `bson:"bind_chain_id"`
	Provider    string      `bson:"provider"`
	BindingType string      `bson:"binding_type"`
	Deposit     store.Coins `bson:"deposit"`
	Prices      store.Coins `bson:"price"`
	Level       Level       `bson:"level"`
	Available   bool        `bson:"available"`
}

type Level struct {
	AvgRspTime int64 `bson:"avg_rsp_time"`
	UsableTime int64 `bson:"usable_time"`
}

func (m SvcBind) Name() string {
	return CollectionNmSvcBind
}

func (m SvcBind) PkKvPair() map[string]interface{} {
	return bson.M{"hash": m.Hash}
}
