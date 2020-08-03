package adapter

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ewangplay/eventbus/adapter/redis"
	c "github.com/ewangplay/eventbus/common"
	"github.com/ewangplay/eventbus/config"
	"github.com/ewangplay/eventbus/i"
)

// JobManager struct define, implement JobManager interface
type JobManager struct {
	i.Logger
	i.Producer
	mutexEvent     sync.RWMutex
	mutexFailEvent sync.RWMutex
	opts           *config.EBOptions
	redisCtx       *redis.Context
	redisClient    *redis.Client
	quit           chan int
}

// NewJobManager ...
func NewJobManager(opts *config.EBOptions, logger i.Logger, producer i.Producer) (*JobManager, error) {
	jm := &JobManager{}
	jm.opts = opts
	jm.Logger = logger
	jm.Producer = producer

	redisCtx, err := redis.GetContext(opts, logger)
	if err != nil {
		return nil, err
	}

	redisClient, err := redisCtx.CreateRedisClient()
	if err != nil {
		return nil, err
	}

	jm.redisCtx = redisCtx
	jm.redisClient = redisClient

	jm.quit = make(chan int)
	go jm.work(jm.quit)

	return jm, nil
}

// Close ...
func (jm *JobManager) Close() error {
	jm.Info("Job manager will close")
	close(jm.quit)
	return jm.redisClient.Close()
}

// Define redis key name
const (
	EventKey     = "EVENTBUS:EVENT"
	EventFailKey = "EVENTBUS:EVENT:FAIL"
	MutexKey     = "EVENTBUS:MUTEX"
)

// Set event to job manager
func (jm *JobManager) Set(job i.Event) error {
	jm.mutexEvent.Lock()
	defer jm.mutexEvent.Unlock()

	err := jm.redisClient.HSet(EventKey, job.GetID(), job.GetData())
	if err != nil {
		jm.Error("Set event[%s] into status table error: %v", job.GetData(), err)
		return err
	}
	return nil
}

// Fail add jm job into FAIL table
func (jm *JobManager) Fail(job i.Event) error {
	var field string

	eventID := job.GetID()
	retryPolicy := job.GetRetryPolicy()
	retryCount := job.GetRetryCount()
	retryInterval := job.GetRetryInterval()
	retryTimeout := job.GetRetryTimeout()
	createTime := job.GetCreateTime()
	updateTime := job.GetUpdateTime()
	deadline := createTime + retryTimeout

	if retryPolicy == c.CountRetryPolicy {
		field = fmt.Sprintf("%s:%d:%d:%d:%d",
			eventID,
			retryPolicy,
			retryCount,
			updateTime,
			retryInterval)
	} else if retryPolicy == c.ExpiredRetryPolicy {
		field = fmt.Sprintf("%s:%d:%d:%d:%d",
			eventID,
			retryPolicy,
			deadline,
			updateTime,
			retryInterval)
	} else {
		return fmt.Errorf("invalid retry policy: %d", retryPolicy)
	}

	jm.Debug("Add event[%s] into FAIL table...", job.GetID())

	jm.mutexFailEvent.Lock()
	err := jm.redisClient.HSet(EventFailKey, field, job.GetData())
	if err != nil {
		jm.Error("Add event[%s] into FAIL table error: %v", job.GetID(), err)
		jm.mutexFailEvent.Unlock()
		return err
	}
	jm.mutexFailEvent.Unlock()

	return nil
}

// Get event info from job manager
func (jm *JobManager) Get(eventID string) ([]byte, error) {
	jm.mutexEvent.RLock()
	defer jm.mutexEvent.RUnlock()
	return jm.redisClient.HGet(EventKey, eventID)
}

func (jm *JobManager) work(quit chan int) {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-quit:
			jm.Debug("quit job manager bg worker")
			return
		case <-ticker.C:
			jm.Debug("Start to clear up job...")
			ok, err := jm.getMutexFlag()
			if err == nil && ok {
				jm.ProcessFailedJobs()
				jm.resetMutexFlag()
			}
		}
	}
}

// ProcessFailedJobs ...
func (jm *JobManager) ProcessFailedJobs() error {
	jm.mutexFailEvent.RLock()
	failedJobs, err := jm.redisClient.HGetAll(EventFailKey)
	if err != nil {
		jm.Error("Get failed events error: %v", err)
		jm.mutexFailEvent.RUnlock()
		return err
	}
	jm.mutexFailEvent.RUnlock()

	var event c.EBEvent
	var valid, trigger bool

	for field, data := range failedJobs {

		//Parse event from data
		err = json.Unmarshal([]byte(data), &event)
		if err != nil {
			jm.Error("Unmarshal data[%s] error：%v", data, err)
			return err
		}

		//Check the validity of failed event
		valid, trigger = jm.checkValidity(field)
		if valid {
			if !trigger {
				continue
			}

			jm.Debug("Reschedule the failed event[%s]...", field)

			//Re-register jm job
			event.Status = c.EventStatusInit
			event.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
			err = jm.Set(&event)
			if err != nil {
				jm.Error("Re-init event[%s] error：%v", field, err)
				return err
			}

			//Publish jm event to MQ
			//topic := fmt.Sprintf("%s.%s", event.GetType(), event.GetSubject())
			err = jm.Publish(&c.EBMessage{Subject: event.GetType(), Data: event.GetData()})
			if err != nil {
				jm.Error("Publish event[%s] error：%v", field, err)
				return err
			}

		} else {
			jm.Debug("Drop the failed event[%s]...", field)

			//drop jm job
			event.Status = c.EventStatusDrop
			event.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
			err = jm.Set(&event)
			if err != nil {
				jm.Error("Drop event[%s] error：%v", field, err)
				return err
			}
		}

		//Remove from FAIL table
		jm.mutexFailEvent.Lock()
		err = jm.redisClient.HDel(EventFailKey, field)
		if err != nil {
			jm.Error("Delete event[%s] from FAIL table error: %v", field, err)
			jm.mutexFailEvent.Unlock()
			return err
		}
		jm.mutexFailEvent.Unlock()
	}
	return nil
}

func (jm *JobManager) getMutexFlag() (bool, error) {
	exist, err := jm.redisClient.Exists(MutexKey)
	if err != nil {
		jm.Error("Check clearup flag[%s] error: %v", MutexKey, err)
		return false, err
	}

	if !exist {
		err = jm.redisClient.Set(MutexKey, "1")
		if err != nil {
			jm.Error("Set clearup flag[%s] error: %v", MutexKey, err)
			return false, err
		}
		return true, nil
	}

	flag, err := jm.redisClient.Get(MutexKey)
	if err != nil {
		jm.Error("Get clearup flag[%s] error: %v", MutexKey, err)
		return false, err
	}

	if flag == "0" {
		err = jm.redisClient.Set(MutexKey, "1")
		if err != nil {
			jm.Error("Set clearup flag[%s] error: %v", MutexKey, err)
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func (jm *JobManager) resetMutexFlag() error {
	err := jm.redisClient.Set(MutexKey, "0")
	if err != nil {
		jm.Error("Set clearup flag[%s] error: %v", MutexKey, err)
		return err
	}
	return nil

}

func (jm *JobManager) checkValidity(field string) (valid bool, trigger bool) {
	var err error

	elems := strings.Split(field, ":")
	if elems == nil {
		return false, false
	}
	if len(elems) != 5 {
		return false, false
	}

	//eventID := elems[0]

	var retryPolicy, retryCount, retryInterval, updateTime, deadline uint64

	retryPolicy, err = strconv.ParseUint(elems[1], 10, 64)
	if err != nil {
		return false, false
	}

	if int(retryPolicy) == c.CountRetryPolicy {
		retryCount, err = strconv.ParseUint(elems[2], 10, 64)
		if err != nil {
			return false, false
		}
		if retryCount == 0 {
			return false, false
		}

		updateTime, err = strconv.ParseUint(elems[3], 10, 64)
		if err != nil {
			return false, false
		}
		retryInterval, err = strconv.ParseUint(elems[4], 10, 64)
		if err != nil {
			return false, false
		}

		now := uint64(time.Now().Unix())

		if now-updateTime >= retryInterval {
			return true, true
		}
		return true, false

	} else if int(retryPolicy) == c.ExpiredRetryPolicy {
		deadline, err = strconv.ParseUint(elems[2], 10, 64)
		if err != nil {
			return false, false
		}
		updateTime, err = strconv.ParseUint(elems[3], 10, 64)
		if err != nil {
			return false, false
		}
		retryInterval, err = strconv.ParseUint(elems[4], 10, 64)
		if err != nil {
			return false, false
		}

		if updateTime >= deadline {
			return false, false
		}

		now := uint64(time.Now().Unix())
		if now >= deadline {
			return false, false
		}
		if now-updateTime >= retryInterval {
			return true, true
		}
		return true, false

	}

	return false, false
}
