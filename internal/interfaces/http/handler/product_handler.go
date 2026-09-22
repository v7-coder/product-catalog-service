package handler

import (
	"encoding/json"
	"net/http"

	"github.com/v7-coder/product-catalog-service/internal/application/product"
	httpdto "github.com/v7-coder/product-catalog-service/internal/interfaces/http/dto"
	"github.com/v7-coder/product-catalog-service/internal/pkg/app"
	"github.com/v7-coder/product-catalog-service/internal/pkg/validator"
)

type ProductHandler struct {
	createUC *product.CreateProductUseCase
}

func NewCreateProductHandler(createUC *product.CreateProductUseCase) *ProductHandler {
	return &ProductHandler{createUC: createUC}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var createRequest httpdto.CreateProductRequest

	// Декодируем JSON
	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		app.SendError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if err := validator.Validate(createRequest); err != nil {
		app.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	cmd := product.CreateProductCommand{
		Name:        createRequest.Name,
		Description: createRequest.Description,
		Price:       createRequest.Price,
		CategoryId:  createRequest.CategoryId,
	}

	result, err := h.createUC.Execute(r.Context(), cmd)
	if err != nil {
		app.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	app.SendSuccess(w, http.StatusCreated, result)
}
