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

type JobManager struct {
	i.ILogger
	i.IProducer
	mutexEvent     sync.RWMutex
	mutexFailEvent sync.RWMutex
	opts           *config.EB_Options
	redisCtx       *redis.RedisContext
	redisClient    *redis.RedisClient
	quit           chan int
}

func NewJobManager(opts *config.EB_Options, logger i.ILogger, producer i.IProducer) (*JobManager, error) {
	this := &JobManager{}
	this.opts = opts
	this.ILogger = logger
	this.IProducer = producer

	redisCtx, err := redis.GetRedisContext(opts, logger)
	if err != nil {
		return nil, err
	}

	redisClient, err := redisCtx.CreateRedisClient()
	if err != nil {
		return nil, err
	}

	this.redisCtx = redisCtx
	this.redisClient = redisClient

	this.quit = make(chan int)
	go this.work(this.quit)

	return this, nil
}

func (this *JobManager) Close() error {
	this.Info("Job manager will close")
	close(this.quit)
	return this.redisClient.Close()
}

const (
	EVENT_KEY      = "EVENTBUS:EVENT"
	EVENT_FAIL_KEY = "EVENTBUS:EVENT:FAIL"
	MUTEX_KEY      = "EVENTBUS:MUTEX"
)

func (this *JobManager) Set(job i.IEvent) error {
	this.mutexEvent.Lock()
	defer this.mutexEvent.Unlock()

	err := this.redisClient.HSet(EVENT_KEY, job.GetId(), job.GetData())
	if err != nil {
		this.Error("Set event[%s] into status table error: %v", job.GetData(), err)
		return err
	}
	return nil
}

//Add this job into FAIL table
func (this *JobManager) Fail(job i.IEvent) error {
	var field string

	event_id := job.GetId()
	retry_policy := job.GetRetryPolicy()
	retry_count := job.GetRetryCount()
	retry_interval := job.GetRetryInterval()
	retry_timeout := job.GetRetryTimeout()
	create_time := job.GetCreateTime()
	update_time := job.GetUpdateTime()
	deadline := create_time + retry_timeout

	if retry_policy == c.ERP_COUNT {
		field = fmt.Sprintf("%s:%d:%d:%d:%d",
			event_id,
			retry_policy,
			retry_count,
			update_time,
			retry_interval)
	} else if retry_policy == c.ERP_TIMEOUT {
		field = fmt.Sprintf("%s:%d:%d:%d:%d",
			event_id,
			retry_policy,
			deadline,
			update_time,
			retry_interval)
	} else {
		return fmt.Errorf("invalid retry policy: %d", retry_policy)
	}

	this.Debug("Add event[%s] into FAIL table...", job.GetId())

	this.mutexFailEvent.Lock()
	err := this.redisClient.HSet(EVENT_FAIL_KEY, field, job.GetData())
	if err != nil {
		this.Error("Add event[%s] into FAIL table error: %v", job.GetId(), err)
		this.mutexFailEvent.Unlock()
		return err
	}
	this.mutexFailEvent.Unlock()

	return nil
}

func (this *JobManager) Get(event_id string) ([]byte, error) {
	this.mutexEvent.RLock()
	defer this.mutexEvent.RUnlock()
	return this.redisClient.HGet(EVENT_KEY, event_id)
}

func (this *JobManager) work(quit chan int) {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-quit:
			this.Debug("quit job manager bg worker")
			return
		case <-ticker.C:
			this.Debug("Start to clear up job...")
			ok, err := this.getMutexFlag()
			if err == nil && ok {
				this.ProcessFailedJobs()
				this.resetMutexFlag()
			}
		}
	}
}

func (this *JobManager) ProcessFailedJobs() error {
	this.mutexFailEvent.RLock()
	failed_jobs, err := this.redisClient.HGetAll(EVENT_FAIL_KEY)
	if err != nil {
		this.Error("Get failed events error: %v", err)
		this.mutexFailEvent.RUnlock()
		return err
	}
	this.mutexFailEvent.RUnlock()

	var event c.EB_Event
	var valid, trigger bool

	for field, data := range failed_jobs {

		//Parse event from data
		err = json.Unmarshal([]byte(data), &event)
		if err != nil {
			this.Error("Unmarshal data[%s] error：%v", data, err)
			return err
		}

		//Check the validity of failed event
		valid, trigger = this.checkValidity(field)
		if valid {
			if !trigger {
				continue
			}

			this.Debug("Reschedule the failed event[%s]...", field)

			//Re-register this job
			event.Status = c.ES_INIT
			event.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
			err = this.Set(&event)
			if err != nil {
				this.Error("Re-init event[%s] error：%v", field, err)
				return err
			}

			//Publish this event to MQ
			//topic := fmt.Sprintf("%s.%s", event.GetType(), event.GetSubject())
			err = this.Publish(&c.EB_Message{Subject: event.GetType(), Data: event.GetData()})
			if err != nil {
				this.Error("Publish event[%s] error：%v", field, err)
				return err
			}

		} else {
			this.Debug("Drop the failed event[%s]...", field)

			//drop this job
			event.Status = c.ES_DROP
			event.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
			err = this.Set(&event)
			if err != nil {
				this.Error("Drop event[%s] error：%v", field, err)
				return err
			}
		}

		//Remove from FAIL table
		this.mutexFailEvent.Lock()
		err = this.redisClient.HDel(EVENT_FAIL_KEY, field)
		if err != nil {
			this.Error("Delete event[%s] from FAIL table error: %v", field, err)
			this.mutexFailEvent.Unlock()
			return err
		}
		this.mutexFailEvent.Unlock()
	}
	return nil
}

func (this *JobManager) getMutexFlag() (bool, error) {
	exist, err := this.redisClient.Exists(MUTEX_KEY)
	if err != nil {
		this.Error("Check clearup flag[%s] error: %v", MUTEX_KEY, err)
		return false, err
	}

	if !exist {
		err = this.redisClient.Set(MUTEX_KEY, "1")
		if err != nil {
			this.Error("Set clearup flag[%s] error: %v", MUTEX_KEY, err)
			return false, err
		}
		return true, nil
	}

	flag, err := this.redisClient.Get(MUTEX_KEY)
	if err != nil {
		this.Error("Get clearup flag[%s] error: %v", MUTEX_KEY, err)
		return false, err
	}

	if flag == "0" {
		err = this.redisClient.Set(MUTEX_KEY, "1")
		if err != nil {
			this.Error("Set clearup flag[%s] error: %v", MUTEX_KEY, err)
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func (this *JobManager) resetMutexFlag() error {
	err := this.redisClient.Set(MUTEX_KEY, "0")
	if err != nil {
		this.Error("Set clearup flag[%s] error: %v", MUTEX_KEY, err)
		return err
	}
	return nil

}

func (this *JobManager) checkValidity(field string) (valid bool, trigger bool) {
	var err error

	elems := strings.Split(field, ":")
	if elems == nil {
		return false, false
	}
	if len(elems) != 5 {
		return false, false
	}

	//event_id := elems[0]

	var retry_policy, retry_count, retry_interval, update_time, deadline uint64

	retry_policy, err = strconv.ParseUint(elems[1], 10, 64)
	if err != nil {
		return false, false
	}

	if int(retry_policy) == c.ERP_COUNT {
		retry_count, err = strconv.ParseUint(elems[2], 10, 64)
		if err != nil {
			return false, false
		}
		if retry_count == 0 {
			return false, false
		}

		update_time, err = strconv.ParseUint(elems[3], 10, 64)
		if err != nil {
			return false, false
		}
		retry_interval, err = strconv.ParseUint(elems[4], 10, 64)
		if err != nil {
			return false, false
		}

		now := uint64(time.Now().Unix())

		if now-update_time >= retry_interval {
			return true, true
		} else {
			return true, false
		}

	} else if int(retry_policy) == c.ERP_TIMEOUT {
		deadline, err = strconv.ParseUint(elems[2], 10, 64)
		if err != nil {
			return false, false
		}
		update_time, err = strconv.ParseUint(elems[3], 10, 64)
		if err != nil {
			return false, false
		}
		retry_interval, err = strconv.ParseUint(elems[4], 10, 64)
		if err != nil {
			return false, false
		}

		if update_time >= deadline {
			return false, false
		}

		now := uint64(time.Now().Unix())
		if now >= deadline {
			return false, false
		}
		if now-update_time >= retry_interval {
			return true, true
		} else {
			return true, false
		}

	}

	return false, false
}
