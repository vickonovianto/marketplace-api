package delivery

import (
	"errors"
	"marketplace-api/helper"
	"marketplace-api/model"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type trxDelivery struct {
	trxUsecase model.TrxUsecase
}

type TrxDelivery interface {
	MountProtectedRoutes(jwtMiddleware func(*fiber.Ctx) error, group fiber.Router)
}

func NewTrxDelivery(trxUsecase model.TrxUsecase) TrxDelivery {
	return &trxDelivery{trxUsecase: trxUsecase}
}

func (p *trxDelivery) MountProtectedRoutes(jwtMiddleware func(*fiber.Ctx) error, group fiber.Router) {
	group.Post("", jwtMiddleware, p.StoreTrxHandler)
	group.Get("", jwtMiddleware, p.FetchTrxHandler)
	group.Get("/:id", jwtMiddleware, p.GetTrxByIDHandler)
}

func (p *trxDelivery) StoreTrxHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	var req model.TrxStoreRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}

	req.MethodBayar = strings.TrimSpace(req.MethodBayar)
	if req.MethodBayar == "" {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("method bayar must not be empty"))
	}
	if len(req.MethodBayar) > 255 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("method bayar must not exceed 255"))
	}

	if req.AlamatPengiriman <= 0 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid alamat kirim"))
	}

	for _, detailTrxRequest := range req.DetailTrxRequests {
		produkId := detailTrxRequest.ProductId
		if produkId <= 0 {
			return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid product id"))
		}

		kuantitas := detailTrxRequest.Kuantitas
		if kuantitas <= 0 {
			return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid kuantitas"))
		}
	}

	userId, err := helper.GetUserIdFromToken(c)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}

	trxGetByIDResponse, err := p.trxUsecase.StoreTrx(ctx, &req, userId)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}

	return helper.ResponseSuccessJson(c, trxGetByIDResponse)
}

func (p *trxDelivery) FetchTrxHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	req := new(model.TrxFetchRequest)

	namaProdukSearch := strings.TrimSpace(c.Query("search"))
	if len(namaProdukSearch) > 255 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("nama produk cannot exceed 255 characters"))
	}
	req.Search = namaProdukSearch

	limitString := strings.TrimSpace(c.Query("limit"))
	pageString := strings.TrimSpace(c.Query("page"))
	limitInt, pageInt := -1, 1
	var err error
	if limitString != "" {
		limitInt, err = strconv.Atoi(limitString)
		if err != nil {
			return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("limit must be integer"))
		}
		if limitInt < 1 {
			return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("limit must be greater than zero"))
		}
	}
	if pageString != "" {
		if limitString == "" {
			return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("must use limit query param when using page query param"))
		}
		pageInt, err = strconv.Atoi(pageString)
		if err != nil {
			return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("page must be integer"))
		}
		if pageInt < 1 {
			return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("page must be greater than zero"))
		}
	}
	req.Limit = limitInt
	req.Page = pageInt

	userId, err := helper.GetUserIdFromToken(c)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}

	trxFetchResponse, err := p.trxUsecase.FetchTrx(ctx, req, userId)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	return helper.ResponseSuccessJson(c, trxFetchResponse)
}

func (p *trxDelivery) GetTrxByIDHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	trxIdString := c.Params("id")
	trxIdInt, err := strconv.Atoi(trxIdString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid id"))
	}

	userId, err := helper.GetUserIdFromToken(c)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}

	trxGetByIdResponse, err := p.trxUsecase.GetTrxByID(ctx, trxIdInt, userId)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	return helper.ResponseSuccessJson(c, trxGetByIdResponse)
}
