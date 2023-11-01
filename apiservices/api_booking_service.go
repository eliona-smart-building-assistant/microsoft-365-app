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

package apiservices

import (
	"context"
	"errors"
	"fmt"
	"microsoft-365/apiserver"
	"microsoft-365/appdb"
	"microsoft-365/conf"
	"microsoft-365/msgraph"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/eliona-smart-building-assistant/go-utils/log"
)

// BookingAPIService is a service that implements the logic for the BookingAPIServicer
// This service should implement the business logic for every endpoint for the BookingAPI API.
// Include any external packages or services that will be required by this service.
type BookingAPIService struct {
	sessions map[string]authorizedSession
	mu       sync.Mutex
}

type authorizedSession struct {
	asset   *appdb.Asset
	graph   *msgraph.GraphHelper
	created time.Time
}

func (s *BookingAPIService) addSession(deviceCode string, session authorizedSession) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.sessions == nil {
		s.sessions = make(map[string]authorizedSession)
	}
	s.sessions[deviceCode] = session
}

func (s *BookingAPIService) cleanupSessions() {
	const sessionDuration = 7 * 24 * time.Hour

	s.mu.Lock()
	defer s.mu.Unlock()

	for sessionID, session := range s.sessions {
		if time.Since(session.created) > sessionDuration {
			delete(s.sessions, sessionID)
		}
	}
}

// NewBookingAPIService creates a default api service
func NewBookingAPIService() apiserver.BookingAPIServicer {
	return &BookingAPIService{}
}

func fetchDBData(ctx context.Context, assetId string) (*appdb.Asset, *apiserver.Configuration, apiserver.ImplResponse, error) {
	assetIdInt64, err := strconv.ParseInt(assetId, 10, 32)
	if err != nil {
		return nil, nil, apiserver.Response(http.StatusBadRequest, nil), fmt.Errorf("parsing asset ID: %v", err)
	}
	assetIdInt32 := int32(assetIdInt64)
	asset, err := conf.GetAsset(ctx, assetIdInt32)
	if err != nil {
		return nil, nil, apiserver.Response(http.StatusInternalServerError, nil), fmt.Errorf("finding asset: %v", err)
	}
	if asset == nil {
		return nil, nil, apiserver.Response(http.StatusNotFound, nil), errors.New("asset not found")
	}
	config, err := conf.GetConfig(ctx, asset.ConfigurationID)
	if err != nil {
		return nil, nil, apiserver.Response(http.StatusInternalServerError, nil), fmt.Errorf("finding configuration %v: %v", asset.ConfigurationID, err)
	}
	if config == nil {
		return nil, nil, apiserver.Response(http.StatusNotFound, nil), errors.New("configuration not found")
	}
	return asset, config, apiserver.ImplResponse{}, nil
}

// BookingsAuthorizeGet - Authorize user for managing bookings
func (s *BookingAPIService) BookingsAuthorizeGet(ctx context.Context, assetId string) (apiserver.ImplResponse, error) {
	asset, config, resp, err := fetchDBData(ctx, assetId)
	if err != nil {
		return resp, err
	}

	graph := msgraph.NewGraphHelper()
	userCodeChannel := make(chan string)
	if err := graph.InitializeGraphForUserAuth(config.ClientId, config.TenantId, userCodeChannel); err != nil {
		log.Error("microsoft-365", "initializing graph for user auth: %v", err)
		return apiserver.Response(http.StatusInternalServerError, nil), fmt.Errorf("internal server error")
	}
	// Ensure this doesn't hang indefinitely waiting for user action
	ctxTimeout, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()

	errChan := make(chan error)
	go func() {
		err := graph.InitiateAuthorization(context.Background())
		if err != nil {
			log.Debug("msgraph", "testing user request: %v", err)
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return apiserver.Response(http.StatusBadRequest, nil), fmt.Errorf("testing authorization failed: server responded with error: %v", err)
		}
	case userCode := <-userCodeChannel:
		s.addSession(userCode, authorizedSession{
			asset:   asset,
			graph:   graph,
			created: time.Now(),
		})
		go s.cleanupSessions()
		return apiserver.Response(http.StatusOK, userCode), nil
	case <-ctxTimeout.Done():
		return apiserver.Response(http.StatusRequestTimeout, nil), errors.New("the code was not authorized before context done")
	}
	return apiserver.Response(http.StatusBadRequest, nil), fmt.Errorf("undefined: authorization finished without device code")
}

// BookingsGet - List bookings
func (s *BookingAPIService) BookingsGet(ctx context.Context, start, end string, assetId string) (apiserver.ImplResponse, error) {
	asset, config, resp, err := fetchDBData(ctx, assetId)
	if err != nil {
		return resp, err
	}
	graph := msgraph.NewGraphHelper()
	if config.ClientSecret == nil || config.Username == nil || config.Password == nil {
		log.Error("conf", "Shouldn't happen: some values are nil")
		return apiserver.Response(http.StatusInternalServerError, nil), fmt.Errorf("internal server error")
	}
	if err := graph.InitializeGraph(config.ClientId, config.TenantId, *config.ClientSecret, *config.Username, *config.Password); err != nil {
		log.Error("microsoft-365", "initializing graph for user auth: %v", err)
		return apiserver.Response(http.StatusInternalServerError, nil), fmt.Errorf("internal server error")
	}

	bookings, err := graph.ListBookings(ctx, asset.Email, start, end)
	if err != nil {
		log.Error("microsoft-365", "getting events from MS Graph: %v", err)
		return apiserver.Response(http.StatusInternalServerError, nil), fmt.Errorf("internal server error")
	}
	return apiserver.Response(http.StatusOK, bookings), nil
}

// BookingsPost - Create a booking
func (s *BookingAPIService) BookingsPost(ctx context.Context, createBookingRequest apiserver.CreateBookingRequest) (apiserver.ImplResponse, error) {
	session, ok := s.sessions[createBookingRequest.DeviceCode]
	if !ok {
		return apiserver.Response(http.StatusBadRequest, nil), errors.New("invalid device code")
	}

	if err := session.graph.CreateBooking(ctx, createBookingRequest.Start, createBookingRequest.End, session.asset.Email, createBookingRequest.EventName, createBookingRequest.EventName); err != nil {
		log.Error("microsoft-365", "creating event: %v", err)
		return apiserver.Response(http.StatusBadRequest, nil), fmt.Errorf("server responded with error: %v", err)
	}

	return apiserver.Response(http.StatusOK, nil), nil
}

// BookingsBookingIdDelete - Cancel a booking
func (s *BookingAPIService) BookingsDeletePost(ctx context.Context, deleteBookingRequest apiserver.DeleteBookingRequest) (apiserver.ImplResponse, error) {
	session, ok := s.sessions[deleteBookingRequest.DeviceCode]
	if !ok {
		return apiserver.Response(http.StatusBadRequest, nil), errors.New("invalid device code")
	}

	if err := session.graph.DeleteBooking(ctx, deleteBookingRequest.BookingId); err != nil {
		log.Error("microsoft-365", "deleting event %v: %v", deleteBookingRequest.BookingId, err)
		return apiserver.Response(http.StatusBadRequest, nil), fmt.Errorf("server responded with error: %v", err)
	}

	return apiserver.Response(http.StatusOK, nil), nil
}

// BookingsBookingIdRegisterGuestPost - Notify event organizer that a guest came for the event.
func (s *BookingAPIService) BookingsBookingIdRegisterGuestPost(ctx context.Context, bookingId string, bookingsBookingIdRegisterGuestPostRequest apiserver.BookingsBookingIdRegisterGuestPostRequest) (apiserver.ImplResponse, error) {
	// TODO - update BookingsBookingIdRegisterGuestPost with the required logic for this service method.
	// Add api_booking_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(204, {}) or use other options such as http.Ok ...
	// return Response(204, nil),nil

	// TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	// return Response(404, nil),nil

	// TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	// return Response(400, nil),nil

	return apiserver.Response(http.StatusNotImplemented, nil), errors.New("BookingsBookingIdRegisterGuestPost method not implemented")
}
