package v1

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ashhforrd/backend-engineering-lab/01-api-and-http/api-versioning/internal/httpapi"
	"github.com/ashhforrd/backend-engineering-lab/01-api-and-http/api-versioning/internal/user"
)

type Handler struct {
	service *user.Service
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func NewHandler(service *user.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(
	mux *http.ServeMux,
) {
	mux.HandleFunc(
		"GET /api/v1/users/{id}",
		h.getUser,
	)
}

func (h *Handler) getUser(
	writer http.ResponseWriter,
	request *http.Request,
) {
	id := request.PathValue("id")

	addDeprecationHeaders(writer, id)

	currentUser, err := h.service.GetByID(
		request.Context(),
		id,
	)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidID):
			httpapi.WriteError(
				writer,
				http.StatusBadRequest,
				err.Error(),
			)
		case errors.Is(err, user.ErrNotFound):
			httpapi.WriteError(
				writer,
				http.StatusNotFound,
				err.Error(),
			)
		default:
			httpapi.WriteError(
				writer,
				http.StatusInternalServerError,
				"internal server error",
			)
		}

		return
	}

	httpapi.WriteJSON(
		writer,
		http.StatusOK,
		toResponse(currentUser),
	)
}

func toResponse(currentUser user.User) UserResponse {
	return UserResponse{
		ID: currentUser.ID,
		Name: strings.TrimSpace(
			currentUser.FirstName + " " +
				currentUser.LastName,
		),
		Email: currentUser.Email,
	}
}

func addDeprecationHeaders(
	writer http.ResponseWriter,
	userID string,
) {
	successorURL := fmt.Sprintf(
		"/api/v2/users/%s",
		url.PathEscape(userID),
	)

	writer.Header().Set(
		"Deprecation",
		"@1788480000",
	)
	writer.Header().Set(
		"Sunset",
		"Thu, 31 Dec 2026 23:59:59 GMT",
	)
	writer.Header().Set(
		"Link",
		fmt.Sprintf(
			"<%s>; rel=\"successor-version\"",
			successorURL,
		),
	)
}
