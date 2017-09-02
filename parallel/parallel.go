package parallel

import (
	"sync"
)

type Parallelizer interface {
	Do(f func(interface{}, int) interface{}) []interface{}
}

type StringParallelizer struct {
	arr []string
}

func NewStringParallelizer(arr []string) Parallelizer {
	return StringParallelizer{
		arr: arr,
	}
}

func (p StringParallelizer) Do(f func(interface{}, int) interface{}) []interface{} {
	wait := sync.WaitGroup{}
	lock := sync.Mutex{}

	results := []interface{}{}

	for i, _ := range p.arr {
		wait.Add(1)
		v := p.arr[i]
		go func() {
			result := f(v, i)
			lock.Lock()
			results = append(results, result)
			lock.Unlock()
			wait.Done()
		}()
	}
	wait.Wait()

	return results
}
