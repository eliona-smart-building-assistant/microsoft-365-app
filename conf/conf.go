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

package conf

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"microsoft-365/apiserver"
	"microsoft-365/appdb"

	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
)

var ErrBadRequest = errors.New("bad request")

func InsertConfig(ctx context.Context, config apiserver.Configuration) (apiserver.Configuration, error) {
	dbConfig, err := dbConfigFromApiConfig(config)
	if err != nil {
		return apiserver.Configuration{}, fmt.Errorf("creating DB config from API config: %v", err)
	}
	if err := dbConfig.InsertG(ctx, boil.Infer()); err != nil {
		return apiserver.Configuration{}, fmt.Errorf("inserting DB config: %v", err)
	}
	return config, nil
}

func UpsertConfig(ctx context.Context, config apiserver.Configuration) (apiserver.Configuration, error) {
	dbConfig, err := dbConfigFromApiConfig(config)
	if err != nil {
		return apiserver.Configuration{}, fmt.Errorf("creating DB config from API config: %v", err)
	}
	if err := dbConfig.UpsertG(ctx, true, []string{"id"}, boil.Blacklist("id"), boil.Infer()); err != nil {
		return apiserver.Configuration{}, fmt.Errorf("inserting DB config: %v", err)
	}
	return config, nil
}

func GetConfig(ctx context.Context, configID int64) (*apiserver.Configuration, error) {
	dbConfig, err := appdb.Configurations(
		appdb.ConfigurationWhere.ID.EQ(configID),
	).OneG(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching config from database: %v", err)
	}
	if dbConfig == nil {
		return nil, ErrBadRequest
	}
	apiConfig, err := apiConfigFromDbConfig(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("creating API config from DB config: %v", err)
	}
	return &apiConfig, nil
}

func DeleteConfig(ctx context.Context, configID int64) error {
	if _, err := appdb.Assets(
		appdb.AssetWhere.ConfigurationID.EQ(configID),
	).DeleteAllG(ctx); err != nil {
		return fmt.Errorf("deleting assets from database: %v", err)
	}
	count, err := appdb.Configurations(
		appdb.ConfigurationWhere.ID.EQ(configID),
	).DeleteAllG(ctx)
	if err != nil {
		return fmt.Errorf("deleting config from database: %v", err)
	}
	if count > 1 {
		return fmt.Errorf("shouldn't happen: deleted more (%v) configs by ID", count)
	}
	if count == 0 {
		return ErrBadRequest
	}
	return nil
}

func dbConfigFromApiConfig(apiConfig apiserver.Configuration) (dbConfig appdb.Configuration, err error) {
	dbConfig.ID = null.Int64FromPtr(apiConfig.Id).Int64
	dbConfig.ClientID = apiConfig.ClientId
	if apiConfig.ClientSecret != nil {
		dbConfig.ClientSecret = *apiConfig.ClientSecret
	}
	dbConfig.TenantID = apiConfig.TenantId
	if apiConfig.Username != nil {
		dbConfig.Username = *apiConfig.Username
	}
	if apiConfig.Password != nil {
		dbConfig.Password = *apiConfig.Password
	}
	dbConfig.ForEliona = apiConfig.ForEliona
	dbConfig.ForProxy = apiConfig.ForProxy
	dbConfig.Enable = null.BoolFromPtr(apiConfig.Enable)
	dbConfig.RefreshInterval = apiConfig.RefreshInterval
	if apiConfig.RequestTimeout != nil {
		dbConfig.RequestTimeout = *apiConfig.RequestTimeout
	}
	af, err := json.Marshal(apiConfig.AssetFilter)
	if err != nil {
		return appdb.Configuration{}, fmt.Errorf("marshalling assetFilter: %v", err)
	}
	dbConfig.AssetFilter = null.JSONFrom(af)
	dbConfig.Active = null.BoolFromPtr(apiConfig.Active)
	if apiConfig.ProjectIDs != nil {
		dbConfig.ProjectIds = *apiConfig.ProjectIDs
	}

	return dbConfig, nil
}

func apiConfigFromDbConfig(dbConfig *appdb.Configuration) (apiConfig apiserver.Configuration, err error) {
	apiConfig.Id = &dbConfig.ID
	apiConfig.ClientId = dbConfig.ClientID
	apiConfig.ClientSecret = &dbConfig.ClientSecret
	apiConfig.TenantId = dbConfig.TenantID
	apiConfig.Username = &dbConfig.Username
	apiConfig.Password = &dbConfig.Password
	apiConfig.ForEliona = dbConfig.ForEliona
	apiConfig.ForProxy = dbConfig.ForProxy
	apiConfig.Enable = dbConfig.Enable.Ptr()
	apiConfig.RefreshInterval = dbConfig.RefreshInterval
	apiConfig.RequestTimeout = &dbConfig.RequestTimeout
	if dbConfig.AssetFilter.Valid {
		var af [][]apiserver.FilterRule
		if err := json.Unmarshal(dbConfig.AssetFilter.JSON, &af); err != nil {
			return apiserver.Configuration{}, fmt.Errorf("unmarshalling assetFilter: %v", err)
		}
		apiConfig.AssetFilter = af
	}
	apiConfig.Active = dbConfig.Active.Ptr()
	apiConfig.ProjectIDs = common.Ptr[[]string](dbConfig.ProjectIds)
	return apiConfig, nil
}

func GetConfigs(ctx context.Context) ([]apiserver.Configuration, error) {
	dbConfigs, err := appdb.Configurations().AllG(ctx)
	if err != nil {
		return nil, err
	}
	var apiConfigs []apiserver.Configuration
	for _, dbConfig := range dbConfigs {
		ac, err := apiConfigFromDbConfig(dbConfig)
		if err != nil {
			return nil, fmt.Errorf("creating API config from DB config: %v", err)
		}
		apiConfigs = append(apiConfigs, ac)
	}
	return apiConfigs, nil
}

func GetConfigsForEliona(ctx context.Context) ([]apiserver.Configuration, error) {
	dbConfigs, err := appdb.Configurations(
		appdb.ConfigurationWhere.ForEliona.EQ(true),
	).AllG(ctx)
	if err != nil {
		return nil, err
	}
	var apiConfigs []apiserver.Configuration
	for _, dbConfig := range dbConfigs {
		ac, err := apiConfigFromDbConfig(dbConfig)
		if err != nil {
			return nil, fmt.Errorf("creating API config from DB config: %v", err)
		}
		apiConfigs = append(apiConfigs, ac)
	}
	return apiConfigs, nil
}

func GetConfigsForProxy(ctx context.Context) ([]apiserver.Configuration, error) {
	dbConfigs, err := appdb.Configurations(
		appdb.ConfigurationWhere.Enable.EQ(null.BoolFrom(true)),
		appdb.ConfigurationWhere.ForProxy.EQ(true),
	).AllG(ctx)
	if err != nil {
		return nil, err
	}
	var apiConfigs []apiserver.Configuration
	for _, dbConfig := range dbConfigs {
		ac, err := apiConfigFromDbConfig(dbConfig)
		if err != nil {
			return nil, fmt.Errorf("creating API config from DB config: %v", err)
		}
		apiConfigs = append(apiConfigs, ac)
	}
	return apiConfigs, nil
}

func GetConfigsForProxyWithProjectId(ctx context.Context, projectId string) ([]apiserver.Configuration, error) {
	var dbConfigs []*appdb.Configuration
	q := fmt.Sprintf(`
		SELECT *
		FROM microsoft_365.configuration
		WHERE
			enabled = true AND
			for_proxy = true AND
			$1 = ANY (project_ids)
	`)
	if err := queries.RawG(q, projectId).BindG(ctx, dbConfigs); err != nil {
		return nil, fmt.Errorf("fetching configuration: %v", err)
	}

	var apiConfigs []apiserver.Configuration
	for _, dbConfig := range dbConfigs {
		ac, err := apiConfigFromDbConfig(dbConfig)
		if err != nil {
			return nil, fmt.Errorf("creating API config from DB config: %v", err)
		}
		apiConfigs = append(apiConfigs, ac)
	}
	return apiConfigs, nil
}

func SetConfigActiveState(ctx context.Context, config apiserver.Configuration, state bool) (int64, error) {
	return appdb.Configurations(
		appdb.ConfigurationWhere.ID.EQ(null.Int64FromPtr(config.Id).Int64),
	).UpdateAllG(ctx, appdb.M{
		appdb.ConfigurationColumns.Active: state,
	})
}

func ProjIds(config apiserver.Configuration) []string {
	if config.ProjectIDs == nil {
		return []string{}
	}
	return *config.ProjectIDs
}

func IsConfigActive(config apiserver.Configuration) bool {
	return config.Active == nil || *config.Active
}

func IsConfigEnabled(config apiserver.Configuration) bool {
	return config.Enable == nil || *config.Enable
}

func SetAllConfigsInactive(ctx context.Context) (int64, error) {
	return appdb.Configurations().UpdateAllG(ctx, appdb.M{
		appdb.ConfigurationColumns.Active: false,
	})
}

func InsertAsset(ctx context.Context, config apiserver.Configuration, projId string, globalAssetID string, assetId int32, email string) error {
	dbAsset := appdb.Asset{
		ConfigurationID: null.Int64FromPtr(config.Id).Int64,
		ProjectID:       projId,
		GlobalAssetID:   globalAssetID,
		AssetID:         null.Int32From(assetId),
		Email:           email,
	}
	return dbAsset.InsertG(ctx, boil.Infer())
}

func GetAssetId(ctx context.Context, config apiserver.Configuration, projId string, globalAssetID string) (*int32, error) {
	dbAsset, err := appdb.Assets(
		appdb.AssetWhere.ConfigurationID.EQ(null.Int64FromPtr(config.Id).Int64),
		appdb.AssetWhere.ProjectID.EQ(projId),
		appdb.AssetWhere.GlobalAssetID.EQ(globalAssetID),
	).AllG(ctx)
	if err != nil || len(dbAsset) == 0 {
		return nil, err
	}
	return common.Ptr(dbAsset[0].AssetID.Int32), nil
}

func GetAsset(ctx context.Context, assetId int32) (*appdb.Asset, error) {
	return appdb.Assets(
		appdb.AssetWhere.AssetID.EQ(null.Int32From(assetId)),
	).OneG(ctx)
}
