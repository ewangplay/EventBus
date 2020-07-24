package nsqd

import (
	"log"
	"math/rand"
	"time"

	"github.com/ewangplay/eventbus/config"
	"github.com/nsqio/nsq/nsqd"
)

type NsqdService struct {
	nsqd *nsqd.NSQD
}

func NewNsqdService() (*NsqdService, error) {
	return &NsqdService{}, nil
}

func (this *NsqdService) Start(opts *config.EB_Options) error {

	//set rand seed
	rand.Seed(time.Now().UTC().UnixNano())

	nsqOpts := this.convertToNSQOptions(opts)

	var err error
	nsqd, err := nsqd.New(nsqOpts)
	if err != nil {
		log.Fatalf("ERROR: %s", err.Error())
	}

	err = nsqd.LoadMetadata()
	if err != nil {
		log.Fatalf("ERROR: %s", err.Error())
	}
	err = nsqd.PersistMetadata()
	if err != nil {
		log.Fatalf("ERROR: failed to persist metadata - %s", err.Error())
	}
	nsqd.Main()

	this.nsqd = nsqd
	return nil
}

func (this *NsqdService) Stop() error {
	if this.nsqd != nil {
		this.nsqd.Exit()
	}
	return nil
}

func (this *NsqdService) convertToNSQOptions(opts *config.EB_Options) *nsqd.Options {
	nsqOpts := nsqd.NewOptions()

	nsqOpts.TCPAddress = opts.NSQTCPAddress
	nsqOpts.HTTPAddress = opts.NSQHTTPAddress
	nsqOpts.HTTPSAddress = opts.NSQHTTPSAddress
	nsqOpts.BroadcastAddress = opts.NSQBroadcastAddress

	// Convert EventBus LogLevel type to NSQ LogLevel type
	switch opts.LogLevel {
	case config.DEBUG:
		nsqOpts.LogLevel = 1
	case config.INFO:
		nsqOpts.LogLevel = 2
	case config.WARN:
		nsqOpts.LogLevel = 3
	case config.ERROR:
		nsqOpts.LogLevel = 4
	case config.FATAL:
		nsqOpts.LogLevel = 5
	default:
		nsqOpts.LogLevel = 2
	}
	nsqOpts.LogPrefix = opts.LogPrefix
	nsqOpts.Logger = opts.Logger

	nsqOpts.NSQLookupdTCPAddresses = opts.NSQLookupdTCPAddresses
	nsqOpts.AuthHTTPAddresses = opts.NSQAuthHTTPAddresses
	nsqOpts.HTTPClientConnectTimeout = opts.HTTPClientConnectTimeout
	nsqOpts.HTTPClientRequestTimeout = opts.HTTPClientRequestTimeout

	nsqOpts.DataPath = opts.DataPath
	nsqOpts.MemQueueSize = opts.MemQueueSize
	nsqOpts.MaxBytesPerFile = opts.MaxBytesPerFile
	nsqOpts.SyncEvery = opts.SyncEvery
	nsqOpts.SyncTimeout = opts.SyncTimeout

	nsqOpts.QueueScanInterval = opts.QueueScanInterval
	nsqOpts.QueueScanRefreshInterval = opts.QueueScanRefreshInterval
	nsqOpts.QueueScanSelectionCount = opts.QueueScanSelectionCount
	nsqOpts.QueueScanWorkerPoolMax = opts.QueueScanWorkerPoolMax
	nsqOpts.QueueScanDirtyPercent = opts.QueueScanDirtyPercent

	nsqOpts.MsgTimeout = opts.MsgTimeout
	nsqOpts.MaxMsgTimeout = opts.MaxMsgTimeout
	nsqOpts.MaxMsgSize = opts.MaxMsgSize
	nsqOpts.MaxBodySize = opts.MaxBodySize
	nsqOpts.MaxReqTimeout = opts.MaxReqTimeout
	nsqOpts.ClientTimeout = opts.ClientTimeout

	nsqOpts.MaxHeartbeatInterval = opts.MaxHeartbeatInterval
	nsqOpts.MaxRdyCount = opts.MaxRdyCount
	nsqOpts.MaxOutputBufferSize = opts.MaxOutputBufferSize
	nsqOpts.MaxOutputBufferTimeout = opts.MaxOutputBufferTimeout

	nsqOpts.StatsdAddress = opts.StatsdAddress
	nsqOpts.StatsdPrefix = opts.StatsdPrefix
	nsqOpts.StatsdInterval = opts.StatsdInterval
	nsqOpts.StatsdMemStats = opts.StatsdMemStats

	nsqOpts.E2EProcessingLatencyWindowTime = opts.E2EProcessingLatencyWindowTime
	nsqOpts.E2EProcessingLatencyPercentiles = opts.E2EProcessingLatencyPercentiles

	nsqOpts.TLSCert = opts.TLSCert
	nsqOpts.TLSKey = opts.TLSKey
	nsqOpts.TLSClientAuthPolicy = opts.TLSClientAuthPolicy
	nsqOpts.TLSRootCAFile = opts.TLSRootCAFile
	nsqOpts.TLSRequired = opts.TLSRequired
	nsqOpts.TLSMinVersion = opts.TLSMinVersion

	nsqOpts.DeflateEnabled = opts.DeflateEnabled
	nsqOpts.MaxDeflateLevel = opts.MaxDeflateLevel
	nsqOpts.SnappyEnabled = opts.SnappyEnabled

	return nsqOpts
}
