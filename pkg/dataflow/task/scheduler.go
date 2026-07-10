package task

import (
	"errors"
	"sync"
	"time"

	"github.com/dellinger2023/net-flux/pkg/logger"
)

var (
	once     sync.Once
	instance TaskManager
)

type Task interface {
	ID() string
	Run() error
	Next() time.Time
}

type TaskManager interface {
	Start() error
	Stop()

	AddTask(task Task)
	RemoveTask(id string)
	GetTasks() []Task
}

type scheduler struct {
	mu      sync.RWMutex
	wg      sync.WaitGroup
	stop    chan struct{}
	started bool
	tasks   []Task
}

func (s *scheduler) AddTask(task Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks = append(s.tasks, task)
}

func (s *scheduler) GetTasks() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	tasks := make([]Task, len(s.tasks))
	copy(tasks, s.tasks)
	return tasks
}

func (s *scheduler) RemoveTask(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, t := range s.tasks {
		if t.ID() == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			break
		}
	}
}

func (s *scheduler) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.started {
		return errors.New("task manager already started")
	}

	s.stop = make(chan struct{})
	s.wg.Add(1)
	go s.run()
	s.started = true
	return nil
}

func (s *scheduler) Stop() {
	s.mu.Lock()
	if !s.started {
		s.mu.Unlock()
		return
	}
	stop := s.stop
	s.started = false
	s.mu.Unlock()

	close(stop)
	s.wg.Wait()
}

func (s *scheduler) run() {
	defer s.wg.Done()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.stop:
			logger.Infof("task manager stopped")
			return
		case <-ticker.C:
			now := time.Now()
			s.mu.RLock()
			tasks := make([]Task, len(s.tasks))
			copy(tasks, s.tasks)
			s.mu.RUnlock()

			for _, task := range tasks {
				next := task.Next()
				if !next.IsZero() && next.After(now) {
					continue
				}
				task.Run()
			}
		}
	}
}

func Manager() TaskManager {
	once.Do(func() {
		instance = &scheduler{}
	})
	return instance
}
