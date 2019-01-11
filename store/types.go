// interface for a document

package store

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	CollectionNameTxn = "mgo_txn"
)

type Docs interface {
	// collection name
	Name() string
	// primary key pair(used to find a unique record)
	PkKvPair() map[string]interface{}
}

type Coin struct {
	Denom  string  `json:"denom"`
	Amount float64 `json:"amount"`
}

type Coins []Coin

type Fee struct {
	Amount Coins
	Gas    int64
}

type ActualFee struct {
	Denom  string  `json:"denom"`
	Amount float64 `json:"amount"`
}

type Msg interface {
	Type() string
	String() string
}

func ParseCoins(coinsStr string) (coins Coins) {
	coinsStr = strings.TrimSpace(coinsStr)
	if len(coinsStr) == 0 {
		return
	}

	coinStrs := strings.Split(coinsStr, ",")
	for _, coinStr := range coinStrs {
		coin := ParseCoin(coinStr)
		coins = append(coins, coin)
	}
	return coins
}

func ParseCoin(coinStr string) (coin Coin) {
	var (
		reDnm  = `[A-Za-z\-]{2,15}`
		reAmt  = `[0-9]+[.]?[0-9]*`
		reSpc  = `[[:space:]]*`
		reCoin = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reAmt, reSpc, reDnm))
	)

	coinStr = strings.TrimSpace(coinStr)

	matches := reCoin.FindStringSubmatch(coinStr)
	if matches == nil {
		return coin
	}
	denom, amount := matches[2], matches[1]

	amt, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return coin
	}

	return Coin{
		Denom:  denom,
		Amount: amt,
	}
}

//func BuildFee(fee auth.StdFee) store.Fee {
//	return store.Fee{
//		Amount: ParseCoins(fee.Amount.String()),
//		Gas:    int64(fee.Gas),
//	}
//}
