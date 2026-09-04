package v2

import (
	"errors"
	"net/http"

	"github.com/ashhforrd/backend-engineering-lab/01-api-and-http/api-versioning/internal/httpapi"
	"github.com/ashhforrd/backend-engineering-lab/01-api-and-http/api-versioning/internal/user"
)

type Handler struct {
	service *user.Service
}

type UserResponse struct {
	ID        string          `json:"id"`
	Profile   ProfileResponse `json:"profile"`
	Email     string          `json:"email"`
	Status    string          `json:"status"`
	CreatedAt string          `json:"createdAt"`
}

type ProfileResponse struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
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
		"GET /api/v2/users/{id}",
		h.getUser,
	)
}

func (h *Handler) getUser(
	writer http.ResponseWriter,
	request *http.Request,
) {
	currentUser, err := h.service.GetByID(
		request.Context(),
		request.PathValue("id"),
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
		Profile: ProfileResponse{
			FirstName: currentUser.FirstName,
			LastName:  currentUser.LastName,
		},
		Email:  currentUser.Email,
		Status: string(currentUser.Status),
		CreatedAt: currentUser.CreatedAt.UTC().Format(
			"2006-01-02T15:04:05Z",
		),
	}
}
