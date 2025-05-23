package responses

import (
	"github.com/hs-zavet/sso-oauth/internal/app/models"
	"github.com/hs-zavet/sso-oauth/resources"
)

func SessionCollection(sessions []models.Session) resources.SessionsCollection {
	var data []resources.SessionData
	for _, session := range sessions {
		data = append(data, Session(session).Data)
	}
	return resources.SessionsCollection{
		Data: resources.SessionsCollectionData{
			Type: resources.AccountSessionsType,
			Attributes: resources.SessionsCollectionDataAttributes{
				Sessions: data,
			},
		},
	}
}
