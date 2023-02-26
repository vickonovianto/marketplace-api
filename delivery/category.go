package delivery

import (
	"errors"
	"marketplace-api/helper"
	"marketplace-api/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type categoryDelivery struct {
	categoryUsecase model.CategoryUsecase
}

type CategoryDelivery interface {
	MountUnprotectedRoutes(group fiber.Router)
	MountProtectedRoutes(jwtMiddleware func(*fiber.Ctx) error, group fiber.Router)
}

func NewCategoryDelivery(categoryUsecase model.CategoryUsecase) CategoryDelivery {
	return &categoryDelivery{categoryUsecase: categoryUsecase}
}

func (p *categoryDelivery) MountUnprotectedRoutes(group fiber.Router) {
	group.Get("", p.FetchCategoryHandler)
}

func (p *categoryDelivery) MountProtectedRoutes(jwtMiddleware func(*fiber.Ctx) error, group fiber.Router) {
	group.Post("", jwtMiddleware, helper.CheckAdminTokenHandler, p.StoreCategoryHandler)
	group.Get("/:id", jwtMiddleware, helper.CheckAdminTokenHandler, p.DetailCategoryHandler)
	group.Put("/:id", jwtMiddleware, helper.CheckAdminTokenHandler, p.EditCategoryHandler)
	group.Delete("/:id", jwtMiddleware, helper.CheckAdminTokenHandler, p.DeleteCategoryHandler)
}

func (p *categoryDelivery) StoreCategoryHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	var req model.CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	req.NamaCategory = strings.TrimSpace(req.NamaCategory)
	if req.NamaCategory == "" {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("nama category must not be empty"))
	}
	if len(req.NamaCategory) > 255 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("nama category cannot exceed 255 characters"))
	}

	categoryResponse, err := p.categoryUsecase.StoreCategory(ctx, &req)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusInternalServerError, err)
	}
	return helper.ResponseSuccessJson(c, categoryResponse)
}

func (p *categoryDelivery) FetchCategoryHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	categoryResponses, err := p.categoryUsecase.FetchAllCategory(ctx)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusInternalServerError, err)
	}
	return helper.ResponseSuccessJson(c, categoryResponses)
}

func (p *categoryDelivery) DetailCategoryHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	idString := c.Params("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid id"))
	}
	categoryResponse, err := p.categoryUsecase.GetCategoryByID(ctx, idInt)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	return helper.ResponseSuccessJson(c, categoryResponse)
}

func (p *categoryDelivery) EditCategoryHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	var req model.CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	req.NamaCategory = strings.TrimSpace(req.NamaCategory)
	if req.NamaCategory == "" {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("nama category must not be empty"))
	}
	if len(req.NamaCategory) > 255 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("nama category cannot exceed 255 characters"))
	}

	idString := c.Params("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid id"))
	}

	categoryResponse, err := p.categoryUsecase.EditCategory(ctx, idInt, &req)
	if err != nil {
		return helper.ResponseErrorJson(c, http.StatusBadRequest, err)
	}
	return helper.ResponseSuccessJson(c, categoryResponse)
}

func (p *categoryDelivery) DeleteCategoryHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	idString := c.Params("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid id"))
	}
	err = p.categoryUsecase.DestroyCategory(ctx, idInt)
	if err != nil {
		return helper.ResponseErrorJson(c, http.StatusBadRequest, err)
	}
	return helper.ResponseSuccessJson(c, "")
}
