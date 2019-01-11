// package for parse tx struct from binary data

package rpc

import (
	itypes "github.com/irisnet/irishub-sync/types"
)

// get tx status and log by query txHash
func GetTxResult(txHash []byte) (itypes.ResponseDeliverTx, error) {
	var resDeliverTx itypes.ResponseDeliverTx

	client := GetClient()
	defer client.Release()

	res, err := client.Tx(txHash, false)
	if err != nil {
		return resDeliverTx, err
	}
	result := res.TxResult

	return result, nil
}
