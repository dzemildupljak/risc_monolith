package mail_usecase

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/dzemildupljak/risc_monolith/server/usecase"
)

type Mail struct {
	Reciever  string
	MailTitle string
	MailBody  string
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

func (mi *MailInteractor) SendEmail(mail Mail) {
	// go func(mail Mail) {
	fmt.Println("============== mail ==============", mail)
	from := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	// from := "centarnitnp@gmail.com"
	// password := "centarnitnp11"
	// from := "risc_app@centarnit.com"
	// password := "risc_appmonolith"

	// Receiver email address.
	to := []string{mail.Reciever}

	// smtp server configuration.
	smtpHost := os.Getenv("MAIL_SMTP")
	smtpPort := os.Getenv("MAIL_PORT")
	// smtpHost := "smtp.gmail.com"
	// smtpPort := "587"
	// smtpHost := "mail.centarnit.com"
	// smtpPort := "587"

	fmt.Println("Message")

	// Message.
	message := []byte(mail.MailBody)

	fmt.Println("Authentication")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	fmt.Println("Sending email")
	fmt.Print("\n\n", smtpHost+":"+smtpPort, auth, from, to, message, "\n\n")

	// Sending email.
	go func() {
		err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
		if err != nil {
			// mi.Logger.LogError("Email Sent Failed:", err)
			fmt.Println("err", err)
		}
		// mi.Logger.LogAccess("Email Sent Successfully!")
		fmt.Println("Email Sent Successfully!")
	}()
	// }(mail)
}
