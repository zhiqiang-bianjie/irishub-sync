package document

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

const CollectionNmSvcInvocation = "nj_scf_svc_invocation"

type SvcInvocation struct {
	Hash     string    `bson:"hash"`
	ReqId    string    `bson:"req_id"`
	TxType   string    `bson:"tx_type"`
	Consumer string    `bson:"consumer"`
	Provider string    `bson:"provider"`
	Height   int64     `bson:"height"`
	Data     string    `bson:"data"`
	Time     time.Time `bson:"time"`
}

func (m SvcInvocation) Name() string {
	return "nj_scf_svc_invocation"
}

func (m SvcInvocation) PkKvPair() map[string]interface{} {
	return bson.M{"hash": m.Hash}
}
