package document

import (
	"github.com/irisnet/irishub-sync/store"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const CollectionNmSvcInvocation = "nj_scf_svc_invocation"

type SvcInvocation struct {
	Hash        string    `bson:"hash"`
	ReqId       string    `bson:"req_id"`
	DefChainID  string    `bson:"def_chain_id"`
	DefName     string    `bson:"def_name"`
	BindChainID string    `bson:"bind_chain_id"`
	ReqChainID  string    `bson:"req_chain_id"`
	MethodID    int16     `bson:"method_id"`
	TxType      string    `bson:"tx_type"`
	Consumer    string    `bson:"consumer"`
	Provider    string    `bson:"provider"`
	Height      int64     `bson:"height"`
	Data        string    `bson:"data"`
	Time        time.Time `bson:"time"`
}

func (m SvcInvocation) Name() string {
	return "nj_scf_svc_invocation"
}

func (m SvcInvocation) PkKvPair() map[string]interface{} {
	return bson.M{"hash": m.Hash}
}

func (m SvcInvocation) QueryByReqId() (result []SvcInvocation, err error) {
	q := bson.M{}
	q["req_id"] = m.ReqId
	fn := func(c *mgo.Collection) error {
		return c.Find(q).All(&result)
	}

	err = store.ExecCollection(m.Name(), fn)

	if err != nil {
		return nil, err
	}
	return
}
