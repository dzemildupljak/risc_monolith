package auth_usecase

import (
	"log"
	"os"
)

type Configurations struct {
	AccessTokenPrivateKeyPath  string
	AccessTokenPublicKeyPath   string
	RefreshTokenPrivateKeyPath string
	RefreshTokenPublicKeyPath  string
	JwtExpiration              int
	JwtRefreshExpiration       int
	MailVerifTemplateID        string
	PassResetTemplateID        string
}

func newConfigurations() *Configurations {
	curDir, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}
	return &Configurations{
		AccessTokenPrivateKeyPath:  curDir + "/access-private.pem",
		AccessTokenPublicKeyPath:   curDir + "/access-public.pem",
		RefreshTokenPrivateKeyPath: curDir + "/refresh-private.pem",
		RefreshTokenPublicKeyPath:  curDir + "/refresh-public.pem",
		JwtExpiration:              1,
		JwtRefreshExpiration:       1,
	}
}
