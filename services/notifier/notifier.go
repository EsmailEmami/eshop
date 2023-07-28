package notifier

import "github.com/esmailemami/eshop/services/notifier/sms"

type NorifierDriver interface {
	Send(to, message string) (result interface{}, err error)
	GetDriverName() string
}

func NewNorifier(driverName string) notifier {
	switch driverName {
	case "sms":
		{
			driver := sms.NewSmsNotifier("abfa")
			return notifier{
				dirver: driver,
			}
		}
	}

	panic("invalid notifier driver")
}

type notifier struct {
	dirver NorifierDriver
}

func (n notifier) Send(to, message string) (result interface{}, err error) {
	return n.dirver.Send(to, message)
}
