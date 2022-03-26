package mail_usecase

type MailUsecase interface {
	SendEmail(mail Mail, verifyCode, user_name string)
}
