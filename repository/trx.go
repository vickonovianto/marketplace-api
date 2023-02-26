package repository

import (
	"context"
	"errors"
	"marketplace-api/config"
	"marketplace-api/model"
	"strings"

	"gorm.io/gorm"
)

type trxRepository struct {
	Cfg config.Config
}

func NewTrxRepository(cfg config.Config) model.TrxRepository {
	return &trxRepository{Cfg: cfg}
}

func (t *trxRepository) CreateTrx(
	ctx context.Context,
	trx *model.Trx,
	detailTrxWithLogProdukList []*model.DetailTrxWithLogProduk,
) (*model.Trx, error) {

	transaction := t.Cfg.Database().WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			transaction.Rollback()
		}
	}()

	if err := transaction.Error; err != nil {
		return nil, err
	}

	if err := transaction.Create(&trx).Error; err != nil {
		transaction.Rollback()
		return nil, err
	}

	for _, detailTrxWithLogProduk := range detailTrxWithLogProdukList {
		logProduk := detailTrxWithLogProduk.LogProduk
		detailTrx := detailTrxWithLogProduk.DetailTrx

		produk := new(model.Produk)
		produkId := logProduk.IdProduk
		if err := transaction.First(produk, produkId).Error; err != nil {
			transaction.Rollback()
			return nil, err
		}

		if detailTrx.Kuantitas > produk.Stok {
			return nil, errors.New("kuantitas melebihi stok produk")
		}
		newProdukStok := produk.Stok - detailTrx.Kuantitas
		if err := transaction.
			Model(&model.Produk{ID: produkId}).Update("stok", newProdukStok).Error; err != nil {
			transaction.Rollback()
			return nil, err
		}

		if err := transaction.Create(&logProduk).Error; err != nil {
			transaction.Rollback()
			return nil, err
		}

		detailTrx.IdTrx = trx.ID
		detailTrx.IdLogProduk = logProduk.ID
		if err := transaction.Create(&detailTrx).Error; err != nil {
			transaction.Rollback()
			return nil, err
		}
	}

	return trx, transaction.Commit().Error
}

// can search invoice code, product name, or toko name
func (t *trxRepository) Fetch(ctx context.Context, req *model.TrxFetchRequest, userId int) ([]*model.Trx, error) {
	var data []*model.Trx

	offset := (req.Page - 1) * req.Limit

	if strings.HasPrefix(req.Search, model.KODE_INVOICE_PREFIX) {
		// search invoice code
		if err := t.Cfg.Database().WithContext(ctx).
			Where("id_user = ? AND kode_invoice LIKE ?", userId, req.Search).
			Limit(req.Limit).Offset(offset).
			Find(&data).Error; err != nil {
			return nil, err
		}
	} else {
		// search product name or toko name
		type TrxIdStruct struct {
			ID int
		}
		var trxIdStructList []*TrxIdStruct

		err := t.Cfg.Database().WithContext(ctx).
			Model(&model.Trx{}).
			Where("trx.id_user = ?", userId).
			Joins("JOIN detail_trx ON detail_trx.id_trx = trx.id").
			Joins("JOIN log_produk ON log_produk.id = detail_trx.id_log_produk").
			Joins("JOIN toko ON toko.id = log_produk.id_toko").
			Where("nama_produk LIKE ? OR nama_toko LIKE ?", "%"+req.Search+"%", "%"+req.Search+"%").
			Select("trx.id").
			Find(&trxIdStructList).Error
		if err != nil {
			return nil, err
		}

		trxIdList := []int{}
		for _, trxIdStruct := range trxIdStructList {
			trxIdList = append(trxIdList, trxIdStruct.ID)
		}

		if err := t.Cfg.Database().WithContext(ctx).
			Where("id in ?", trxIdList).
			Limit(req.Limit).Offset(offset).
			Find(&data).Error; err != nil {
			return nil, err
		}
	}

	return data, nil
}

func (t *trxRepository) FindByID(ctx context.Context, trxId int) (*model.Trx, error) {
	trx := new(model.Trx)

	if err := t.Cfg.Database().
		WithContext(ctx).
		First(trx, trxId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return trx, nil
}
