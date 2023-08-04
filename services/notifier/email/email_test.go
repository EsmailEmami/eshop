package email

import (
	"testing"

	"github.com/spf13/viper"
)

func TestGmailSendEmail(t *testing.T) {

	viper.SetConfigName("config")
	viper.AddConfigPath("../../../.")
	viper.SetConfigType("toml")
	viper.AutomaticEnv()

	viper.ReadInConfig()

	notifier := NewEmailNotifier("gmail")
	err := notifier.Send(
		[]string{"esmailemami84@gmail.com", "alireza83safarii@gmail.com"},
		KeyForgotPassword,
		ForgotPassword{
			Username:    "esmailemami",
			RecoveryUrl: "https://google.com",
		},
	)

	if err != nil {
		t.Error(err.Error())
	}
}
