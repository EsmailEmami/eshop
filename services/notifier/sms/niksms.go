package sms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

type NikSmsDriver struct {
}

func (d NikSmsDriver) getBaseURL() string {
	return "https://niksms.com/fa"
}

func (d NikSmsDriver) ptpSmsURL() string {
	return d.getBaseURL() + "/publicapi/PtpSms"
}

func (d NikSmsDriver) getUsername() string {
	return viper.GetString("niksms.username")
}

func (d NikSmsDriver) getPassword() string {
	return viper.GetString("niksms.password")
}

func (d NikSmsDriver) getSenderNumber() string {
	return viper.GetString("niksms.sender-number")
}

func (d NikSmsDriver) GetDriverName() string {
	return "niksms"
}

func (d NikSmsDriver) Send(to, message string) (result interface{}, err error) {
	reqModel := nikSmsSendModel{
		Username:     d.getUsername(),
		Password:     d.getPassword(),
		SenderNumber: d.getSenderNumber(),
		Numbers:      to,
		Message:      message,
	}

	bts, _ := json.Marshal(&reqModel)

	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, d.ptpSmsURL(), bytes.NewReader(bts))
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bts, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("NikSmsDriver@Send, response:" + string(bts))
	}

	b, _ := io.ReadAll(res.Body)
	log.Println(string(b))

	respModel := struct {
		Status         int64   `json:"Status"`
		Id             *string `json:"Id"`
		WarningMessage *string `json:"WarningMessage"`
		NikIds         *string `json:"NikIds"`
		Data           *string `json:"Data"`
	}{}

	err = json.NewDecoder(res.Body).Decode(&respModel)
	if err != nil {
		return nil, err
	}

	if respModel.Status != 0 {
		return nil, fmt.Errorf("NikSmsDriver@Send response status: %v", res.Status)
	}
	return SmsResult{
		To:         to,
		TrackID:    *respModel.Data,
		DriverName: d.GetDriverName(),
	}, nil
}

type nikSmsSendModel struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	SenderNumber string `json:"senderNumber"`
	Numbers      string `json:"numbers"`
	Message      string `json:"message"`
}
