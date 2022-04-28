package notification

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type notificationStorageImpl struct {
	notifications  map[int][]Notification
	pubSubChannels map[int]chan struct{}
	mutex          sync.Mutex
}

var nStorage = notificationStorageImpl{
	notifications:  make(map[int][]Notification),
	mutex:          sync.Mutex{},
	pubSubChannels: make(map[int]chan struct{}),
}
var ErrNotificationTimeout = errors.New("timeout for notification wait time")
var ErrNoNotifications = errors.New("no notifications when calling get")

func Get(userId int) ([]Notification, error) {
	nStorage.mutex.Lock()
	defer nStorage.mutex.Unlock()

	n, ok := nStorage.notifications[userId]
	if !ok {
		return nil, ErrNoNotifications
	}
	nStorage.notifications[userId] = nStorage.notifications[userId][:0]
	return n, nil
}

func Push(userId int, n Notification) {
	nStorage.mutex.Lock()
	defer nStorage.mutex.Unlock()
	fmt.Println("here1")
	if _, ok := nStorage.notifications[userId]; !ok {
		nStorage.notifications[userId] = make([]Notification, 0)
	}
	fmt.Println("here2")
	if _, ok := nStorage.pubSubChannels[userId]; !ok {
		nStorage.pubSubChannels[userId] = make(chan struct{}, 1)
	}
	fmt.Println("here3")
	nStorage.notifications[userId] = append(nStorage.notifications[userId], n)
	fmt.Println("here4")
	if len(nStorage.pubSubChannels[userId]) == 0 { // only store info that notifications were appended
		nStorage.pubSubChannels[userId] <- struct{}{}
	}
	fmt.Println("here5")
}

func WaitForNotifications(userId int) ([]Notification, error) {
	if _, ok := nStorage.pubSubChannels[userId]; !ok {
		nStorage.pubSubChannels[userId] = make(chan struct{}, 1)
	}

	select {
	case <-time.After(time.Second * 20):
		return nil, ErrNotificationTimeout
	case <-nStorage.pubSubChannels[userId]:
		return Get(userId)
	}
}
