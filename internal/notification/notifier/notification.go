package notifier

type NotificationPriority int

const (
	Hight  NotificationPriority = 0
	Medium NotificationPriority = 1
	Low    NotificationPriority = 2
)

type Notification struct {
	Title     string
	Text      string
}
