package mail_usecase

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"

	"github.com/dzemildupljak/risc_monolith/server/usecase"
)

var LinkToVerifyMail = "<p>Thanks for using our aplication</p><p>Confirm your account <a href=\"localhost:8080/verify/mail?email=dzemildupljak@mail.com&code=vaIujDpH&type=1\" target=\"_blank\">here</a></p>"

type templatedata struct {
	Name string
	URL  string
}

type Mail struct {
	Reciever  string
	MailTitle string
	MailBody  string
	Type      int
}

func NewMail() *Mail {
	return &Mail{}
}

type MailInteractor struct {
	Logger usecase.Logger
}

func NewMailInteractor(l usecase.Logger) *MailInteractor {
	return &MailInteractor{
		Logger: l,
	}
}

func (mi *MailInteractor) SendEmail(mail Mail, verifyCode, user_name string) {
	from := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")

	// Receiver email address.
	to := []string{mail.Reciever}

	// smtp server configuration.
	smtpHost := os.Getenv("MAIL_SMTP")
	smtpPort := os.Getenv("MAIL_PORT")
	host_adress := os.Getenv("HOST_ADRESS")

	fmt.Println("Authentication")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Generate template for email
	templateData := templatedata{}

	if mail.Type == 1 {

		templateData = templatedata{
			Name: user_name,
			URL:  host_adress + "/v1/verify/mail?email=" + mail.Reciever + "&code=" + verifyCode + "&type=" + fmt.Sprint(mail.Type),
		}
	} else {
		templateData = templatedata{
			Name: user_name,
			URL:  host_adress + "/password-reset?email=" + mail.Reciever + "&code=" + verifyCode + "&type=" + fmt.Sprint(mail.Type),
		}
	}
	t, err := template.ParseFiles("templates/email_vrification.html")

	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, templateData); err != nil {
		return
	}
	mailbody := buf.String()
	// Generate template for email

	fmt.Println("Sending email")

	var msg []byte

	if mail.Type == 1 {
		msg = []byte(
			"From: RISC Novi Pazar <" + from + ">\r\n" +
				"To: " + to[0] + "\r\n" +
				"Subject: RISC Novi Pazar - " + mail.MailTitle + "!\r\n" +
				"MIME: MIME-version: 1.0\r\n" +
				"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
				"\r\n" + mailbody)
	} else {
		templateData = templatedata{
			Name: user_name,
			URL:  host_adress + "/password-reset?email=" + mail.Reciever + "&code=" + verifyCode + "&type=" + fmt.Sprint(mail.Type),
		}
		msg = []byte(
			"From: RISC Novi Pazar <" + from + ">\r\n" +
				"To: " + to[0] + "\r\n" +
				"Subject: RISC Novi Pazar - " + mail.MailTitle + "!\r\n" +
				"MIME: MIME-version: 1.0\r\n" +
				"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
				"\r\n" + mailbody)
	}

	// Sending email.
	go func() {
		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(msg))
		if err != nil {
			fmt.Println("err", err)
		}
		fmt.Println("Email Sent Successfully!")
	}()

}
