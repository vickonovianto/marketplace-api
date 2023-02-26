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

type cityRepository struct {
	Cfg config.Config
}

func NewCityRepository(cfg config.Config) model.CityRepository {
	return &cityRepository{Cfg: cfg}
}

func (c *cityRepository) FetchAll(ctx context.Context, provinceId string) ([]*model.City, error) {
	client := &http.Client{}
	url := "https://www.emsifa.com/api-wilayah-indonesia/api/regencies/"
	url += provinceId + ".json"
	req, err := http.NewRequest("GET", url, nil)
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
	cities := []*model.City{}
	err = json.Unmarshal(bodyBytes, &cities)
	if err != nil {
		return nil, errors.New("invalid prov_id")
	}
	return cities, nil
}

func (c *cityRepository) FindByID(ctx context.Context, provinceId string, cityId string) (*model.City, error) {
	cities, err := c.FetchAll(ctx, provinceId)
	if err != nil {
		return nil, errors.New("city not found")
	}
	for _, cityPointer := range cities {
		city := *cityPointer
		if city.ID == cityId {
			return cityPointer, nil
		}
	}
	return nil, errors.New("city not found")
}
