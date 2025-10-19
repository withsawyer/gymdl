package utils

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	intervalTasks map[string]*intervalTask
	cronTasks     map[string]cron.EntryID
	delayTasks    map[string]*time.Timer
	cronRunner    *cron.Cron
	mu            sync.RWMutex
}

type intervalTask struct {
	Name     string
	Interval time.Duration
	Job      func()
	stopChan chan struct{}
}

// NewScheduler 创建调度器
func NewScheduler() *Scheduler {
	return &Scheduler{
		intervalTasks: make(map[string]*intervalTask),
		cronTasks:     make(map[string]cron.EntryID),
		delayTasks:    make(map[string]*time.Timer),
		cronRunner:    cron.New(cron.WithSeconds()), // 支持秒级 cron
	}
}

// AddIntervalTask 添加固定间隔任务
func (s *Scheduler) AddIntervalTask(name string, interval time.Duration, job func()) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.intervalTasks[name]; exists {
		return fmt.Errorf("任务 %s 已存在", name)
	}

	task := &intervalTask{
		Name:     name,
		Interval: interval,
		Job:      job,
		stopChan: make(chan struct{}),
	}

	s.intervalTasks[name] = task

	go func(t *intervalTask) {
		ticker := time.NewTicker(t.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.safeRun(t.Name, t.Job)
			case <-t.stopChan:
				return
			}
		}
	}(task)

	return nil
}

// AddCronTask 添加 Cron 表达式任务
func (s *Scheduler) AddCronTask(name string, expr string, job func()) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.cronTasks[name]; exists {
		return fmt.Errorf("任务 %s 已存在", name)
	}

	id, err := s.cronRunner.AddFunc(expr, func() { s.safeRun(name, job) })
	if err != nil {
		return err
	}

	s.cronTasks[name] = id
	return nil
}

// AddDelayTask 添加一次性延迟任务
func (s *Scheduler) AddDelayTask(name string, delay time.Duration, job func()) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.delayTasks[name]; exists {
		return fmt.Errorf("任务 %s 已存在", name)
	}

	timer := time.AfterFunc(delay, func() {
		s.safeRun(name, job)
		s.mu.Lock()
		delete(s.delayTasks, name)
		s.mu.Unlock()
	})

	s.delayTasks[name] = timer
	return nil
}

// safeRun 执行任务并捕获 panic，保证调度器稳定运行
func (s *Scheduler) safeRun(name string, job func()) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[任务 %s] panic: %v", name, r)
		}
	}()
	job()
}

// RemoveTask 删除任务（无论是哪种类型）
func (s *Scheduler) RemoveTask(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if t, ok := s.intervalTasks[name]; ok {
		close(t.stopChan)
		delete(s.intervalTasks, name)
	}

	if id, ok := s.cronTasks[name]; ok {
		s.cronRunner.Remove(id)
		delete(s.cronTasks, name)
	}

	if timer, ok := s.delayTasks[name]; ok {
		timer.Stop()
		delete(s.delayTasks, name)
	}
}

// Start 启动 Cron 任务调度器
func (s *Scheduler) Start() {
	s.cronRunner.Start()
}

// Stop 停止所有任务
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for name, t := range s.intervalTasks {
		close(t.stopChan)
		delete(s.intervalTasks, name)
	}

	for name := range s.cronTasks {
		delete(s.cronTasks, name)
	}
	s.cronRunner.Stop()

	for name, timer := range s.delayTasks {
		timer.Stop()
		delete(s.delayTasks, name)
	}
}
