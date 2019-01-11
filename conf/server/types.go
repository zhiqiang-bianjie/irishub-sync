package server

import (
	"github.com/irisnet/irishub-sync/types"
	"os"
	"strconv"
	"strings"

	"github.com/irisnet/irishub-sync/logger"
)

var (
	BlockChainMonitorUrl = []string{"tcp://127.0.0.1:26657"}
	ChainId              = "rainbow-dev"

	WorkerNumCreateTask  = 2
	WorkerNumExecuteTask = 60

	InitConnectionNum        = 50              // fast init num of tendermint client pool
	MaxConnectionNum         = 100             // max size of tendermint client pool
	CronWatchBlock           = "*/1 * * * * *" // every 1 seconds
	CronCalculateUpTime      = "0 */1 * * * *" // every minute
	CronCalculateTxGas       = "0 */5 * * * *" // every five minute
	SyncProposalStatus       = "0 */1 * * * *" // every minute
	CronSaveValidatorHistory = "@daily"        // every day
	CronUpdateDelegator      = "0/5 * * * * *" // every ten minute

	// deprecated
	SyncMaxGoroutine = 60 // max go routine in server
	// deprecated
	SyncBlockNumFastSync = 8000 // sync block num each goroutine

	ConsulAddr    = "192.168.150.7:8500"
	SyncWithDLock = false
	Bech32        = Bech32AddrPrefix{
		PrefixAccAddr:  "faa",
		PrefixAccPub:   "fap",
		PrefixValAddr:  "fva",
		PrefixValPub:   "fvp",
		PrefixConsAddr: "fca",
		PrefixConsPub:  "fcp",
	}
)

type Bech32AddrPrefix struct {
	PrefixAccAddr  string
	PrefixAccPub   string
	PrefixValAddr  string
	PrefixValPub   string
	PrefixConsAddr string
	PrefixConsPub  string
}

// get value of env var
func init() {
	nodeUrl, found := os.LookupEnv(types.EnvNameSerNetworkFullNode)
	if found {
		BlockChainMonitorUrl = strings.Split(nodeUrl, ",")
	}

	logger.Info("Env Value", logger.Any(types.EnvNameSerNetworkFullNode, BlockChainMonitorUrl))

	chainId, found := os.LookupEnv(types.EnvNameSerNetworkChainId)
	if found {
		ChainId = chainId
	}
	logger.Info("Env Value", logger.String(types.EnvNameSerNetworkChainId, ChainId))

	consulAddr, found := os.LookupEnv(types.EnvNameConsulAddr)
	if found {
		ConsulAddr = consulAddr
	}
	logger.Info("Env Value", logger.String(types.EnvNameConsulAddr, ConsulAddr))

	withDLock, found := os.LookupEnv(types.EnvNameSyncWithDLock)
	if found {
		flag, err := strconv.ParseBool(withDLock)
		if err != nil {
			logger.Fatal("Env Value", logger.String(types.EnvNameSyncWithDLock, withDLock))
		}
		SyncWithDLock = flag
	}
	logger.Info("Env Value", logger.Bool(types.EnvNameSyncWithDLock, SyncWithDLock))

	cronSaveValidatorHistory, found := os.LookupEnv(types.EnvNameCronSaveValidatorHistory)
	if found {
		CronSaveValidatorHistory = cronSaveValidatorHistory
	}
	logger.Info("Env Value", logger.String(types.EnvNameCronSaveValidatorHistory, cronSaveValidatorHistory))

	workerNumCreateTask, found := os.LookupEnv(types.EnvNameWorkerNumCreateTask)
	if found {
		var err error
		WorkerNumCreateTask, err = strconv.Atoi(workerNumCreateTask)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(types.EnvNameWorkerNumCreateTask, workerNumCreateTask))
		}
	}
	logger.Info("Env Value", logger.Int(types.EnvNameWorkerNumCreateTask, WorkerNumCreateTask))

	workerNumExecuteTask, found := os.LookupEnv(types.EnvNameWorkerNumExecuteTask)
	if found {
		var err error
		WorkerNumExecuteTask, err = strconv.Atoi(workerNumExecuteTask)
		if err != nil {
			logger.Fatal("Can't convert str to int", logger.String(types.EnvNameWorkerNumExecuteTask, workerNumCreateTask))
		}
	}
	logger.Info("Env Value", logger.Int(types.EnvNameWorkerNumExecuteTask, WorkerNumExecuteTask))

	loadBe32Prefix()
}

func loadBe32Prefix() {
	prefixAccAddr, found := os.LookupEnv(types.EnvNamePrefixAccAddr)
	if found {
		Bech32.PrefixAccAddr = prefixAccAddr
	}
	logger.Info("Env Value", logger.String(types.EnvNamePrefixAccAddr, Bech32.PrefixAccAddr))

	prefixAccPub, found := os.LookupEnv(types.EnvNamePrefixAccPub)
	if found {
		Bech32.PrefixAccPub = prefixAccPub
	}
	logger.Info("Env Value", logger.String(types.EnvNamePrefixAccPub, Bech32.PrefixAccPub))

	prefixValAddr, found := os.LookupEnv(types.EnvNamePrefixValAddr)
	if found {
		Bech32.PrefixValAddr = prefixValAddr
	}
	logger.Info("Env Value", logger.String(types.EnvNamePrefixValAddr, Bech32.PrefixValAddr))

	prefixValPub, found := os.LookupEnv(types.EnvNamePrefixValPub)
	if found {
		Bech32.PrefixValPub = prefixValPub
	}
	logger.Info("Env Value", logger.String(types.EnvNamePrefixValPub, Bech32.PrefixValPub))

	prefixConsAddr, found := os.LookupEnv(types.EnvNamePrefixConsAddr)
	if found {
		Bech32.PrefixConsAddr = prefixConsAddr
	}
	logger.Info("Env Value", logger.String(types.EnvNamePrefixConsAddr, Bech32.PrefixConsAddr))

	prefixConsPub, found := os.LookupEnv(types.EnvNamePrefixConsPub)
	if found {
		Bech32.PrefixConsPub = prefixConsPub
	}
	logger.Info("Env Value", logger.String(types.EnvNamePrefixConsPub, Bech32.PrefixConsPub))
}
