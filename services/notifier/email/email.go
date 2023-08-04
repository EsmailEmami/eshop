package email

type EmailDriver interface {
	Send(to []string, key Key, data any) error
	GetDriverName() string
}

func NewEmailNotifier(driverName string) emailNotifier {
	switch driverName {
	case "gmail":
		{
			return emailNotifier{driver: GmailDriver{}}
		}
	}
	panic("invalid sms driver")
}

type emailNotifier struct {
	driver EmailDriver
}

func (e emailNotifier) Send(to []string, key Key, data any) error {
	return e.driver.Send(to, key, data)
}
func (e emailNotifier) GetDriverName() string {
	return "email"
}
