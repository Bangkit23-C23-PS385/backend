package email

import (
	"bytes"
	"html/template"
	"os"
	"strconv"
	"strings"
	"ta/backend/src/constant"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

func LoadTemplate(templateType string, params map[string]interface{}) (loadedTemplate string, err error) {
	var tmpl *template.Template
	if templateType == string(constant.Register) {
		tmpl, err = template.New("register-html").Parse(Register())
	} else if templateType == string(constant.VerifySuccess) {
		tmpl, err = template.New("verifySuccess-html").Parse(VerifySuccess())
	}
	if err != nil {
		err = errors.Wrap(err, "template: parse html file")
		return
	}

	b := bytes.NewBufferString("")
	err = tmpl.Execute(b, params)
	if err != nil {
		err = errors.Wrap(err, "template: load params to template")
		return
	}

	loadedTemplate = b.String()
	if strings.Contains(loadedTemplate, "<no value>") {
		err = constant.ErrTemplateNoValue
	}
	return
}

func SendEmail(reciever, subject, content string) (err error) {
	_ = godotenv.Load()

	mailer := gomail.NewMessage()
	mailer.SetHeaders(map[string][]string{
		"From":    []string{mailer.FormatAddress(os.Getenv("EMAIL_AUTH_USERNAME"), os.Getenv("EMAIL_SENDER_NAME"))},
		"To":      []string{reciever},
		"Subject": []string{subject},
	})
	mailer.SetBody("text/html", content)

	port, _ := strconv.Atoi(os.Getenv("EMAIL_SMTP_PORT"))
	dialer := gomail.NewDialer(
		os.Getenv("EMAIL_SMTP_HOST"),
		port,
		os.Getenv("EMAIL_AUTH_USERNAME"),
		os.Getenv("EMAIL_AUTH_PASSWORD"),
	)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		err = errors.Wrap(err, "send email")
	}

	return
}
