package metrics

import (
	"sync"
	"time"

	"github.com/alexcesaro/statsd"
	"github.com/pkg/errors"
	log "github.com/xlab/suplog"
)

var client Statter
var clientMux = new(sync.RWMutex)
var config *StatterConfig

type StatterConfig struct {
	EnvName              string
	HostName             string
	StuckFunctionTimeout time.Duration
	MockingEnabled       bool
}

func (m *StatterConfig) BaseTags() []string {
	var baseTags []string

	if len(config.EnvName) > 0 {
		baseTags = append(baseTags, "env", config.EnvName)
	}
	if len(config.HostName) > 0 {
		baseTags = append(baseTags, "machine", config.HostName)
	}

	return baseTags
}

type Statter interface {
	Count(bucket string, n interface{})
	Increment(bucket string)
	Gauge(bucket string, value interface{})
	Timing(bucket string, value interface{})
	Histogram(bucket string, value interface{})
	Unique(bucket string, value string)
	Close()
}

func Close() {
	clientMux.RLock()
	defer clientMux.RUnlock()
	if client == nil {
		return
	}
	client.Close()
}

func Disable() {
	config = checkConfig(nil)
	clientMux.Lock()
	client = newMockStatter(true)
	clientMux.Unlock()
}

func Init(addr string, prefix string, cfg *StatterConfig) error {
	config = checkConfig(cfg)
	if config.MockingEnabled {
		// init a mock statter instead of real statsd client
		clientMux.Lock()
		client = newMockStatter(false)
		clientMux.Unlock()
		return nil
	}
	statter, err := statsd.New(
		statsd.Address(addr),
		statsd.Prefix(prefix),
		statsd.ErrorHandler(errHandler),
		statsd.Tags(config.BaseTags()...),
	)
	if err != nil {
		err = errors.Wrap(err, "statsd init failed")
		return err
	}
	clientMux.Lock()
	client = statter
	clientMux.Unlock()
	return nil
}

func checkConfig(cfg *StatterConfig) *StatterConfig {
	if cfg == nil {
		cfg = &StatterConfig{}
	}
	if cfg.StuckFunctionTimeout < time.Second {
		cfg.StuckFunctionTimeout = 5 * time.Minute
	}
	if len(cfg.EnvName) == 0 {
		cfg.EnvName = "local"
	}
	return cfg
}

func errHandler(err error) {
	log.WithError(err).Errorln("statsd error")
}

func newMockStatter(noop bool) Statter {
	return &mockStatter{
		noop: noop,
		fields: log.Fields{
			"module": "mock_statter",
		},
	}
}

type mockStatter struct {
	fields log.Fields
	noop   bool
}

func (s *mockStatter) Count(bucket string, n interface{}) {
	if s.noop {
		return
	}
	log.WithFields(log.WithFn(s.fields)).Debugf("Bucket %s: %v", bucket, n)
}

func (s *mockStatter) Increment(bucket string) {
	if s.noop {
		return
	}
	log.WithFields(log.WithFn(s.fields)).Debugf("Bucket %s", bucket)
}

func (s *mockStatter) Gauge(bucket string, value interface{}) {
	if s.noop {
		return
	}
	log.WithFields(log.WithFn(s.fields)).Debugf("Bucket %s: %v", bucket, value)
}

func (s *mockStatter) Timing(bucket string, value interface{}) {
	if s.noop {
		return
	}
	log.WithFields(log.WithFn(s.fields)).Debugf("Bucket %s: %v", bucket, value)
}

func (s *mockStatter) Histogram(bucket string, value interface{}) {
	if s.noop {
		return
	}
	log.WithFields(log.WithFn(s.fields)).Debugf("Bucket %s: %v", bucket, value)
}

func (s *mockStatter) Unique(bucket string, value string) {
	if s.noop {
		return
	}
	log.WithFields(log.WithFn(s.fields)).Debugf("Bucket %s: %v", bucket, value)
}

func (s *mockStatter) Close() {
	if s.noop {
		return
	}
	log.WithFields(log.WithFn(s.fields)).Debugf("closed at %s", time.Now())
}
