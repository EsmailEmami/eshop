package email

import "github.com/esmailemami/eshop/app/services/file"

type Key int

const (
	KeyForgotPassword Key = iota
)

func (k Key) getTemplatePath() string {
	path := file.GetPath("services", "notifier", "email", "templates")
	switch k {
	case KeyForgotPassword:
		return path + "/forgot-password.html"
	}

	panic("invalid template key!")
}

func (k Key) validateData(data any) (ok bool) {
	switch k {
	case KeyForgotPassword:
		_, ok = data.(ForgotPassword)
		return
	}

	panic("invalid template key!")
}

func (k Key) GetSubject() string {
	switch k {
	case KeyForgotPassword:
		return "Recovery Password"
	}

	panic("invalid template key!")
}
