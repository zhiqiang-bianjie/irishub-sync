package document

import "gopkg.in/mgo.v2/bson"

const CollectionNmSvcDef = "nj_scf_svc_define"

type SvcDef struct {
	Hash              string `bson:"hash"`
	Code              string `bson:"code"`
	ChainId           string `bson:"chain_id"`
	Description       string `bson:"description"`
	Author            string `bson:"author"`
	AuthorDescription string `bson:"author_description"`
	IDLContent        string `bson:"idl_content"`
	Status            string `bson:"status"`
}

func (m SvcDef) Name() string {
	return CollectionNmSvcDef
}

func (m SvcDef) PkKvPair() map[string]interface{} {
	return bson.M{"hash": m.Hash}
}
