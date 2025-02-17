package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/core/sectools"
	"github.com/recovery-flow/sso-oauth/internal/service/domain/core/sqlerr"
	"github.com/recovery-flow/sso-oauth/internal/service/transport/rest/requests"
	"github.com/recovery-flow/sso-oauth/internal/service/transport/rest/responses"
)

func (h *Handlers) Refresh(w http.ResponseWriter, r *http.Request) {
	svc := h.svc
	log := svc.Logger

	req, err := requests.NewRefresh(r)
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}
	refreshToken := req.Data.Attributes.RefreshToken

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Debugf("Missing Authorization header")
		httpkit.RenderErr(w, problems.Unauthorized("Missing Authorization header"))
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		log.Debugf("Invalid Authorization header format")
		httpkit.RenderErr(w, problems.Unauthorized("Invalid Authorization header format"))
		return
	}

	tokenString := parts[1]

	log.Debugf("Token received: %s", tokenString)

	userData, err := svc.TokenManager.VerifyJWT(r.Context(), tokenString)
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		log.Warnf("Token validation failed: %v", err)
		httpkit.RenderErr(w, problems.Unauthorized())
		return
	}
	if userData == nil {
		log.Debugf("Token validation failed")
		httpkit.RenderErr(w, problems.Unauthorized("Token validation failed"))
		return
	}

	userID, err := uuid.Parse(userData.ID)
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	sessionID, err := uuid.Parse(*userData.SessionID)
	if err != nil {
		httpkit.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	user, err := svc.DB.Accounts.GetByID(r, userID)
	if err != nil {
		render.RenderSelectErr(w, log, err, "Failed to get user")
		return
	}

	session, err := svc.DB.Sessions.GetByID(r, sessionID)
	if err != nil {
		render.RenderSelectErr(w, log, err, "Failed to get session")
		return
	}

	decryptedToken, err := tools.DecryptToken(session.Token, svc.Config.JWT.RefreshToken.EncryptionKey)
	if err != nil {
		log.Errorf("Failed to decrypt refresh token: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	if decryptedToken != refreshToken {
		svc.Logger.Warn("Provided refresh token does not match the stored token")
		httpkit.RenderErr(w, problems.Conflict())
		return
	}

	sesIDStr := sessionID.String()
	tokenAccess, err := svc.TokenManager.GenerateJWT(
		svc.Config.Server.Name,
		userID.String(),
		svc.Config.JWT.AccessToken.TokenLifetime,
		nil,
		&user.Role,
		&sesIDStr,
	)
	if err != nil {
		svc.Logger.Errorf("Error generating access token: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	tokenRefresh, err := svc.TokenManager.GenerateJWT(
		svc.Config.Server.Name,
		userID.String(),
		svc.Config.JWT.RefreshToken.TokenLifetime,
		nil,
		&user.Role,
		&sesIDStr,
	)
	if err != nil {
		svc.Logger.Errorf("Error generating refresh token: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	encryptedToken, err := tools.EncryptToken(tokenRefresh, svc.Config.JWT.RefreshToken.EncryptionKey)
	if err != nil {
		log.Errorf("Failed to encrypt refresh token: %v", err)
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	_, err = svc.DB.Sessions.UpdateToken(r, userID, encryptedToken)
	if err != nil {
		render.RenderSelectErr(w, log, err, "Failed to update session token")
		return
	}

	httpkit.Render(w, responses.TokensPair(tokenAccess, tokenRefresh))
}
