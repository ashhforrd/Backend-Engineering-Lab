package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ashhforrd/backend-engineering-lab/05-caching-and-performance/cache-aside/internal/product"
)

type Handler struct {
	service *product.Service
}

type UpdateProductRequest struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
	Stock int    `json:"stock"`
}

func NewHandler(
	service *product.Service,
) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetProduct(
	writer http.ResponseWriter,
	request *http.Request,
) {
	productID, err := parseProductID(request)
	if err != nil {
		writeError(
			writer,
			http.StatusBadRequest,
			"invalid product id",
		)
		return
	}

	result, err := h.service.Get(
		request.Context(),
		productID,
	)
	if err != nil {
		writeServiceError(writer, err)
		return
	}

	writeJSON(writer, http.StatusOK, result)
}

func (h *Handler) UpdateProduct(
	writer http.ResponseWriter,
	request *http.Request,
) {
	productID, err := parseProductID(request)
	if err != nil {
		writeError(
			writer,
			http.StatusBadRequest,
			"invalid product id",
		)
		return
	}

	var input UpdateProductRequest

	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeError(
			writer,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	if input.Name == "" ||
		input.Price < 0 ||
		input.Stock < 0 {
		writeError(
			writer,
			http.StatusBadRequest,
			"invalid product data",
		)
		return
	}

	updatedProduct, err := h.service.Update(
		request.Context(),
		product.Product{
			ID:    productID,
			Name:  input.Name,
			Price: input.Price,
			Stock: input.Stock,
		},
	)
	if err != nil {
		writeServiceError(writer, err)
		return
	}

	writeJSON(
		writer,
		http.StatusOK,
		updatedProduct,
	)
}

func parseProductID(
	request *http.Request,
) (int64, error) {
	return strconv.ParseInt(
		request.PathValue("id"),
		10,
		64,
	)
}

func writeServiceError(
	writer http.ResponseWriter,
	err error,
) {
	if errors.Is(err, product.ErrNotFound) {
		writeError(
			writer,
			http.StatusNotFound,
			"product not found",
		)
		return
	}

	writeError(
		writer,
		http.StatusInternalServerError,
		"internal server error",
	)
}

func writeError(
	writer http.ResponseWriter,
	statusCode int,
	message string,
) {
	writeJSON(
		writer,
		statusCode,
		map[string]string{
			"error": message,
		},
	)
}

func writeJSON(
	writer http.ResponseWriter,
	statusCode int,
	value any,
) {
	writer.Header().Set(
		"Content-Type",
		"application/json",
	)
	writer.WriteHeader(statusCode)

	_ = json.NewEncoder(writer).Encode(value)
}
