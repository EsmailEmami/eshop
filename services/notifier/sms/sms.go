package sms

type SmsDriver interface {
	Send(to, message string) (result interface{}, err error)
	GetDriverName() string
}

func NewSmsNotifier(driverName string) smsNotifier {
	switch driverName {
	case "niksms":
		{
			return smsNotifier{driver: NikSmsDriver{}}
		}
	}
	panic("invalid sms driver")
}

type smsNotifier struct {
	driver SmsDriver
}

func (n smsNotifier) Send(to, message string) (result interface{}, err error) {
	return n.driver.Send(to, message)
}

func (n smsNotifier) GetDriverName() string {
	return "sms"
}

type SmsResult struct {
	To         string `json:"to"`
	TrackID    string `json:"trackId"`
	DriverName string `json:"driverName"`
}
