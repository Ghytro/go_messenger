package notification

type Notification struct {
	Action      string                 `json:"type"`
	Description map[string]interface{} `json:"description"`
	Timestamp   int64                  `json:"timestamp"`
}

type Notifications struct {
	List []Notification `json:"notifications"`
}
