package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
	"github.com/recovery-flow/sso-oauth/internal/service/transport/responses"
)

func (h *Handlers) AdminSessionsGet(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "account_id"))
	if err != nil {
		h.Log.WithError(err).Warn("Invalid account_id")
		httpkit.RenderErr(w, problems.BadRequest(validation.Errors{
			"account_id": validation.NewError("account_id", "Invalid account_id"),
		})...)
		return
	}

	sessions, err := h.Domain.SessionsListByAccount(r.Context(), userID)
	if err != nil {
		httpkit.RenderErr(w, problems.InternalError())
		return
	}

	httpkit.Render(w, responses.SessionCollection(sessions))
}
