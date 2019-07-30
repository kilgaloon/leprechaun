package notifier

import (
	"github.com/kilgaloon/leprechaun/log"
	notis "github.com/kilgaloon/leprechaun/notifier/notifications"
)

//Config interface specifis what methods we need to build Notifier struct
type Config interface {
	notis.EmailConfig
}

// Notifier holds all informations for notifing users
type Notifier struct {
	Methods map[string]notis.Notification
	log.Logs
}

// NotifyWithOptions users
func (n Notifier) NotifyWithOptions(o notis.Options) {
	for _, notif := range n.Methods {
		// send notifications
		notif.SetMessage(o.GetBody())
		notif.SetTitle(o.GetTitle())
		err := notif.Send()
		if err != nil {
			n.Error(err.Error())
		}
	}
}

// New create new notifier to send notifications about workers and jobs
func New(cfg Config, log log.Logs) *Notifier {
	n := &Notifier{
		Methods: make(map[string]notis.Notification),
		Logs:    log,
	}

	if cfg.SMTPHost() != "" {
		n.Methods["mail"] = notis.NewEmail(cfg)
	}

	return n
}
