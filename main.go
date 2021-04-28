package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var cache = map[int]Book{}
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func main() {
	wg := &sync.WaitGroup{}
	m := &sync.RWMutex{}
	cachech := make(chan Book)
	dbch := make(chan Book)
	for i := 0; i < 10; i++ {
		id := rnd.Intn(10) + 1
		wg.Add(2)
		go func(id int, wg *sync.WaitGroup, m *sync.RWMutex, cachech chan Book) {
			if b, ok := queryCache(id, m); ok {
				cachech <- b
			}
			wg.Done()
		}(id, wg, m, cachech)

		go func(id int, wg *sync.WaitGroup, m *sync.RWMutex, dbch chan Book) {
			if b, ok := queryDatabase(id, m); ok {
				dbch <- b
			}
			wg.Done()
		}(id, wg, m, dbch)

		go func(cachech, dbch chan Book) {
			select {
			case b := <-cachech:
				fmt.Println("From cache")
				fmt.Println(b)
				<-dbch
			case b := <-dbch:
				fmt.Println("From database")
				fmt.Println(b)
			}

		}(cachech, dbch)

		time.Sleep(150 * time.Millisecond)

	}

	wg.Wait()
}

func queryCache(id int, m *sync.RWMutex) (Book, bool) {
	m.RLock()
	b, ok := cache[id]
	m.RUnlock()
	return b, ok
}

func queryDatabase(id int, m *sync.RWMutex) (Book, bool) {
	time.Sleep(100 * time.Millisecond)
	for _, b := range books {
		if b.ID == id {
			m.Lock()
			cache[id] = b
			m.Unlock()
			return b, true
		}
	}
	return Book{}, false
}
