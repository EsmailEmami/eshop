package email

type Key int

const (
	KeyForgotPassword Key = iota
)

func (k Key) getTemplatePath() string {
	switch k {
	case KeyForgotPassword:
		return "./templates/forgot-password.html"
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
