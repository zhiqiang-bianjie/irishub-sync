package rpc

import (
	"context"
	conf "github.com/irisnet/irishub-sync/conf/server"
	"github.com/irisnet/irishub-sync/logger"
	gcp "github.com/jolestar/go-commons-pool"
	"github.com/robfig/cron"
	"math/rand"
	"sync"
)

var (
	factory PoolFactory
	pool    *NodePool
	ctx     = context.Background()
)

func init() {
	var syncMap sync.Map
	for _, url := range conf.BlockChainMonitorUrl {
		key := generateId(url)
		endPoint := EndPoint{
			Address:   url,
			Available: true,
		}

		syncMap.Store(key, endPoint)
	}
	factory = PoolFactory{
		cron:     cron.New(),
		peersMap: syncMap,
	}
	config := gcp.NewDefaultPoolConfig()

	config.MaxTotal = conf.MaxConnectionNum
	config.MaxIdle = conf.InitConnectionNum
	config.MinIdle = conf.InitConnectionNum
	config.TestOnBorrow = true
	config.TestOnCreate = true
	config.TestWhileIdle = true

	logger.Info("PoolConfig", logger.Int("config.MaxTotal", config.MaxTotal), logger.Int("config.MaxIdle", config.MaxIdle))
	pool = &NodePool{
		gcp.NewObjectPool(ctx, &factory, config),
	}
}

type EndPoint struct {
	Address   string
	Available bool
}

type NodePool struct {
	*gcp.ObjectPool
}

type PoolFactory struct {
	peersMap sync.Map
	cron     *cron.Cron
}

func Close() {
	logger.Info("release resource nodePool")
	pool.Close(ctx)
	factory.cron.Stop()
}

func (f *PoolFactory) MakeObject(ctx context.Context) (*gcp.PooledObject, error) {
	endpoint := f.GetEndPoint()
	logger.Debug("PoolFactory MakeObject peer", logger.Any("endpoint", endpoint))
	return gcp.NewPooledObject(newClient(endpoint.Address)), nil
}

func (f *PoolFactory) DestroyObject(ctx context.Context, object *gcp.PooledObject) error {
	logger.Debug("PoolFactory DestroyObject peer", logger.Any("peer", object.Object))
	c := object.Object.(*Client)
	if c.IsRunning() {
		c.Stop()
	}
	return nil
}

func (f *PoolFactory) ValidateObject(ctx context.Context, object *gcp.PooledObject) bool {
	// do validate
	logger.Debug("PoolFactory ValidateObject peer", logger.Any("peer", object.Object))
	c := object.Object.(*Client)
	if c.HeartBeat() != nil {
		value, ok := f.peersMap.Load(c.Id)
		if ok {
			endPoint := value.(EndPoint)
			endPoint.Available = true
			f.peersMap.Store(c.Id, endPoint)
		}
		return false
	}
	return true
}

func (f *PoolFactory) ActivateObject(ctx context.Context, object *gcp.PooledObject) error {
	logger.Debug("PoolFactory ActivateObject peer", logger.Any("peer", object.Object))
	return nil
}

func (f *PoolFactory) PassivateObject(ctx context.Context, object *gcp.PooledObject) error {
	logger.Debug("PoolFactory PassivateObject peer", logger.Any("peer", object.Object))
	return nil
}

func (f *PoolFactory) GetEndPoint() EndPoint {
	var (
		keys        []string
		selectedKey string
	)

	f.peersMap.Range(func(k, value interface{}) bool {
		key := k.(string)
		endPoint := value.(EndPoint)
		if endPoint.Available {
			keys = append(keys, key)
		}
		selectedKey = key

		return true
	})

	if len(keys) > 0 {
		index := rand.Intn(len(keys))
		selectedKey = keys[index]
	}
	value, ok := f.peersMap.Load(selectedKey)
	if ok {
		return value.(EndPoint)
	} else {
		logger.Error("Can't get selected end point", logger.String("selectedKey", selectedKey))
	}
	return EndPoint{}
}
