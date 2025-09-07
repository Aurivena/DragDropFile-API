package domain

import (
	"DragDrop-Files/internal/domain/entity"
	"sync"
)

type Pool struct {
	wg sync.WaitGroup
}

func (p *Pool) Run(n int, jobs <-chan entity.File, fn func(file *entity.File)) {
	p.wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer p.wg.Done()
			for j := range jobs {
				job := j
				fn(&job)
			}
		}()
	}
}

func (p *Pool) Wait() { p.wg.Wait() }
