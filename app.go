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

package main

import (
	"context"
	"fmt"
	"microsoft-365/apiserver"
	"microsoft-365/apiservices"
	"microsoft-365/conf"
	"microsoft-365/eliona"
	"microsoft-365/msgraph"
	"net/http"
	"sync"
	"time"

	"github.com/eliona-smart-building-assistant/go-eliona/app"
	"github.com/eliona-smart-building-assistant/go-eliona/asset"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/db"
	utilshttp "github.com/eliona-smart-building-assistant/go-utils/http"
	"github.com/eliona-smart-building-assistant/go-utils/log"
	"github.com/gorilla/mux"
)

var once sync.Once

func initialize() {
	ctx := context.Background()

	// Necessary to close used init resources
	conn := db.NewInitConnectionWithContextAndApplicationName(ctx, app.AppName())
	defer conn.Close(ctx)

	// Init the app before the first run.
	app.Init(conn, app.AppName(),
		app.ExecSqlFile("conf/init.sql"),
		asset.InitAssetTypeFiles("eliona/asset-type-*.json"),
	)
}

// collectData is the main app function which is called periodically
func collectData() {
	configs, err := conf.GetConfigsForEliona(context.Background())
	if err != nil {
		log.Fatal("conf", "Couldn't read configs from DB: %v", err)
		return
	}
	if len(configs) == 0 {
		once.Do(func() {
			log.Info("conf", "No configs in DB. Please configure the app in Eliona.")
		})
		return
	}

	for _, config := range configs {
		if !conf.IsConfigEnabled(config) {
			if conf.IsConfigActive(config) {
				conf.SetConfigActiveState(context.Background(), config, false)
			}
			continue
		}

		if !conf.IsConfigActive(config) {
			conf.SetConfigActiveState(context.Background(), config, true)
			log.Info("conf", "Collecting initialized with Configuration %d:\n"+
				"Enable: %t\n"+
				"Refresh Interval: %d\n"+
				"Request Timeout: %d\n"+
				"Project IDs: %v\n",
				*config.Id,
				*config.Enable,
				config.RefreshInterval,
				*config.RequestTimeout,
				*config.ProjectIDs)
		}

		common.RunOnceWithParam(func(config apiserver.Configuration) {
			log.Info("main", "Collecting %d started", *config.Id)

			if err := collectResources(config); err != nil {
				return // Error is handled in the method itself.
			}

			log.Info("main", "Collecting %d finished", *config.Id)

			time.Sleep(time.Second * time.Duration(config.RefreshInterval))
		}, config, *config.Id)
	}
}

func collectResources(config apiserver.Configuration) error {
	graph := msgraph.NewGraphHelper()
	if config.ClientSecret == nil || config.Username == nil || config.Password == nil {
		log.Error("conf", "Shouldn't happen: some values are nil")
		return fmt.Errorf("shouldn't happen: some values are nil")
	}
	if err := graph.InitializeGraph(config.ClientId, config.TenantId, *config.ClientSecret, *config.Username, *config.Password); err != nil {
		log.Error("microsoft-365", "initializing graph for user auth: %v", err)
		return err
	}

	rooms, err := graph.GetRooms(config)
	if err != nil {
		log.Error("microsoft-365", "getting rooms: %v", err)
		return err
	}
	fmt.Printf("got %v rooms.\n", len(rooms))
	if err := eliona.CreateRoomsAssetsIfNecessary(config, rooms); err != nil {
		log.Error("eliona", "creating room assets: %v", err)
		return err
	}

	assets := make([]eliona.Asset, len(rooms))
	for i, v := range rooms {
		assets[i] = v
	}

	equipment, err := graph.GetEquipment(config)
	if err != nil {
		log.Error("microsoft-365", "getting equipment: %v", err)
		return err
	}
	fmt.Printf("got %v equipment.\n", len(equipment))
	if err := eliona.CreateEquipmentAssetsIfNecessary(config, equipment); err != nil {
		log.Error("eliona", "creating equipment assets: %v", err)
		return err
	}

	for _, v := range equipment {
		assets = append(assets, v)
	}

	if err := eliona.UpsertAssetData(config, assets); err != nil {
		log.Error("eliona", "inserting room data into Eliona: %v", err)
		return err
	}
	return nil
}

// listenApi starts the API server and listen for requests
func listenApi() {
	r := mux.NewRouter()
	msproxyUrl := "/v1/msproxy/"
	r.PathPrefix(msproxyUrl).Handler(http.StripPrefix(msproxyUrl, &msgraph.Proxy{}))

	r.PathPrefix("/").Handler(utilshttp.NewCORSEnabledHandler(apiserver.NewRouter(
		apiserver.NewConfigurationAPIController(apiservices.NewConfigurationApiService()),
		apiserver.NewVersionAPIController(apiservices.NewVersionApiService()),
		apiserver.NewCustomizationAPIController(apiservices.NewCustomizationApiService()),
		apiserver.NewBookingAPIController(apiservices.NewBookingAPIService()),
	)))

	err := http.ListenAndServe(":"+common.Getenv("API_SERVER_PORT", "3000"), r)
	log.Fatal("main", "API server: %v", err)
}
