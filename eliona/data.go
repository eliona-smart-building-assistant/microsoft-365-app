package eliona

import (
	"context"
	"fmt"
	"github.com/eliona-smart-building-assistant/go-eliona/v2/asset"
	"microsoft-365/v2/conf"

	"github.com/eliona-smart-building-assistant/go-utils/log"
)

type Asset interface {
	AssetType() string
	Id() string
}

func UpsertAssetData(config conf.Configuration, assets []Asset) error {
	for _, a := range assets {
		log.Debug("Eliona", "upserting data for asset: config %d and asset '%v'", config.Id, a.Id())
		assetId, err := conf.GetAssetId(context.Background(), config, a.Id())
		if err != nil {
			return err
		}
		if assetId == nil {
			return fmt.Errorf("unable to find asset ID")
		}

		data := asset.Data{
			AssetId: *assetId,
			Data:    a,
		}
		if err := asset.UpsertAssetDataIfAssetExists(config.ApiEndpoint(), config.ApiKey(), data); err != nil {
			return fmt.Errorf("upserting data: %v", err)
		}
	}
	return nil
}
