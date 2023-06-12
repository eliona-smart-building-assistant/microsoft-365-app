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
	"net/http"
	"template/apiserver"
	"template/apiservices"
	"template/conf"
	"template/msgraph"
	"time"

	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

type ConfigurationWithHelper struct {
	Config apiserver.Configuration
	Graph  *msgraph.GraphHelper
}

var configCache map[int64]ConfigurationWithHelper = make(map[int64]ConfigurationWithHelper)

// collectData is the main app function which is called periodically
func collectData() {
	configs, err := conf.GetConfigs(context.Background())
	if err != nil {
		log.Fatal("conf", "Couldn't read configs from DB: %v", err)
		return
	}
	if len(configs) == 0 {
		log.Info("conf", "No configs in DB")
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

		cachedConfigWithHelper, found := configCache[*config.Id]
		if !found || !sameLogin(cachedConfigWithHelper.Config, config) {
			// todo: refresh also if same login, but different config
			graph := msgraph.NewGraphHelper()
			if err := graph.InitializeGraphForUserAuth(config.ClientId, config.TenantId, []string{}); err != nil {
				log.Error("ms-graph", "initializing graph for user auth: %v", err)
				continue
			}

			// Update the cache.
			configCache[*config.Id] = ConfigurationWithHelper{
				Config: config,
				Graph:  graph,
			}
		}

		common.RunOnceWithParam(func(configH ConfigurationWithHelper) {
			log.Info("main", "Collecting %d started", *configH.Config.Id)

			if err := collectRooms(configH); err != nil {
				return // Error is handled in the method itself.
			}

			log.Info("main", "Collecting %d finished", *configH.Config.Id)

			time.Sleep(time.Second * time.Duration(configH.Config.RefreshInterval))
		}, configCache[*config.Id], *config.Id)
	}
}

func sameLogin(a, b apiserver.Configuration) bool {
	return a.ClientSecret == b.ClientSecret &&
		a.TenantId == b.TenantId &&
		equalStringPtr(a.Username, b.Username) &&
		equalStringPtr(a.Password, b.Password)
}

func equalStringPtr(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func collectRooms(configH ConfigurationWithHelper) error {
	err := configH.Graph.GetRooms()
	if err != nil {
		log.Error("ms-graph", "getting rooms: %v", err)
		return err
	}
	fmt.Printf("")
	// if err := eliona.CreateRoomsAssetsIfNecessary(config, rooms); err != nil {
	// 	log.Error("eliona", "creating location assets: %v", err)
	// 	return err
	// }

	// if err := eliona.UpsertRoomData(config, rooms); err != nil {
	// 	log.Error("eliona", "inserting location data into Eliona: %v", err)
	// 	return err
	// }
	return nil
}

// listenApi starts the API server and listen for requests
func listenApi() {
	err := http.ListenAndServe(":"+common.Getenv("API_SERVER_PORT", "3000"), apiserver.NewRouter(
		apiserver.NewConfigurationApiController(apiservices.NewConfigurationApiService()),
		apiserver.NewVersionApiController(apiservices.NewVersionApiService()),
		apiserver.NewCustomizationApiController(apiservices.NewCustomizationApiService()),
	))
	log.Fatal("main", "API server: %v", err)
}
