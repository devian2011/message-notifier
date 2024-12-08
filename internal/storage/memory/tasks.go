package memory

import (
	"context"
	"github.com/sirupsen/logrus"
	"notifier/pkg/swap"
	"sync"
	"time"

	"github.com/google/uuid"

	"notifier/internal/entity"
	cfgParser "notifier/pkg/config"
)

const (
	defaultSwapDelay = 10 * time.Second
	defaultSwapFile  = "./store.bin"

	defaultHoldDoneDelay = 1 * time.Hour
	defaultHoldDoneTasks = 24 * time.Hour
)

type config struct {
	swapFile      string
	swapDelay     time.Duration
	holdDoneDelay time.Duration
	holdDoneTasks time.Duration
}

type DataStore struct {
	Process map[uuid.UUID]entity.MessageTask
	Closed  map[uuid.UUID]entity.MessageTask
}

type Storage struct {
	ctx  context.Context
	cfg  *config
	data *DataStore

	processRwMtx *sync.RWMutex
	closedRwMtx  *sync.RWMutex
}

func NewStorage(ctx context.Context, params map[string]string) (*Storage, error) {
	cfg, cfgErr := buildCfg(params)
	if cfgErr != nil {
		return nil, cfgErr
	}

	storage := &Storage{
		ctx:          ctx,
		cfg:          cfg,
		processRwMtx: &sync.RWMutex{},
		closedRwMtx:  &sync.RWMutex{},
		data: &DataStore{
			Process: make(map[uuid.UUID]entity.MessageTask),
			Closed:  make(map[uuid.UUID]entity.MessageTask),
		},
	}

	loadErr := storage.Load()
	logrus.WithFields(logrus.Fields{
		"error":    loadErr,
		"context":  "storage",
		"provider": "memory",
	}).Errorln("error on load data from disk")

	return storage, nil
}

func (s *Storage) GetData() *DataStore {
	return s.data
}

func (s *Storage) SetData(data *DataStore) {
	s.data = data
}

func (s *Storage) Save(msg entity.MessageTask) error {
	if msg.Status == entity.TaskStatusToProcess || msg.Status == entity.TaskStatusInProgress {
		s.processRwMtx.Lock()
		defer s.processRwMtx.Unlock()
		s.data.Process[msg.Id] = msg
	}
	if msg.Status == entity.TaskStatusDone || msg.Status == entity.TaskStatusError {
		s.closedRwMtx.Lock()
		s.processRwMtx.Lock()
		defer s.closedRwMtx.Unlock()
		defer s.processRwMtx.Unlock()

		s.data.Closed[msg.Id] = msg

		if _, exists := s.data.Process[msg.Id]; exists {
			delete(s.data.Process, msg.Id)
		}
	}

	return nil
}

func (s *Storage) Run() {
	go func() {
	loop:
		for {
			select {
			case <-s.ctx.Done():
				break loop
			case <-time.After(s.cfg.swapDelay):
				swapErr := s.Swap()
				logrus.WithFields(logrus.Fields{
					"error":    swapErr,
					"context":  "storage",
					"provider": "memory",
				}).Errorln("error on swap data to disk")
			case <-time.After(s.cfg.holdDoneDelay):
				s.ClearExpired()
			}
		}
	}()
}

func (s *Storage) Swap() error {
	s.processRwMtx.Lock()
	s.closedRwMtx.Lock()
	defer s.processRwMtx.Unlock()
	defer s.closedRwMtx.Unlock()

	return swap.Swap(s.cfg.swapFile, s.data)
}

func (s *Storage) Load() error {
	s.processRwMtx.Lock()
	s.closedRwMtx.Lock()
	defer s.processRwMtx.Unlock()
	defer s.closedRwMtx.Unlock()

	return swap.Load(s.cfg.swapFile, s.data)
}

func (s *Storage) ClearExpired() {
	s.closedRwMtx.Lock()
	defer s.closedRwMtx.Unlock()
	now := time.Now()
	for i := range s.data.Closed {
		t := s.data.Closed[i].LastExecutionTime.Add(s.cfg.holdDoneTasks)
		if now.After(t) {
			delete(s.data.Closed, i)
		}
	}
}

func buildCfg(params map[string]string) (*config, error) {
	cfg := &config{}

	if swapFile, exists := params["swapFile"]; exists && swapFile != "" {
		cfg.swapFile = cfgParser.GetValueWithDefault(swapFile, defaultSwapFile)
	} else {
		cfg.swapFile = defaultSwapFile
	}

	if swapDelay, exists := params["swapDelay"]; exists && swapDelay != "" {
		var errParseDuration error
		cfg.swapDelay, errParseDuration = time.ParseDuration(swapDelay)
		if errParseDuration != nil {
			cfg.swapDelay = defaultSwapDelay
		}
	} else {
		cfg.swapDelay = defaultSwapDelay
	}

	if holdDelay, exists := params["holdDelay"]; exists && holdDelay != "" {
		var errParseDuration error
		cfg.holdDoneDelay, errParseDuration = time.ParseDuration(holdDelay)
		if errParseDuration != nil {
			cfg.holdDoneDelay = defaultHoldDoneDelay
		}
	} else {
		cfg.holdDoneDelay = defaultHoldDoneDelay
	}

	if holdDoneTasks, exists := params["holdTasks"]; exists && holdDoneTasks != "" {
		var errParseDuration error
		cfg.holdDoneTasks, errParseDuration = time.ParseDuration(holdDoneTasks)
		if errParseDuration != nil {
			cfg.holdDoneTasks = defaultHoldDoneTasks
		}
	} else {
		cfg.holdDoneTasks = defaultHoldDoneTasks
	}

	return cfg, nil
}
