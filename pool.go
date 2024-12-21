package main

import "sync"

type taskFunc (func(interface{}) []Report)

type TaskPool struct {
	reportCh chan []Report
	done     chan interface{}
	wg       sync.WaitGroup
}

func NewTaskPool() *TaskPool {
	pool := TaskPool{
		done:     make(chan interface{}),
		reportCh: make(chan []Report, 100),
	}
	return &pool
}

func (tp *TaskPool) AddTask(param interface{}, f taskFunc) {
	tp.wg.Add(1)
	go func() {
		defer tp.wg.Done()
		d := f(param)
		tp.report(d)
	}()
}

func (tp *TaskPool) report(r []Report) {
	tp.reportCh <- r
}

func (tp *TaskPool) Run() {
	for {
		select {
		case ds := <-tp.reportCh:
			for _, d := range ds {
				reports.TryInsert(d)
			}
		case <-tp.done:
			return
		}
	}
}

func (tp *TaskPool) Stop() {
	close(tp.done)
	tp.wg.Wait()
}
