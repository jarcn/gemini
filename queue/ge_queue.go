package queue

import (
	"context"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"log"
	"sync"
)

type GeminiQueue struct {
	items []genai.GenerativeModel
	lock  sync.Mutex
	cond  *sync.Cond
}

func NewSafeQueue() *GeminiQueue {
	q := &GeminiQueue{}
	q.cond = sync.NewCond(&q.lock)
	return q
}

func (q *GeminiQueue) Enqueue(item genai.GenerativeModel) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.items = append(q.items, item)
	q.cond.Signal()
}

func (q *GeminiQueue) Dequeue() genai.GenerativeModel {
	q.lock.Lock()
	defer q.lock.Unlock()
	for len(q.items) == 0 {
		q.cond.Wait()
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *GeminiQueue) IsEmpty() bool {
	return len(q.items) == 0
}

func (q *GeminiQueue) Size() int {
	return len(q.items)
}

func InitQueue(keys []string, q *GeminiQueue) {
	for _, key := range keys {
		ctx := context.Background()
		client, err := genai.NewClient(ctx, option.WithAPIKey(key))
		if err != nil {
			log.Fatal(err)
		}
		model := client.GenerativeModel("gemini-pro")
		q.Enqueue(*model)
	}
	fmt.Printf("gemini queue init success,queue size is [%d]. \r\n", q.Size())
}
