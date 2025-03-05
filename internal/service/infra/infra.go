package infra

import (
	"github.com/recovery-flow/sso-oauth/internal/config"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/data"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/events"
	"github.com/recovery-flow/sso-oauth/internal/service/infra/jwtmanager"
	"github.com/sirupsen/logrus"
)

type Infra struct {
	Tokens jwtmanager.JWTManager
	Kafka  events.Kafka
	Data   *data.Data
}

func NewInfra(cfg *config.Config, log *logrus.Logger) (*Infra, error) {
	jwtManager := jwtmanager.NewJWTManager(cfg)
	eve := events.NewBroker(cfg, log)
	NewData, err := data.NewData(cfg, log)
	if err != nil {
		return nil, err
	}
	return &Infra{
		Tokens: jwtManager,
		Kafka:  eve,
		Data:   NewData,
	}, nil
}
