package usecase

import (
	"context"
	"errors"
	"marketplace-api/model"
	"math/rand"
	"strconv"

	"github.com/jinzhu/copier"
)

type trxUsecase struct {
	trxRepository        model.TrxRepository
	alamatRepository     model.AlamatRepository
	detailTrxRepository  model.DetailTrxRepository
	logProdukRepository  model.LogProdukRepository
	tokoRepository       model.TokoRepository
	categoryRepository   model.CategoryRepository
	fotoProdukRepository model.FotoProdukRepository
	produkRepository     model.ProdukRepository
}

func NewTrxUsecase(
	trxRepository model.TrxRepository,
	alamatRepository model.AlamatRepository,
	detailTrxRepository model.DetailTrxRepository,
	logProdukRepository model.LogProdukRepository,
	tokoRepository model.TokoRepository,
	categoryRepository model.CategoryRepository,
	fotoProdukRepository model.FotoProdukRepository,
	produkRepository model.ProdukRepository,
) model.TrxUsecase {
	return &trxUsecase{
		trxRepository:        trxRepository,
		alamatRepository:     alamatRepository,
		detailTrxRepository:  detailTrxRepository,
		logProdukRepository:  logProdukRepository,
		tokoRepository:       tokoRepository,
		categoryRepository:   categoryRepository,
		fotoProdukRepository: fotoProdukRepository,
		produkRepository:     produkRepository,
	}
}

func (t *trxUsecase) StoreTrx(ctx context.Context, req *model.TrxStoreRequest, userId int) (*model.TrxGetByIDResponse, error) {
	detailTrxWithLogProdukList := []*model.DetailTrxWithLogProduk{}
	trxHargaTotal := 0
	for _, detailTrxRequest := range req.DetailTrxRequests {
		detailTrxWithLogProduk := new(model.DetailTrxWithLogProduk)

		produkId := detailTrxRequest.ProductId
		produk, err := t.produkRepository.FindByID(ctx, produkId)
		if err != nil {
			return nil, err
		}
		logProdukRequest := new(model.LogProdukRequest)
		copier.Copy(logProdukRequest, produk)
		logProduk := new(model.LogProduk)
		copier.Copy(logProduk, logProdukRequest)
		logProduk.IdProduk = produkId

		toko, err := t.tokoRepository.FindByTokoID(ctx, produk.IdToko)
		if err != nil {
			return nil, err
		}
		if toko.IdUser == userId {
			return nil, errors.New("cannot buy product on self-owned store")
		}

		detailTrx := new(model.DetailTrx)
		detailTrx.IdToko = produk.IdToko
		detailTrx.Kuantitas = detailTrxRequest.Kuantitas
		hargaKonsumenInt, err := strconv.Atoi(produk.HargaKonsumen)
		if err != nil {
			return nil, err
		}
		detailTrx.HargaTotal = detailTrx.Kuantitas * hargaKonsumenInt
		trxHargaTotal += detailTrx.HargaTotal

		detailTrxWithLogProduk.LogProduk = logProduk
		detailTrxWithLogProduk.DetailTrx = detailTrx
		detailTrxWithLogProdukList = append(detailTrxWithLogProdukList, detailTrxWithLogProduk)
	}

	trx := new(model.Trx)
	trx.IdUser = userId
	trx.HargaTotal = trxHargaTotal
	trx.MethodBayar = req.MethodBayar
	trx.KodeInvoice = model.KODE_INVOICE_PREFIX + strconv.Itoa(1000000000+rand.Intn(9999999999-1000000000))

	alamat, err := t.alamatRepository.FindByID(ctx, req.AlamatPengiriman)
	if err != nil {
		return nil, err
	}
	if alamat.IdUser != userId {
		return nil, errors.New("unauthorized")
	}
	trx.AlamatPengiriman = req.AlamatPengiriman

	trx, err = t.trxRepository.CreateTrx(ctx, trx, detailTrxWithLogProdukList)
	if err != nil {
		return nil, err
	}
	trxGetByIDResponse := new(model.TrxGetByIDResponse)
	copier.Copy(trxGetByIDResponse, trx)

	alamatResponse := new(model.AlamatResponse)
	copier.Copy(alamatResponse, alamat)
	trxGetByIDResponse.AlamatPengiriman = alamatResponse

	detailTrxResponses := []*model.DetailTrxResponse{}
	detailTrxList, err := t.detailTrxRepository.FindByTrxID(ctx, trx.ID)
	if err != nil {
		return nil, err
	}
	for _, detailTrx := range detailTrxList {
		detailTrxResponse := new(model.DetailTrxResponse)
		copier.Copy(detailTrxResponse, detailTrx)

		logProduk, err := t.logProdukRepository.FindByID(ctx, detailTrx.IdLogProduk)
		if err != nil {
			return nil, err
		}
		logProdukResponse := new(model.LogProdukResponse)
		copier.Copy(logProdukResponse, logProduk)

		toko, err := t.tokoRepository.FindByTokoID(ctx, logProduk.IdToko)
		if err != nil {
			return nil, err
		}
		tokoLogProdukResponse := new(model.TokoLogProdukResponse)
		copier.Copy(tokoLogProdukResponse, toko)
		logProdukResponse.Toko = tokoLogProdukResponse

		category, err := t.categoryRepository.FindByID(ctx, logProduk.IdCategory)
		if err != nil {
			return nil, err
		}
		categoryResponse := new(model.CategoryResponse)
		copier.Copy(categoryResponse, category)
		logProdukResponse.Category = categoryResponse

		fotoProdukResponses := []*model.FotoProdukResponse{}
		fotoProdukList, err := t.fotoProdukRepository.FetchByProdukId(ctx, logProduk.IdProduk)
		if err != nil {
			return nil, err
		}
		copier.Copy(&fotoProdukResponses, &fotoProdukList)
		logProdukResponse.Photos = fotoProdukResponses

		detailTrxResponse.LogProduk = logProdukResponse

		tokoGetByIDResponse := new(model.TokoGetByIDResponse)
		copier.Copy(tokoGetByIDResponse, toko)
		detailTrxResponse.Toko = tokoGetByIDResponse

		detailTrxResponses = append(detailTrxResponses, detailTrxResponse)
	}

	trxGetByIDResponse.DetailTrxResponses = detailTrxResponses

	return trxGetByIDResponse, nil
}

func (t *trxUsecase) FetchTrx(ctx context.Context, req *model.TrxFetchRequest, userId int) (*model.TrxFetchResponse, error) {
	trxList, err := t.trxRepository.Fetch(ctx, req, userId)
	if err != nil {
		return nil, err
	}
	trxFetchResponse := new(model.TrxFetchResponse)
	copier.Copy(trxFetchResponse, req)

	trxGetByIDResponses := []*model.TrxGetByIDResponse{}

	for _, trx := range trxList {
		trxGetByIDResponse := new(model.TrxGetByIDResponse)
		copier.Copy(trxGetByIDResponse, trx)

		alamat, err := t.alamatRepository.FindByID(ctx, trx.AlamatPengiriman)
		if err != nil {
			return nil, err
		}
		alamatResponse := new(model.AlamatResponse)
		copier.Copy(alamatResponse, alamat)
		trxGetByIDResponse.AlamatPengiriman = alamatResponse

		detailTrxResponses := []*model.DetailTrxResponse{}
		detailTrxList, err := t.detailTrxRepository.FindByTrxID(ctx, trx.ID)
		if err != nil {
			return nil, err
		}
		for _, detailTrx := range detailTrxList {
			detailTrxResponse := new(model.DetailTrxResponse)
			copier.Copy(detailTrxResponse, detailTrx)

			logProduk, err := t.logProdukRepository.FindByID(ctx, detailTrx.IdLogProduk)
			if err != nil {
				return nil, err
			}
			logProdukResponse := new(model.LogProdukResponse)
			copier.Copy(logProdukResponse, logProduk)

			toko, err := t.tokoRepository.FindByTokoID(ctx, logProduk.IdToko)
			if err != nil {
				return nil, err
			}
			tokoLogProdukResponse := new(model.TokoLogProdukResponse)
			copier.Copy(tokoLogProdukResponse, toko)
			logProdukResponse.Toko = tokoLogProdukResponse

			category, err := t.categoryRepository.FindByID(ctx, logProduk.IdCategory)
			if err != nil {
				return nil, err
			}
			categoryResponse := new(model.CategoryResponse)
			copier.Copy(categoryResponse, category)
			logProdukResponse.Category = categoryResponse

			fotoProdukResponses := []*model.FotoProdukResponse{}
			fotoProdukList, err := t.fotoProdukRepository.FetchByProdukId(ctx, logProduk.IdProduk)
			if err != nil {
				return nil, err
			}
			copier.Copy(&fotoProdukResponses, &fotoProdukList)
			logProdukResponse.Photos = fotoProdukResponses

			detailTrxResponse.LogProduk = logProdukResponse

			tokoGetByIDResponse := new(model.TokoGetByIDResponse)
			copier.Copy(tokoGetByIDResponse, toko)
			detailTrxResponse.Toko = tokoGetByIDResponse

			detailTrxResponses = append(detailTrxResponses, detailTrxResponse)
		}

		trxGetByIDResponse.DetailTrxResponses = detailTrxResponses

		trxGetByIDResponses = append(trxGetByIDResponses, trxGetByIDResponse)
	}

	trxFetchResponse.Data = trxGetByIDResponses

	return trxFetchResponse, nil
}

func (t *trxUsecase) GetTrxByID(ctx context.Context, trxId int, userId int) (*model.TrxGetByIDResponse, error) {
	trx, err := t.trxRepository.FindByID(ctx, trxId)
	if err != nil {
		return nil, err
	}
	if trx.IdUser != userId {
		return nil, errors.New("unauthorized")
	}

	trxGetByIDResponse := new(model.TrxGetByIDResponse)
	copier.Copy(trxGetByIDResponse, trx)

	alamat, err := t.alamatRepository.FindByID(ctx, trx.AlamatPengiriman)
	if err != nil {
		return nil, err
	}
	alamatResponse := new(model.AlamatResponse)
	copier.Copy(alamatResponse, alamat)
	trxGetByIDResponse.AlamatPengiriman = alamatResponse

	detailTrxResponses := []*model.DetailTrxResponse{}
	detailTrxList, err := t.detailTrxRepository.FindByTrxID(ctx, trx.ID)
	if err != nil {
		return nil, err
	}
	for _, detailTrx := range detailTrxList {
		detailTrxResponse := new(model.DetailTrxResponse)
		copier.Copy(detailTrxResponse, detailTrx)

		logProduk, err := t.logProdukRepository.FindByID(ctx, detailTrx.IdLogProduk)
		if err != nil {
			return nil, err
		}
		logProdukResponse := new(model.LogProdukResponse)
		copier.Copy(logProdukResponse, logProduk)

		toko, err := t.tokoRepository.FindByTokoID(ctx, logProduk.IdToko)
		if err != nil {
			return nil, err
		}
		tokoLogProdukResponse := new(model.TokoLogProdukResponse)
		copier.Copy(tokoLogProdukResponse, toko)
		logProdukResponse.Toko = tokoLogProdukResponse

		category, err := t.categoryRepository.FindByID(ctx, logProduk.IdCategory)
		if err != nil {
			return nil, err
		}
		categoryResponse := new(model.CategoryResponse)
		copier.Copy(categoryResponse, category)
		logProdukResponse.Category = categoryResponse

		fotoProdukResponses := []*model.FotoProdukResponse{}
		fotoProdukList, err := t.fotoProdukRepository.FetchByProdukId(ctx, logProduk.IdProduk)
		if err != nil {
			return nil, err
		}
		copier.Copy(&fotoProdukResponses, &fotoProdukList)
		logProdukResponse.Photos = fotoProdukResponses

		detailTrxResponse.LogProduk = logProdukResponse

		tokoGetByIDResponse := new(model.TokoGetByIDResponse)
		copier.Copy(tokoGetByIDResponse, toko)
		detailTrxResponse.Toko = tokoGetByIDResponse

		detailTrxResponses = append(detailTrxResponses, detailTrxResponse)
	}

	trxGetByIDResponse.DetailTrxResponses = detailTrxResponses

	return trxGetByIDResponse, nil
}
