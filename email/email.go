package email

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	gomail "gopkg.in/gomail.v2"
)

type EmailObject struct {
	Recipient string
	BodyText  string
}

func SendEmail(data *EmailObject) error {

	if data.Recipient == "" {
		return errors.New("email recipient can't be blank")
	}
	m := gomail.NewMessage()
	m.SetHeader("From", "FROM_EMAIL")
	m.SetHeader("To", data.Recipient)

	m.SetHeader("Subject", "MLS - Brand Ambition")
	m.SetBody(
		"text/html",
		fmt.Sprintf("<div style='font-size: 16px;'>%s</div>", data.BodyText),
	)

	d := gomail.NewDialer("EMAIL_DOMAIN", 465, "FROM_EMAIL", "EMAIL_PASSWORD")

	if err := d.DialAndSend(m); err != nil {
		glog.Errorf("Email to %s failed to send: %s", data.Recipient, err.Error())
		return err
	}

	return nil
}
