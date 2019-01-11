package service

import (
	"github.com/irisnet/irishub-sync/store/document"
)

// get tx type
func GetTxType(docTx document.CommonTx) string {
	if docTx.TxHash == "" {
		return ""
	}
	return docTx.Type
}

type Service = func(tx document.CommonTx)

func Execute(docTx document.CommonTx, actions []Service) {
	for _, action := range actions {
		if docTx.TxHash != "" {
			action(docTx)
		}
	}

}
