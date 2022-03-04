package loadbalancer

import (
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type AtomicCounter struct {
	value, minValue, maxValue int
	mu                        sync.Mutex
}

func NewAtomicCounter(minValue int, maxValue int) *AtomicCounter {
	return &AtomicCounter{
		value:    minValue,
		minValue: minValue,
		maxValue: maxValue,
	}
}

func (ac *AtomicCounter) Inc() int {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	oldValue := ac.value
	ac.value++
	if ac.value > ac.maxValue {
		ac.value = ac.minValue
	}
	return oldValue
}

func (ac *AtomicCounter) Value() int {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	return ac.value
}

type LoadBalancer struct {
	workerAddrs []string
	counter     *AtomicCounter
	client      http.Client
}

func NewLoadBalancer(workerAddrs []string) *LoadBalancer {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = len(workerAddrs)
	t.MaxConnsPerHost = len(workerAddrs)
	t.MaxIdleConnsPerHost = len(workerAddrs)
	t.IdleConnTimeout = 90 * time.Second
	lb := LoadBalancer{
		workerAddrs: make([]string, len(workerAddrs)),
		counter:     NewAtomicCounter(0, len(workerAddrs)-1),
		client: http.Client{
			Transport: t,
		},
	}
	copy(lb.workerAddrs, workerAddrs)
	return &lb
}

func handleError(err error, w http.ResponseWriter) {
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (lb *LoadBalancer) SendRequest(w http.ResponseWriter, r *http.Request) {
	apiMethod := r.URL.Path
	adaptedRequest, err := http.NewRequest(r.Method, lb.workerAddrs[lb.counter.Inc()]+apiMethod, r.Body)
	if err != nil {
		handleError(err, w)
		return
	}
	r.Body.Close()
	response, err := lb.client.Do(adaptedRequest)
	if err != nil {
		handleError(err, w)
		return
	}
	w.Header().Set("Content-Type", response.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", response.Header.Get("Content-Length"))
	w.WriteHeader(response.StatusCode)
	io.Copy(w, response.Body)
	response.Body.Close()
}
