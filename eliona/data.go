package eliona

import (
	"context"
	"fmt"
	"microsoft-365/apiserver"
	"microsoft-365/conf"
	"microsoft-365/msgraph"

	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

func UpsertRoomData(config apiserver.Configuration, rooms []msgraph.Room) error {
	for _, projectId := range *config.ProjectIDs {
		for _, room := range rooms {
			log.Debug("Eliona", "upserting data for room: config %d and room '%v'", config.Id, room.EmailAddress)
			assetId, err := conf.GetAssetId(context.Background(), config, projectId, "microsoft_365_room"+fmt.Sprint(room.EmailAddress))
			if err != nil {
				return err
			}
			if assetId == nil {
				return fmt.Errorf("unable to find asset ID")
			}

			data := asset.Data{
				AssetId: *assetId,
				Data:    room,
			}
			if asset.UpsertAssetDataIfAssetExists(data); err != nil {
				return fmt.Errorf("upserting data: %v", err)
			}
		}
	}
	return nil
}
