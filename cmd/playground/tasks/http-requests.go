package tasks

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type WorkerPool struct {
	count      int
	jobs      chan string
	results   chan string
	
	mu sync.Mutex
}

func NewWorkerPool(count, jobsSize, resultsSize int) WorkerPool {
	return WorkerPool{
		count: count,
		jobs: make(chan string, jobsSize),
		results: make(chan string, resultsSize),
	}
}

func (wp *WorkerPool) AddJob(job string, ctx context.Context) (error) {
	defer wp.mu.Unlock()
	wp.mu.Lock()
	
	select {
	case wp.jobs <- job:
	case <-ctx.Done():
		return nil
	}
	
	return nil
}

func processURL(url string) string {
	res, err := http.Get(url)

	if err != nil {
		return fmt.Sprintf("Address %s - failed \n", url)
	}
	defer res.Body.Close()
	
	_, _ = io.ReadAll(res.Body)
	
	if res.StatusCode != http.StatusOK {
		return  fmt.Sprintf("Address %s - failed \n", url)
	}
	
	return fmt.Sprintf("Address %s - ok \n", url)
}

func (wp *WorkerPool) Run(ctx context.Context) {
	wg := sync.WaitGroup{}
	
	wg.Add(wp.count)
	for range wp.count {
		go func() {
			defer wg.Done()
			for url := range wp.jobs {
				msg := processURL(url)

				select {
				case wp.results <- msg:
				case <- ctx.Done():
					return
				}
			}
		}()
	}

	wg.Wait()
	close(wp.results)
}

func HttpRequests() {
	ctx, cancel := context.WithCancel(context.Background())
	wp := NewWorkerPool(3, 2, 2)
	go wp.Run(ctx)

	urls := []string{
		"http://ozon.ru",
		"https://ozon.ru",
		"http://google.com",
		"http://non-existen.domain.ltd",
		"http://ya.ru",
		"https://ya.ru",
		"http://ёёёё",
	}

	urlsChan := make(chan string)
	go func() {
		for _, url := range urls {
			urlsChan <- url
		}
		close(urlsChan)
	}()

	go func() {
		for url := range urlsChan {
			wp.AddJob(url, ctx)
		}
	}()

	go func() {
		_ = <-time.After(10 * time.Millisecond)
		wp.mu.Lock()
		cancel()
		close(wp.jobs)
		wp.mu.Unlock()
	}()

	for res := range wp.results {
		fmt.Print("Result", res)
	}
	
	fmt.Println("END")
}
