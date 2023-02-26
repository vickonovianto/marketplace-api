package delivery

import (
	"errors"
	"marketplace-api/helper"
	"marketplace-api/model"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type tokoDelivery struct {
	tokoUsecase model.TokoUsecase
}

type TokoDelivery interface {
	MountProtectedRoutes(jwtMiddleware func(*fiber.Ctx) error, group fiber.Router)
}

func NewTokoDelivery(tokoUsecase model.TokoUsecase) TokoDelivery {
	return &tokoDelivery{tokoUsecase: tokoUsecase}
}

func (p *tokoDelivery) MountProtectedRoutes(jwtMiddleware func(*fiber.Ctx) error, group fiber.Router) {
	group.Get("", jwtMiddleware, p.FetchAndPaginateTokoHandler)
	group.Get("/my", jwtMiddleware, p.MyTokoHandler)
	group.Get("/:id_toko", jwtMiddleware, p.DetailTokoHandler)
	group.Put("/:id_toko", jwtMiddleware, p.EditTokoHandler)
}

func (p *tokoDelivery) FetchAndPaginateTokoHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	req := new(model.TokoFetchPaginateRequest)
	limitString := strings.TrimSpace(c.Query("limit"))
	pageString := strings.TrimSpace(c.Query("page"))
	namaString := strings.TrimSpace(c.Query("nama"))
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
	req.Nama = namaString
	tokoFetchPaginateResponse, err := p.tokoUsecase.FetchAndPaginateToko(ctx, req)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusInternalServerError, err)
	}
	return helper.ResponseSuccessJson(c, tokoFetchPaginateResponse)
}

func (p *tokoDelivery) DetailTokoHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	idTokoString := c.Params("id_toko")
	idTokoInt, err := strconv.Atoi(idTokoString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid id toko"))
	}
	tokoGetByIDResponse, err := p.tokoUsecase.GetTokoByID(ctx, idTokoInt)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	return helper.ResponseSuccessJson(c, tokoGetByIDResponse)
}

func (p *tokoDelivery) MyTokoHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	userId, err := helper.GetUserIdFromToken(c)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	getMyTokoResponse, err := p.tokoUsecase.GetMyToko(ctx, userId)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	return helper.ResponseSuccessJson(c, getMyTokoResponse)
}

func (p *tokoDelivery) EditTokoHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	var req model.TokoUpdateRequest

	idTokoString := c.Params("id_toko")
	idTokoInt, err := strconv.Atoi(idTokoString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid id"))
	}
	req.ID = idTokoInt

	userId, err := helper.GetUserIdFromToken(c)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	req.IdUser = userId

	namaToko := strings.TrimSpace(c.FormValue("nama_toko"))
	if namaToko == "" {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("nama toko must not be empty"))
	}
	if len(namaToko) > 255 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("nama toko must not exceed 255 characters"))
	}
	req.NamaToko = namaToko

	photo, err := c.FormFile("photo")
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	fileExtension := photo.Filename[strings.LastIndex(photo.Filename, ".")+1:]
	if fileExtension != "jpg" &&
		fileExtension != "jpeg" &&
		fileExtension != "png" &&
		fileExtension != "webp" &&
		fileExtension != "jfif" {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("photo must be jpg/jpeg/png/webp/jfif"))
	}
	photoFilename := idTokoString + "." + fileExtension
	rootFolderPath, err := filepath.Abs("./")
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusInternalServerError, err)
	}
	photoFilePath := filepath.Join(rootFolderPath, "uploads", "toko", photoFilename)
	urlFoto := c.BaseURL() + "/uploads/toko/" + photoFilename
	req.UrlFoto = urlFoto

	err = c.SaveFile(photo, photoFilePath)
	if err != nil {
		return helper.ResponseErrorJson(c, http.StatusInternalServerError, err)
	}

	tokoUpdateResponse, err := p.tokoUsecase.EditToko(ctx, &req)
	if err != nil {
		errRemove := os.Remove(photoFilePath)
		if errRemove != nil {
			return helper.ResponseErrorJson(c, http.StatusInternalServerError, errRemove)
		}
		return helper.ResponseErrorJson(c, http.StatusBadRequest, err)
	}

	return helper.ResponseSuccessJson(c, tokoUpdateResponse)
}
