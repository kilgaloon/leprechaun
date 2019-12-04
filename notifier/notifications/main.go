package notifications

// Notification interface define what we need to send notification
type Notification interface {
	Messenger
	Send() error
}

// Options define available options to be used in function
// NotifyWithOptions
type Options struct {
	Title string
	Body  string
}

// Messenger is interface that will defined how to get message
// that needs to be used for body in notification
type Messenger interface {
	SetMessage(m string)
	GetMessage() string

	SetTitle(t string)
	GetTitle() string
}

// GetTitle returns title if its specified otherwise return default
func (o Options) GetTitle() string {
	if o.Title == "" {
		return "Leprechaun notification"
	}

	return o.Title
}

// GetBody returns body
func (o Options) GetBody() string {
	return o.Body
}
