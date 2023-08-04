package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type GmailDriver struct {
}

func (g GmailDriver) getSmtpServer() string {
	return viper.GetString("smtp.server")
}

func (g GmailDriver) getPost() string {
	return viper.GetString("smtp.port")
}

func (g GmailDriver) getEmail() string {
	return viper.GetString("smtp.email")
}

func (g GmailDriver) getPassword() string {
	return viper.GetString("smtp.password")
}

func (g GmailDriver) GetDriverName() string {
	return "gmail"
}

func (g GmailDriver) Send(to []string, key Key, data any) error {

	// validate data model
	if !key.validateData(data) {
		return errors.New("invalid email data entered")
	}

	var (
		smtpServer = g.getSmtpServer()
		smtpPort   = g.getPost()
		email      = g.getEmail()
		password   = g.getPassword()
	)

	bts, err := generateMessage(key, data, &emailHeader{
		From:    email,
		To:      strings.Join(to, ", "),
		Subject: key.GetSubject(),
	})

	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", email, password, smtpServer)
	conn, err := smtp.Dial(smtpServer + ":" + smtpPort)
	if err != nil {
		return fmt.Errorf("error while dialing the smtp: %s", err.Error())
	}
	defer conn.Close()

	tlsConfig := &tls.Config{
		ServerName: smtpServer,
	}
	if err = conn.StartTLS(tlsConfig); err != nil {
		return errors.New("cannot start TLS")
	}
	if err = conn.Auth(auth); err != nil {
		return fmt.Errorf("error while configure auth: %s", err.Error())
	}

	if err = conn.Mail(email); err != nil {
		return fmt.Errorf("error while configure sender: %s", err.Error())
	}

	var wg sync.WaitGroup
	var errCh = make(chan error, len(to))
	for _, rcpt := range to {
		wg.Add(1)
		go func(rcpt string) {
			defer wg.Done()
			if err = conn.Rcpt(rcpt); err != nil {
				errCh <- fmt.Errorf("error while configure receiver: %s", err.Error())
			}
		}(rcpt)
	}
	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
	}

	wc, err := conn.Data()
	if err != nil {
		return fmt.Errorf("error while sending email: %s", err.Error())
	}
	defer wc.Close()

	_, err = wc.Write(bts)
	if err != nil {
		return fmt.Errorf("error while sending email: %s", err.Error())
	}

	return nil
}
