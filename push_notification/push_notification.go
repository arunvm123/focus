package push

type Notification interface {
	SendPushNotification(registrationTokens []string, payload *Payload) error
}

type Payload struct {
	Title       string            `json:"title"`
	Body        string            `json:"body"`
	Data        map[string]string `json:"data"`
	ClickAction string            `json:"clickAction"`
}
