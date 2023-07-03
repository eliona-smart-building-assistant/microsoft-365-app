//  This file is part of the eliona project.
//  Copyright © 2022 LEICOM iTEC AG. All Rights Reserved.
//  ______ _ _
// |  ____| (_)
// | |__  | |_  ___  _ __   __ _
// |  __| | | |/ _ \| '_ \ / _` |
// | |____| | | (_) | | | | (_| |
// |______|_|_|\___/|_| |_|\__,_|
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING
//  BUT NOT LIMITED  TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NON INFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
//  DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package eliona

import (
	"context"
	"fmt"
	"microsoft-365/apiserver"
	"microsoft-365/conf"
	"microsoft-365/msgraph"

	api "github.com/eliona-smart-building-assistant/go-eliona-api-client/v2"
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

func CreateRoomsAssetsIfNecessary(config apiserver.Configuration, rooms []msgraph.Room) error {
	for _, projectId := range conf.ProjIds(config) {
		_, rootAssetID, err := upsertAsset(assetData{
			config:                  config,
			projectId:               projectId,
			parentLocationalAssetId: nil,
			identifier:              "microsoft_365_root",
			assetType:               "microsoft_365_root",
			name:                    "Microsoft 365",
			description:             "Root asset for Microsoft 365 places",
		})
		if err != nil {
			return fmt.Errorf("upserting root asset: %v", err)
		}
		for _, room := range rooms {
			id := room.EmailAddress
			assetType := "microsoft_365_room"
			name := room.DisplayName
			_, _, err := upsertAsset(assetData{
				config:                  config,
				projectId:               projectId,
				parentLocationalAssetId: &rootAssetID,
				identifier:              fmt.Sprintf("%s_%s", assetType, *id),
				assetType:               assetType,
				name:                    *name,
				description:             fmt.Sprintf("%s (%v)", *name, *id),
			})
			if err != nil {
				return fmt.Errorf("upserting room %s: %v", *id, err)
			}
		}
	}
	return nil
}

type assetData struct {
	config                  apiserver.Configuration
	projectId               string
	parentFunctionalAssetId *int32
	parentLocationalAssetId *int32
	identifier              string
	assetType               string
	name                    string
	description             string
}

func upsertAsset(d assetData) (created bool, assetID int32, err error) {
	// Get known asset id from configuration
	currentAssetID, err := conf.GetAssetId(context.Background(), d.config, d.projectId, d.identifier)
	if err != nil {
		return false, 0, fmt.Errorf("finding asset ID: %v", err)
	}
	if currentAssetID != nil {
		return false, *currentAssetID, nil
	}

	a := api.Asset{
		ProjectId:               d.projectId,
		GlobalAssetIdentifier:   d.identifier,
		Name:                    *api.NewNullableString(common.Ptr(d.name)),
		AssetType:               d.assetType,
		Description:             *api.NewNullableString(common.Ptr(d.description)),
		ParentFunctionalAssetId: *api.NewNullableInt32(d.parentFunctionalAssetId),
		ParentLocationalAssetId: *api.NewNullableInt32(d.parentLocationalAssetId),
		IsTracker:               *api.NewNullableBool(common.Ptr(false)),
	}
	newID, err := asset.UpsertAsset(a)
	if err != nil {
		return false, 0, fmt.Errorf("upserting asset %+v into Eliona: %v", a, err)
	}
	if newID == nil {
		return false, 0, fmt.Errorf("cannot create asset %s", d.name)
	}

	// Remember the asset id for further usage
	if err := conf.InsertAsset(context.Background(), d.config, d.projectId, d.identifier, *newID); err != nil {
		return false, 0, fmt.Errorf("inserting asset to config db: %v", err)
	}

	log.Debug("eliona", "Created new asset for project %s and device %s.", d.projectId, d.identifier)

	return true, *newID, nil
}
