package responses

import (
	"github.com/recovery-flow/sso-oauth/internal/service/domain/core/models"
	"github.com/recovery-flow/sso-oauth/resources"
)

func Session(session models.Session) resources.Session {
	return resources.Session{
		Data: resources.SessionData{
			Type: resources.UserSessionType,
			Id:   session.ID.String(),
			Attributes: resources.SessionAttributes{
				UserId:    session.UserID.String(),
				Client:    session.Client,
				Ip:        session.IP,
				CreatedAt: session.CreatedAt,
				LastUsed:  session.LastUsed,
			},
		},
	}
}
