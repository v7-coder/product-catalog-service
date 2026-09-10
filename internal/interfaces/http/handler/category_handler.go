package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/v7-coder/product-catalog-service/internal/application/category"
	httpdto "github.com/v7-coder/product-catalog-service/internal/interfaces/http/dto"
	"github.com/v7-coder/product-catalog-service/internal/pkg/app"
	"github.com/v7-coder/product-catalog-service/internal/pkg/validator"
)

type CategoryHandler struct {
	createUC  *category.CreateCategoryUseCase
	getUC     *category.GetCategoryUseCase
	getListUC *category.GetCategoriesUseCase
	updateUC  *category.UpdateCategoryUseCase
	deleteUC  *category.DeleteCategoryUseCase
}

func NewCategoryHandler(createUC *category.CreateCategoryUseCase, getUC *category.GetCategoryUseCase, getListUC *category.GetCategoriesUseCase, updateUC *category.UpdateCategoryUseCase, deleteUC *category.DeleteCategoryUseCase) *CategoryHandler {
	return &CategoryHandler{
		createUC:  createUC,
		getUC:     getUC,
		getListUC: getListUC,
		updateUC:  updateUC,
		deleteUC:  deleteUC,
	}
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var createRequest httpdto.CreateCategoryRequest

	// Декодируем JSON
	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		app.SendError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if err := validator.Validate(createRequest); err != nil {
		app.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	cmd := category.CreateCategoryCommand{
		Name:        createRequest.Name,
		Description: createRequest.Description,
	}

	result, err := h.createUC.Execute(r.Context(), cmd)
	if err != nil {
		app.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	app.SendSuccess(w, http.StatusCreated, result)
}

func (h *CategoryHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("categoryId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.SendError(w, http.StatusBadRequest, "invalid category ID")
		return
	}

	cmd := category.GetCategoryCommand{
		ID: id,
	}

	result, err := h.getUC.Execute(r.Context(), cmd)
	if err != nil {
		app.SendError(w, http.StatusNotFound, err.Error())
		return
	}

	app.SendSuccess(w, http.StatusOK, result)
}

func (h *CategoryHandler) GetList(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limitStr = "50"
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		app.SendError(w, http.StatusBadRequest, "invalid limit parameter")
		return
	}

	offsetStr := r.URL.Query().Get("offset")
	if offsetStr == "" {
		offsetStr = "0"
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		app.SendError(w, http.StatusBadRequest, "invalid offset parameter")
		return
	}

	cmd := category.GetCategoriesCommand{
		Limit:  limit,
		Offset: offset,
	}

	result, err := h.getListUC.Execute(r.Context(), cmd)
	if err != nil {
		app.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.SendSuccess(w, http.StatusOK, result)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("categoryId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.SendError(w, http.StatusBadRequest, "invalid category ID")
		return
	}

	defer r.Body.Close()

	var updateRequest httpdto.UpdateCategoryRequest

	// Декодируем JSON
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		app.SendError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}

	if err := validator.Validate(updateRequest); err != nil {
		app.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	cmd := category.UpdateCategoryCommand{
		ID:          id,
		Name:        updateRequest.Name,
		Description: updateRequest.Description,
	}

	result, err := h.updateUC.Execute(r.Context(), cmd)
	if err != nil {
		app.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	app.SendSuccess(w, http.StatusOK, result)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("categoryId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		app.SendError(w, http.StatusBadRequest, "invalid category ID")
		return
	}

	cmd := category.DeleteCategoryCommand{
		ID: id,
	}

	result, err := h.deleteUC.Execute(r.Context(), cmd)
	if err != nil {
		app.SendError(w, http.StatusNotFound, err.Error())
		return
	}

	app.SendSuccess(w, http.StatusOK, result)
}
