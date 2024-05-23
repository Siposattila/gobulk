package validation

import (
	"sync"

	"github.com/Siposattila/gobulk/internal/logger"
	"github.com/schollz/progressbar/v3"
)

type MasterInterface interface {
	NewWork(fn func())
	Start()
	Stop()
	Wait()
}

type master struct {
	maxWorkers int
	pending    chan *work
	stop       chan bool
	bar        *progressbar.ProgressBar
	wg         sync.WaitGroup
}

type work struct {
	fn func()
}

func newMaster(totalWork int64, maxWorkers int) MasterInterface {
	return &master{
		maxWorkers: maxWorkers,
		pending:    make(chan *work),
		stop:       make(chan bool),
		bar:        progressbar.Default(totalWork),
	}
}

func (m *master) NewWork(fn func()) {
	m.wg.Add(1)
	m.pending <- &work{
		fn: fn,
	}
}

func (m *master) worker() {
	for {
		select {
		case <-m.stop:
			logger.LogWarning("Stop signal received stopping worker.")

			return
		default:
			w := <-m.pending
			m.bar.Add(1)
			w.fn()
			m.wg.Done()
		}
	}
}

func (m *master) Start() {
	for i := 0; i < m.maxWorkers; i++ {
		go m.worker()
		logger.LogSuccess("Worker started.")
	}
}

func (m *master) Stop() {
	for i := 0; i < m.maxWorkers; i++ {
		m.stop <- true
	}
}

func (m *master) Wait() {
	m.wg.Wait()
}
