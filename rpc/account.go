// This package is used for Query balance of account

package rpc

import (
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/types"
)

// query account balance from sdk store
func GetBalance(address string) store.Coins {
	cdc := types.GetCodec()

	addr, err := types.AccAddressFromBech32(address)
	if err != nil {
		logger.Error("get addr from hex failed", logger.Any("err", err))
		return nil
	}

	res, err := Query(types.AddressStoreKey(addr), "acc",
		types.StoreDefaultEndPath)

	if err != nil {
		logger.Error("Query balance from tendermint failed", logger.Any("err", err))
		return nil
	}

	// balance is empty
	if len(res) <= 0 {
		return nil
	}

	decoder := types.GetAccountDecoder(cdc)
	account, err := decoder(res)
	if err != nil {
		logger.Error("decode account failed", logger.Any("err", err))
		return nil
	}

	return store.ParseCoins(account.GetCoins().String())
}

func ValAddrToAccAddr(address string) (accAddr string) {
	valAddr, err := types.ValAddressFromBech32(address)
	if err != nil {
		logger.Error("ValAddressFromBech32 decode account failed", logger.String("address", address))
		return
	}

	return types.AccAddress(valAddr.Bytes()).String()
}
