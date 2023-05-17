package domain

import (
	"shortner/data"
	"time"

	"github.com/teris-io/shortid"
)

type shortner struct {
	Short    string `json:"short"`
	Url      string `json:"url"`
	CreateAt string `json:"created_at"`
}

type service struct {
	Storage data.DataStore
}
type Service interface {
	SearchUrl(string) ([]shortner, error)
	Store(string) (shortner, error)
	Find(string) (shortner, error)
	Delete(string) (bool, error)
}

func NewService(d data.DataStore) Service {
	return &service{
		Storage: d,
	}
}

func (S *service) SearchUrl(url string) ([]shortner, error) {
	var snlist []shortner
	data, err := S.Storage.SearchUrl(url)
	if err != nil {
		return nil, err
	}
	for _, i := range data {
		var sn shortner
		sn.Short = i["short"]
		sn.Url = i["url"]
		sn.CreateAt = i["create_at"]
		snlist = append(snlist, sn)

	}
	return snlist, nil

}
func (S *service) Store(url string) (shortner, error) {
	short := shortid.MustGenerate()
	now := time.Now().Format("2006-01-02 15:04:05")
	data := map[string]string{
		"short":     short,
		"url":       url,
		"create_at": now,
	}
	err := S.Storage.Store(data)
	if err != nil {
		return shortner{}, err
	}
	return shortner{
		Short:    short,
		Url:      url,
		CreateAt: now,
	}, nil

}
func (S *service) Find(key string) (shortner, error) {
	data, err := S.Storage.Find(key)
	if err != nil {
		return shortner{}, err
	}
	return shortner{
		Short:    data["short"],
		Url:      data["url"],
		CreateAt: data["create_at"],
	}, nil

}
func (S *service) Delete(key string) (bool, error) {
	return S.Storage.Delete(key)

}
