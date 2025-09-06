package domain

import (
	"DragDrop-Files/internal/domain/entity"
	"sync"
)

type Pool struct {
	wg sync.WaitGroup
}

func (p *Pool) Work(jobs <-chan entity.File, processedFile []entity.File, fn func(file *entity.File, processedFiles []entity.File)) {
	for job := range jobs {
		defer p.wg.Done()
		fn(&job, processedFile)
	}
}

func (p *Pool) Add(count int) {
	p.wg.Add(count)
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
