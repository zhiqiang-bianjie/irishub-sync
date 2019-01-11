package core

import (
	"github.com/irisnet/irishub-sync/core/crons"
	"github.com/irisnet/irishub-sync/core/service"
	"github.com/irisnet/irishub-sync/logger"
	"github.com/robfig/cron"
	"time"
)

var (
	engine *SyncEngine
)

func init() {
	engine = &SyncEngine{
		cron:      cron.New(),
		tasks:     []crons.Task{},
		initFuncs: []func(){},
	}

	engine.AddTask(crons.MakeCalculateAndSaveValidatorUpTimeTask())
	engine.AddTask(crons.MakeCalculateTxGasAndGasPriceTask())
	engine.AddTask(crons.MakeSyncProposalStatusTask())
	engine.AddTask(crons.MakeValidatorHistoryTask())
	engine.AddTask(crons.MakeUpdateDelegatorTask())

	// init delegator for genesis validator
	engine.initFuncs = append(engine.initFuncs, service.InitDelegator)
}

type SyncEngine struct {
	cron      *cron.Cron   //cron
	tasks     []crons.Task // my timer task
	initFuncs []func()     // module init fun
}

func (engine *SyncEngine) AddTask(task crons.Task) {
	engine.tasks = append(engine.tasks, task)
	engine.cron.AddFunc(task.GetCron(), task.GetCommand())
}

func (engine *SyncEngine) init() {
	// init module info
	for _, init := range engine.initFuncs {
		init()
	}
}

func (engine *SyncEngine) Start() {
	engine.init()

	go startCreateTask()
	go startExecuteTask()

	// cron task should start after fast sync finished
	fastSyncChan := make(chan bool, 1)
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			<-ticker.C
			flag, err := assertFastSyncFinished()
			if err != nil {
				logger.Error("assert fast sync finished failed", logger.String("err", err.Error()))
			}
			if flag {
				close(fastSyncChan)
				return
			}
		}
	}()
	<-fastSyncChan
	logger.Info("fast sync finished, now cron task can start")

	engine.cron.Start()
}

func (engine *SyncEngine) Stop() {
	logger.Info("release resource :SyncEngine")
	engine.cron.Stop()
	for _, t := range engine.tasks {
		t.Release()
	}
}

func New() *SyncEngine {
	return engine
}
