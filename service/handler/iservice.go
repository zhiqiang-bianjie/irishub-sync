package handler

import (
	"fmt"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/store"
	"github.com/irisnet/irishub-sync/store/document"
	"github.com/irisnet/irishub-sync/types"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
)

func handleService(txs []document.CommonTx) {
	var batch []txn.Op
	for _, tx := range txs {
		if tx.Status != document.TxStatusSuccess {
			continue
		}
		switch tx.Type {
		case types.MsgSvcDef{}.Type():
			msgSvcDef := tx.MsgSvc.(types.MsgSvcDef)
			var svcDef = document.SvcDef{
				Hash:              tx.TxHash,
				Code:              msgSvcDef.Name,
				ChainId:           msgSvcDef.ChainId,
				Description:       msgSvcDef.Description,
				Author:            msgSvcDef.Author.String(),
				AuthorDescription: msgSvcDef.AuthorDescription,
				IDLContent:        msgSvcDef.IDLContent,
				Status:            "enable",
			}
			txOp := txn.Op{
				C:      document.CollectionNmSvcDef,
				Id:     bson.NewObjectId(),
				Insert: svcDef,
			}
			batch = append(batch, txOp)
		case types.MsgSvcBind{}.Type():
			msgSvcBind := tx.MsgSvc.(types.MsgSvcBind)
			var svcDef = document.SvcBind{
				ID:          GetSvcBindId(msgSvcBind.DefName, msgSvcBind.DefChainID),
				Hash:        tx.TxHash,
				DefName:     msgSvcBind.DefName,
				DefChainID:  msgSvcBind.DefChainID,
				BindChainID: msgSvcBind.BindChainID,
				Provider:    msgSvcBind.Provider.String(),
				BindingType: msgSvcBind.BindingType.String(),
				Deposit:     tx.Amount,
				Prices:      types.BuildCoins(msgSvcBind.Prices),
				Level: document.Level{
					AvgRspTime: msgSvcBind.Level.AvgRspTime,
					UsableTime: msgSvcBind.Level.UsableTime,
				},
				Available: true,
			}
			txOp := txn.Op{
				C:      document.CollectionNmSvcBind,
				Id:     bson.NewObjectId(),
				Insert: svcDef,
			}
			batch = append(batch, txOp)
		case types.MsgSvcBindingUpdate{}.Type():
			msgSvcBindUpdate := tx.MsgSvc.(types.MsgSvcBindingUpdate)
			var msgSvcBind = document.SvcBind{
				DefName:     msgSvcBindUpdate.DefName,
				DefChainID:  msgSvcBindUpdate.DefChainID,
				BindChainID: msgSvcBindUpdate.BindChainID,
				Provider:    msgSvcBindUpdate.Provider.String(),
				BindingType: msgSvcBindUpdate.BindingType.String(),
				Deposit:     tx.Amount,
				Level: document.Level{
					AvgRspTime: msgSvcBindUpdate.Level.AvgRspTime,
					UsableTime: msgSvcBindUpdate.Level.UsableTime,
				},
			}
			txOp := txn.Op{
				C:  document.CollectionNmSvcBind,
				Id: GetSvcBindId(msgSvcBindUpdate.DefName, msgSvcBindUpdate.DefChainID),
				Update: bson.M{
					"$set": msgSvcBind,
				},
			}
			batch = append(batch, txOp)
		case types.MsgSvcDisable{}.Type():
			msgSvcDisable := tx.MsgSvc.(types.MsgSvcDisable)
			var msgSvcBind = document.SvcBind{
				DefName:     msgSvcDisable.DefName,
				DefChainID:  msgSvcDisable.DefChainID,
				BindChainID: msgSvcDisable.BindChainID,
				Provider:    msgSvcDisable.Provider.String(),
				Available:   false,
			}
			txOp := txn.Op{
				C:  document.CollectionNmSvcBind,
				Id: GetSvcBindId(msgSvcDisable.DefName, msgSvcDisable.DefChainID),
				Update: bson.M{
					"$set": msgSvcBind,
				},
			}
			batch = append(batch, txOp)
		case types.MsgSvcEnable{}.Type():
			msgSvcEnable := tx.MsgSvc.(types.MsgSvcEnable)
			var msgSvcBind = document.SvcBind{
				DefName:     msgSvcEnable.DefName,
				DefChainID:  msgSvcEnable.DefChainID,
				BindChainID: msgSvcEnable.BindChainID,
				Provider:    msgSvcEnable.Provider.String(),
				Deposit:     tx.Amount,
				Available:   true,
			}
			txOp := txn.Op{
				C:  document.CollectionNmSvcBind,
				Id: GetSvcBindId(msgSvcEnable.DefName, msgSvcEnable.DefChainID),
				Update: bson.M{
					"$set": msgSvcBind,
				},
			}
			batch = append(batch, txOp)
		case types.MsgSvcRequest{}.Type():
			msgSvcRequest := tx.MsgSvc.(types.MsgSvcRequest)
			var msgInvocation = document.SvcInvocation{
				Hash:        tx.TxHash,
				ReqId:       tx.Tags["request-id"],
				TxType:      msgSvcRequest.Type(),
				DefChainID:  msgSvcRequest.DefChainID,
				DefName:     msgSvcRequest.DefName,
				BindChainID: msgSvcRequest.BindChainID,
				ReqChainID:  msgSvcRequest.ReqChainID,
				MethodID:    msgSvcRequest.MethodID,
				Consumer:    msgSvcRequest.Consumer.String(),
				Provider:    msgSvcRequest.Provider.String(),
				Height:      tx.Height,
				Data:        string(msgSvcRequest.Input),
				Time:        tx.Time,
			}
			txOp := txn.Op{
				C:      document.CollectionNmSvcInvocation,
				Id:     bson.NewObjectId(),
				Insert: msgInvocation,
			}
			batch = append(batch, txOp)
		case types.MsgSvcResponse{}.Type():
			msgSvcResponse := tx.MsgSvc.(types.MsgSvcResponse)
			var msgInvocation = document.SvcInvocation{
				Hash:     tx.TxHash,
				ReqId:    tx.Tags["request-id"],
				TxType:   msgSvcResponse.Type(),
				Provider: msgSvcResponse.Provider.String(),
				Consumer: tx.Tags["consumer"],
				Height:   tx.Height,
				Data:     string(msgSvcResponse.Output),
				Time:     tx.Time,
			}
			if reqs, err := msgInvocation.QueryByReqId(); err == nil && len(reqs) == 1 {
				msgInvocation.DefChainID = reqs[0].DefChainID
				msgInvocation.DefName = reqs[0].DefName
				msgInvocation.BindChainID = reqs[0].BindChainID
				msgInvocation.ReqChainID = reqs[0].ReqChainID
				msgInvocation.MethodID = reqs[0].MethodID
			}
			txOp := txn.Op{
				C:      document.CollectionNmSvcInvocation,
				Id:     bson.NewObjectId(),
				Insert: msgInvocation,
			}
			batch = append(batch, txOp)
		}
	}
	if len(batch) > 0 {
		if err := store.Txn(batch); err != nil {
			logger.Error("save service error", logger.String("err", err.Error()))
		}

	}
}
func GetSvcBindId(code, chainId string) string {
	return fmt.Sprintf("%s-%s", code, chainId)
}
