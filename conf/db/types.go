package db

import (
	"github.com/irisnet/irishub-sync/logger"
	"github.com/irisnet/irishub-sync/types"
	"os"
)

var (
	Addrs    = "127.0.0.1:27017"
	User     = "iris"
	Passwd   = "irispassword"
	Database = "sync-iris"
)

// get value of env var
func init() {
	addrs, found := os.LookupEnv(types.EnvNameDbAddr)
	if found {
		Addrs = addrs
	}
	logger.Info("Env Value", logger.String(types.EnvNameDbAddr, Addrs))

	user, found := os.LookupEnv(types.EnvNameDbUser)
	if found {
		User = user
	}
	logger.Info("Env Value", logger.String(types.EnvNameDbUser, User))

	passwd, found := os.LookupEnv(types.EnvNameDbPassWd)
	if found {
		Passwd = passwd
	}
	logger.Info("Env Value", logger.String(types.EnvNameDbPassWd, Passwd))

	database, found := os.LookupEnv(types.EnvNameDbDataBase)
	if found {
		Database = database
	}
	logger.Info("Env Value", logger.String(types.EnvNameDbDataBase, Database))
}
