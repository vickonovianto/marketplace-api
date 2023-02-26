package repository

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"marketplace-api/config"
	"marketplace-api/model"
	"net/http"
)

type provinceRepository struct {
	Cfg config.Config
}

func NewProvinceRepository(cfg config.Config) model.ProvinceRepository {
	return &provinceRepository{Cfg: cfg}
}

func (p *provinceRepository) FetchAll(ctx context.Context) ([]*model.Province, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.emsifa.com/api-wilayah-indonesia/api/provinces.json", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	provinces := []*model.Province{}
	err = json.Unmarshal(bodyBytes, &provinces)
	if err != nil {
		return nil, err
	}
	return provinces, nil
}

func (p *provinceRepository) FindByID(ctx context.Context, id string) (*model.Province, error) {
	provinces, err := p.FetchAll(ctx)
	if err != nil {
		return nil, err
	}
	for _, provincePointer := range provinces {
		province := *provincePointer
		if province.ID == id {
			return provincePointer, nil
		}
	}
	return nil, errors.New("province not found")
}
