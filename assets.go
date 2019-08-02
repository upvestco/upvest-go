package upvest

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

// Asset is the resource representing your Upvest Tenant asset.
// For more details see https://doc.upvest.co/reference#assets
type Asset struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Symbol   string                 `json:"symbol"`
	Exponent int64                  `json:"exponent"`
	Protocol string                 `json:"protocol"`
	MetaData map[string]interface{} `json:"metadata"`
}

// AssetService handles operations related to the asset
// For more details see https://doc.upvest.co/reference#assets/
type AssetService struct {
	service
}

// AssetList is a list object for assets.
type AssetList struct {
	Meta   ListMeta
	Values []Asset `json:"results"`
}

// Get returns the details of a asset.
// For more details see
func (s *AssetService) Get(assetID string) (*Asset, error) {
	u := fmt.Sprintf("/assets/%s", assetID)
	asset := &Asset{}
	p := &Params{}
	p.SetAuthProvider(s.auth)
	err := s.client.Call(http.MethodGet, u, nil, asset, p)
	return asset, err
}

// List returns list of all assets.
// For more details see https://doc.upvest.co/reference#asset
func (s *AssetService) List() (*AssetList, error) {
	path := "/assets/"
	u, _ := url.Parse(path)

	p := &Params{}
	p.SetAuthProvider(s.auth)

	var results []Asset
	assets := &AssetList{}

	for {
		err := s.client.Call(http.MethodGet, u.String(), nil, assets, p)
		if err != nil {
			return nil, errors.Wrap(err, "Could not retrieve list of assets")
		}
		results = append(results, assets.Values...)

		// append page_size param to the returned params
		u1, err := url.Parse(assets.Meta.Next)
		q := u1.Query()
		q.Set("page_size", string(MaxPageSize))
		u.RawQuery = q.Encode()
		if assets.Meta.Next == "" {
			break
		}
	}

	return &AssetList{Values: results}, nil
}
