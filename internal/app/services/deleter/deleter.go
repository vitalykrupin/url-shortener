package ds

import (
	"sync"
)

type DeleteService struct {
	wg   *sync.WaitGroup
	done chan struct{}
}

func NewDeleteService() *DeleteService {
	return &DeleteService{
		wg:   &sync.WaitGroup{},
		done: make(chan struct{}),
	}
}

func (ds *DeleteService) Add(urls []string) <-chan []string {
	gen := ds.generator(urls)
	return ds.merge(gen)

	//go func() {
	//	for url := range out {
	//		ds.delete(url, userID)
	//	}
	//}()
}

func (ds *DeleteService) generator(urls []string) chan []string {
	ch := make(chan []string)
	ds.wg.Add(1)
	go func() {
		defer ds.wg.Done()
		select {
		case <-ds.done:
			close(ch)
			return
		default:
			ch <- urls
		}
	}()
	return ch
}

func (ds *DeleteService) merge(in ...<-chan []string) <-chan []string {
	out := make(chan []string)

	output := func(c <-chan []string) {
		for s := range c {
			out <- s
		}
		ds.wg.Done()
	}
	ds.wg.Add(len(in))
	for _, c := range in {
		go output(c)
	}
	go func() {
		ds.wg.Wait()
	}()
	return out
}
