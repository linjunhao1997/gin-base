package email

import (
	"gin-base/pkg/logging"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"strconv"
	"strings"
)

var logger = logging.GetLogger("email")

type Client struct {
	*gomail.Dialer
}

func (cli *Client) NewClient(host, username, password string) {
	parts := strings.Split(host, ":")
	if len(parts) < 2 {
		panic("host error")
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}

	cli.Dialer = gomail.NewDialer(parts[0], port, username, password)
}

func (cli *Client) SendHtml(to, cc, bcc []string, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", cli.Username)
	msg.SetHeader("To", to...)
	msg.SetHeader("Cc", cc...)
	msg.SetHeader("Bcc", bcc...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)
	logger.WithFields(log.Fields{
		"From": cli.Username,
		"To":   to,
		"Cc":   cc,
		"Bcc":  bcc,
	})
	err := cli.DialAndSend(msg)
	if err != nil {
		logger.Warn(err)
	}
	return err
}

func (cli *Client) SendHtmlWithAttachments(to, cc, bcc []string, subject, body string, attachments []*Attachment) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", cli.Username)
	msg.SetHeader("To", to...)
	msg.SetHeader("Cc", cc...)
	msg.SetHeader("Bcc", bcc...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)
	for _, attachment := range attachments {
		settings := append(attachment.ReadSetting, FilenameSetting(attachment.Filename))
		msg.Attach(attachment.Filename, settings...)
	}
	logger.WithFields(log.Fields{
		"From":    cli.Username,
		"To":      to,
		"Cc":      cc,
		"Bcc":     bcc,
		"fileLen": len(attachments),
	})
	err := cli.DialAndSend(msg)
	if err != nil {
		logger.Warn(err)
	}
	return err
}
