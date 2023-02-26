package delivery

import (
	"errors"
	"marketplace-api/helper"
	"marketplace-api/model"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type produkDelivery struct {
	produkUsecase model.ProdukUsecase
}

type ProdukDelivery interface {
	MountUnprotectedRoutes(group fiber.Router)
	MountProtectedRoutes(jwtMiddleware func(*fiber.Ctx) error, group fiber.Router)
}

func NewProdukDelivery(produkUsecase model.ProdukUsecase) ProdukDelivery {
	return &produkDelivery{produkUsecase: produkUsecase}
}

func (p *produkDelivery) MountUnprotectedRoutes(group fiber.Router) {
	group.Get("", p.FetchProdukHandler)
	group.Get("/:id", p.DetailProdukHandler)
}

func (p *produkDelivery) MountProtectedRoutes(jwtMiddleware func(*fiber.Ctx) error, group fiber.Router) {
	group.Post("", jwtMiddleware, p.StoreProdukHandler)
	group.Put("/:id", jwtMiddleware, p.EditProdukHandler)
	group.Delete("/:id", jwtMiddleware, p.DeleteProdukHandler)
}

func (p *produkDelivery) StoreProdukHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	var req model.ProdukRequest

	form, err := c.MultipartForm()
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}

	if len(form.Value["nama_produk"]) == 0 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("nama produk must not be empty"))
	}
	namaProduk := strings.TrimSpace(form.Value["nama_produk"][0])
	if namaProduk == "" {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("nama produk must not be empty"))
	}
	if len(namaProduk) > 255 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("nama produk must not exceed 255 characters"))
	}
	req.NamaProduk = namaProduk

	if len(form.Value["category_id"]) == 0 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("category id must not be empty"))
	}
	categoryIdString := strings.TrimSpace(form.Value["category_id"][0])
	if categoryIdString == "" {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("category id must not be empty"))
	}
	categoryIdInt, err := strconv.Atoi(categoryIdString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid category id"))
	}
	req.IdCategory = categoryIdInt

	if len(form.Value["harga_reseller"]) == 0 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("harga reseller must not be empty"))
	}
	hargaReseller := strings.TrimSpace(form.Value["harga_reseller"][0])
	if hargaReseller == "" {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("harga reseller must not be empty"))
	}
	if len(hargaReseller) > 255 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("harga reseller must not exceed 255 characters"))
	}
	hargaResellerInt, err := strconv.Atoi(hargaReseller)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid harga reseller"))
	}
	if hargaResellerInt <= 0 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid harga reseller"))
	}
	req.HargaReseller = hargaReseller

	if len(form.Value["harga_konsumen"]) == 0 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("harga konsumen must not be empty"))
	}
	hargaKonsumen := strings.TrimSpace(form.Value["harga_konsumen"][0])
	if hargaKonsumen == "" {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("harga konsumen must not be empty"))
	}
	if len(hargaKonsumen) > 255 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("harga konsumen must not exceed 255 characters"))
	}
	hargaKonsumenInt, err := strconv.Atoi(hargaKonsumen)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid harga konsumen"))
	}
	if hargaKonsumenInt <= 0 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid harga konsumen"))
	}
	req.HargaKonsumen = hargaKonsumen

	if len(form.Value["stok"]) == 0 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("stok must not be empty"))
	}
	stokString := strings.TrimSpace(form.Value["stok"][0])
	if stokString == "" {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("stok must not be empty"))
	}
	stokInt, err := strconv.Atoi(stokString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid stok"))
	}
	if stokInt < 0 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid stok"))
	}
	req.Stok = stokInt

	if len(form.Value["deskripsi"]) == 0 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("deskripsi must not be empty"))
	}
	deskripsi := strings.TrimSpace(form.Value["deskripsi"][0])
	if deskripsi == "" {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("deskripsi must not be empty"))
	}
	req.Deskripsi = deskripsi

	if len(form.File["photos"]) == 0 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("photos must not be empty"))
	}
	photos := form.File["photos"]
	photoFilePaths := []string{}
	photoUrls := []string{}
	for _, photo := range photos {
		fileExtension := photo.Filename[strings.LastIndex(photo.Filename, ".")+1:]
		if fileExtension != "jpg" &&
			fileExtension != "jpeg" &&
			fileExtension != "png" &&
			fileExtension != "webp" &&
			fileExtension != "jfif" {
			return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("photo must be jpg/jpeg/png/webp/jfif"))
		}
		photoFilename := uuid.NewString() + "." + fileExtension
		rootFolderPath, err := filepath.Abs("./")
		if err != nil {
			return helper.ResponseErrorJson(c, fiber.StatusInternalServerError, err)
		}
		photoFilePath := filepath.Join(rootFolderPath, "uploads", "produk", photoFilename)
		photoUrl := c.BaseURL() + "/uploads/produk/" + photoFilename

		err = c.SaveFile(photo, photoFilePath)
		if err != nil {
			// when a file is failed to be saved on disk, delete already created photo on the disk
			for _, createdPhotoFilePath := range photoFilePaths {
				errRemove := os.Remove(createdPhotoFilePath)
				if errRemove != nil {
					return helper.ResponseErrorJson(c, http.StatusInternalServerError, errRemove)
				}
			}
			return helper.ResponseErrorJson(c, http.StatusInternalServerError, err)
		}
		photoFilePaths = append(photoFilePaths, photoFilePath)
		photoUrls = append(photoUrls, photoUrl)
	}
	req.PhotoUrls = photoUrls

	userId, err := helper.GetUserIdFromToken(c)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}

	produkResponse, err := p.produkUsecase.StoreProduk(ctx, &req, userId)
	if err != nil {
		// delete already created photos on the disk if the store produk operation is failed
		for _, createdPhotoFilePath := range photoFilePaths {
			errRemove := os.Remove(createdPhotoFilePath)
			if errRemove != nil {
				return helper.ResponseErrorJson(c, http.StatusInternalServerError, errRemove)
			}
		}
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}

	return helper.ResponseSuccessJson(c, produkResponse)
}

func (p *produkDelivery) FetchProdukHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	req := new(model.ProdukFetchRequest)

	namaProduk := strings.TrimSpace(c.Query("nama_produk"))
	if len(namaProduk) > 255 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("nama produk cannot exceed 255 characters"))
	}
	req.NamaProduk = namaProduk

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

	categoryIdString := strings.TrimSpace(c.Query("category_id"))
	if len(categoryIdString) > 255 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("category id cannot exceed 255 characters"))
	}
	categoryIdInt := -1
	if categoryIdString != "" {
		categoryIdInt, err = strconv.Atoi(categoryIdString)
		if err != nil {
			return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("category id must be integer"))
		}
	}
	req.CategoryId = categoryIdInt

	tokoIdString := strings.TrimSpace(c.Query("toko_id"))
	if len(tokoIdString) > 255 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("toko id cannot exceed 255 characters"))
	}
	tokoIdInt := -1
	if tokoIdString != "" {
		tokoIdInt, err = strconv.Atoi(tokoIdString)
		if err != nil {
			return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("toko id must be integer"))
		}
	}
	req.TokoId = tokoIdInt

	maxHargaString := strings.TrimSpace(c.Query("max_harga"))
	if len(maxHargaString) > 255 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("max harga cannot exceed 255 characters"))
	}
	maxHargaInt := -1
	if maxHargaString != "" {
		maxHargaInt, err = strconv.Atoi(maxHargaString)
		if err != nil {
			return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("max harga must be integer"))
		}
	}
	req.MaxHarga = maxHargaInt

	minHargaString := strings.TrimSpace(c.Query("min_harga"))
	if len(minHargaString) > 255 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("min harga cannot exceed 255 characters"))
	}
	minHargaInt := -1
	if minHargaString != "" {
		minHargaInt, err = strconv.Atoi(minHargaString)
		if err != nil {
			return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("min harga must be integer"))
		}
	}
	req.MinHarga = minHargaInt

	produkResponse, err := p.produkUsecase.FetchProduk(ctx, req)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusInternalServerError, err)
	}
	return helper.ResponseSuccessJson(c, produkResponse)
}

func (p *produkDelivery) DetailProdukHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	idString := c.Params("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid id"))
	}

	produkResponse, err := p.produkUsecase.GetProdukByID(ctx, idInt)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	return helper.ResponseSuccessJson(c, produkResponse)
}

func (p *produkDelivery) EditProdukHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	var req model.ProdukRequest

	produkIdString := c.Params("id")
	produkIdInt, err := strconv.Atoi(produkIdString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid id"))
	}

	userId, err := helper.GetUserIdFromToken(c)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}

	form, err := c.MultipartForm()
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}

	if len(form.Value["nama_produk"]) == 0 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("nama produk must not be empty"))
	}
	namaProduk := strings.TrimSpace(form.Value["nama_produk"][0])
	if namaProduk == "" {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("nama produk must not be empty"))
	}
	if len(namaProduk) > 255 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("nama produk must not exceed 255 characters"))
	}
	req.NamaProduk = namaProduk

	if len(form.Value["category_id"]) == 0 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("category id must not be empty"))
	}
	categoryIdString := strings.TrimSpace(form.Value["category_id"][0])
	if categoryIdString == "" {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("category id must not be empty"))
	}
	categoryIdInt, err := strconv.Atoi(categoryIdString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid category id"))
	}
	req.IdCategory = categoryIdInt

	if len(form.Value["harga_reseller"]) == 0 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("harga reseller must not be empty"))
	}
	hargaReseller := strings.TrimSpace(form.Value["harga_reseller"][0])
	if hargaReseller == "" {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("harga reseller must not be empty"))
	}
	if len(hargaReseller) > 255 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("harga reseller must not exceed 255 characters"))
	}
	hargaResellerInt, err := strconv.Atoi(hargaReseller)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid harga reseller"))
	}
	if hargaResellerInt <= 0 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid harga reseller"))
	}
	req.HargaReseller = hargaReseller

	if len(form.Value["harga_konsumen"]) == 0 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("harga konsumen must not be empty"))
	}
	hargaKonsumen := strings.TrimSpace(form.Value["harga_konsumen"][0])
	if hargaKonsumen == "" {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("harga konsumen must not be empty"))
	}
	if len(hargaKonsumen) > 255 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("harga konsumen must not exceed 255 characters"))
	}
	hargaKonsumenInt, err := strconv.Atoi(hargaKonsumen)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid harga konsumen"))
	}
	if hargaKonsumenInt <= 0 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid harga konsumen"))
	}
	req.HargaKonsumen = hargaKonsumen

	if len(form.Value["stok"]) == 0 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("stok must not be empty"))
	}
	stokString := strings.TrimSpace(form.Value["stok"][0])
	if stokString == "" {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("stok must not be empty"))
	}
	stokInt, err := strconv.Atoi(stokString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid stok"))
	}
	if stokInt < 0 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid stok"))
	}
	req.Stok = stokInt

	if len(form.Value["deskripsi"]) == 0 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("deskripsi must not be empty"))
	}
	deskripsi := strings.TrimSpace(form.Value["deskripsi"][0])
	if deskripsi == "" {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("deskripsi must not be empty"))
	}
	req.Deskripsi = deskripsi

	if len(form.File["photos"]) == 0 {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("photos must not be empty"))
	}
	photos := form.File["photos"]
	photoFilePaths := []string{}
	photoUrls := []string{}
	for _, photo := range photos {
		fileExtension := photo.Filename[strings.LastIndex(photo.Filename, ".")+1:]
		if fileExtension != "jpg" &&
			fileExtension != "jpeg" &&
			fileExtension != "png" &&
			fileExtension != "webp" &&
			fileExtension != "jfif" {
			return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("photo must be jpg/jpeg/png/webp/jfif"))
		}
		photoFilename := uuid.NewString() + "." + fileExtension
		rootFolderPath, err := filepath.Abs("./")
		if err != nil {
			return helper.ResponseErrorJson(c, fiber.StatusInternalServerError, err)
		}
		photoFilePath := filepath.Join(rootFolderPath, "uploads", "produk", photoFilename)
		photoUrl := c.BaseURL() + "/uploads/produk/" + photoFilename

		err = c.SaveFile(photo, photoFilePath)
		if err != nil {
			// when a file is failed to be saved on disk, delete already created photo on the disk
			for _, createdPhotoFilePath := range photoFilePaths {
				errRemove := os.Remove(createdPhotoFilePath)
				if errRemove != nil {
					return helper.ResponseErrorJson(c, http.StatusInternalServerError, errRemove)
				}
			}
			return helper.ResponseErrorJson(c, http.StatusInternalServerError, err)
		}
		photoFilePaths = append(photoFilePaths, photoFilePath)
		photoUrls = append(photoUrls, photoUrl)
	}
	req.PhotoUrls = photoUrls

	// get before updated photo urls to delete on the disk if update is successful
	oldPhotoFilePaths := []string{}
	oldProdukResponse, err := p.produkUsecase.GetProdukByID(ctx, produkIdInt)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	for _, fotoProdukResponse := range oldProdukResponse.Photos {
		photoUrl := fotoProdukResponse.Url
		photoFilename := path.Base(photoUrl)
		rootFolderPath, err := filepath.Abs("./")
		if err != nil {
			return helper.ResponseErrorJson(c, fiber.StatusInternalServerError, err)
		}
		photoFilePath := filepath.Join(rootFolderPath, "uploads", "produk", photoFilename)
		oldPhotoFilePaths = append(oldPhotoFilePaths, photoFilePath)
	}

	produkResponse, err := p.produkUsecase.EditProdukByID(ctx, produkIdInt, userId, &req)
	if err != nil {
		// delete already created photos on the disk if the store produk operation is failed
		for _, createdPhotoFilePath := range photoFilePaths {
			errRemove := os.Remove(createdPhotoFilePath)
			if errRemove != nil {
				return helper.ResponseErrorJson(c, http.StatusInternalServerError, errRemove)
			}
		}
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}

	// delete old photo files
	for _, oldPhotoFilePath := range oldPhotoFilePaths {
		errRemove := os.Remove(oldPhotoFilePath)
		if errRemove != nil {
			return helper.ResponseErrorJson(c, http.StatusInternalServerError, errRemove)
		}
	}

	return helper.ResponseSuccessJson(c, produkResponse)
}

func (p *produkDelivery) DeleteProdukHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	idString := c.Params("id")
	idInt, err := strconv.Atoi(idString)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, errors.New("invalid id"))
	}

	userId, err := helper.GetUserIdFromToken(c)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}

	// get old photo files to be deleted later when delete is successful
	oldPhotoFilePaths := []string{}
	oldProdukResponse, err := p.produkUsecase.GetProdukByID(ctx, idInt)
	if err != nil {
		return helper.ResponseErrorJson(c, fiber.StatusBadRequest, err)
	}
	for _, fotoProdukResponse := range oldProdukResponse.Photos {
		photoUrl := fotoProdukResponse.Url
		photoFilename := path.Base(photoUrl)
		rootFolderPath, err := filepath.Abs("./")
		if err != nil {
			return helper.ResponseErrorJson(c, fiber.StatusInternalServerError, err)
		}
		photoFilePath := filepath.Join(rootFolderPath, "uploads", "produk", photoFilename)
		oldPhotoFilePaths = append(oldPhotoFilePaths, photoFilePath)
	}

	err = p.produkUsecase.DestroyProduk(ctx, idInt, userId)
	if err != nil {
		return helper.ResponseErrorJson(c, http.StatusBadRequest, err)
	}

	// delete old photo files
	for _, oldPhotoFilePath := range oldPhotoFilePaths {
		errRemove := os.Remove(oldPhotoFilePath)
		if errRemove != nil {
			return helper.ResponseErrorJson(c, http.StatusInternalServerError, errRemove)
		}
	}

	return helper.ResponseSuccessJson(c, "")
}
