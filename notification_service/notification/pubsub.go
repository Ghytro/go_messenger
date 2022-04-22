package notification

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type notificationStorageImpl struct {
	notifications  map[string][]Notification
	pubSubChannels map[string]chan struct{}
	mutex          sync.Mutex
}

var nStorage = notificationStorageImpl{
	notifications:  make(map[string][]Notification),
	mutex:          sync.Mutex{},
	pubSubChannels: make(map[string]chan struct{}),
}
var ErrNotificationTimeout = errors.New("timeout for notification wait time")
var ErrNoNotifications = errors.New("no notifications when calling get")

func Get(token string) ([]Notification, error) {
	nStorage.mutex.Lock()
	defer nStorage.mutex.Unlock()

	n, ok := nStorage.notifications[token]
	if !ok {
		return nil, ErrNoNotifications
	}
	nStorage.notifications[token] = nStorage.notifications[token][:0]
	return n, nil
}

func Push(token string, n Notification) {
	nStorage.mutex.Lock()
	defer nStorage.mutex.Unlock()
	fmt.Println("here1")
	if _, ok := nStorage.notifications[token]; !ok {
		nStorage.notifications[token] = make([]Notification, 0)
	}
	fmt.Println("here2")
	if _, ok := nStorage.pubSubChannels[token]; !ok {
		nStorage.pubSubChannels[token] = make(chan struct{}, 1)
	}
	fmt.Println("here3")
	nStorage.notifications[token] = append(nStorage.notifications[token], n)
	fmt.Println("here4")
	if len(nStorage.pubSubChannels[token]) == 0 { // only store info that notifications were appended
		nStorage.pubSubChannels[token] <- struct{}{}
	}
	fmt.Println("here5")
}

func WaitForNotifications(token string) ([]Notification, error) {
	if _, ok := nStorage.pubSubChannels[token]; !ok {
		nStorage.pubSubChannels[token] = make(chan struct{}, 1)
	}

	select {
	case <-time.After(time.Second * 20):
		return nil, ErrNotificationTimeout
	case <-nStorage.pubSubChannels[token]:
		return Get(token)
	}
}
