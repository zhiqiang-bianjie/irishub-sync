package service

import (
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/rpc"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	"sync"
	"time"
)

// save account
func SaveAccount(docTx document.CommonTx) {
	var (
		address    string
		updateTime time.Time
		height     int64
		methodName = "SaveAccount: "
	)
	logger.Debug("Start", logger.String("method", methodName))

	// save account
	fun := func(address string, updateTime time.Time, height int64) {
		account := document.Account{
			Address: address,
			Time:    updateTime,
			Height:  height,
		}

		err := store.Save(account)

		if err != nil && err.Error() != "Record exists" {
			logger.Error("account Record exists", logger.String("address", account.Address))
		}
	}

	txType := GetTxType(docTx)
	if len(txType) == 0 {
		logger.Error("Tx is valid", logger.Any("Tx", docTx))
		return
	}

	switch txType {
	case types.TxTypeTransfer, types.TxTypeStakeDelegate,
		types.TxTypeStakeBeginUnbonding, types.TxTypeStakeCompleteUnbonding:
		updateTime = docTx.Time
		height = docTx.Height

		fun(docTx.From, updateTime, height)
		fun(docTx.To, updateTime, height)
		break
	case types.TxTypeStakeCreateValidator, types.TxTypeStakeEditValidator:
		address = docTx.From
		updateTime = docTx.Time
		height = docTx.Height

		fun(address, updateTime, height)
		break
	}

	logger.Debug("End", logger.String("method", methodName))
}

// update account balance
func UpdateBalance(docTx document.CommonTx, mutex sync.Mutex) {
	var (
		methodName = "UpdateBalance: "
	)
	logger.Debug("Start", logger.String("method", methodName))

	fun := func(address string) {
		account, err := document.QueryAccount(address)
		if err != nil {
			logger.Error("QueryAccount failed", logger.String("address", address), logger.String("err", err.Error()))
			return
		}

		// query balance of account
		account.Amount = store.ParseCoins(rpc.GetBalance(address))
		if err := store.Update(account); err != nil {
			logger.Error("updateAccountBalance failed", logger.String("address", account.Address), logger.String("err", err.Error()))
		}
	}

	txType := GetTxType(docTx)
	if txType == "" {
		logger.Error("Tx is valid", logger.Any("Tx", docTx))
		return
	}

	switch txType {
	case types.TxTypeTransfer, types.TxTypeStakeDelegate,
		types.TxTypeStakeBeginUnbonding, types.TxTypeStakeCompleteUnbonding:
		fun(docTx.From)
		fun(docTx.To)
		break
	case types.TxTypeStakeCreateValidator, types.TxTypeStakeEditValidator:
		fun(docTx.From)
		break
	}

	logger.Debug("End", logger.String("method", methodName))
}
