package service

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"

	"github.com/lestrrat-go/jwx/jwk"
	"go.uber.org/zap"
)

type KeyService interface {
	PublicKey() (res *jwk.Key, err error)
}

type keyService struct {
	logger *zap.Logger
}

func NewKeyService(log *zap.Logger) KeyService {
	return &userService{
		logger: log,
	}
}

func (svc *userService) PublicKey() (res *jwk.Key, err error) {
	publicKey, _ := ioutil.ReadFile(os.Getenv("AUTH_KEYS_PATH") + "public.pem")

	key, _ := pem.Decode(publicKey)

	pub, err := x509.ParsePKIXPublicKey(key.Bytes)
	if err != nil {
		return nil, err
	}

	s, err := jwk.New(pub)
	if err != nil {
		return nil, err
	}

	err = jwk.AssignKeyID(s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
