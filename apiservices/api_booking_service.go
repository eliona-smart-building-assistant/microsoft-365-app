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
	"microsoft-365/apiserver"
	"net/http"
	"time"
)

// BookingAPIService is a service that implements the logic for the BookingAPIServicer
// This service should implement the business logic for every endpoint for the BookingAPI API.
// Include any external packages or services that will be required by this service.
type BookingAPIService struct {
}

// NewBookingAPIService creates a default api service
func NewBookingAPIService() apiserver.BookingAPIServicer {
	return &BookingAPIService{}
}

// BookingsGet - List bookings
func (s *BookingAPIService) BookingsGet(ctx context.Context, start time.Time, end time.Time, assetId string) (apiserver.ImplResponse, error) {
	// TODO - update BookingsGet with the required logic for this service method.
	// Add api_booking_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, []Booking{}) or use other options such as http.Ok ...
	// return Response(200, []Booking{}), nil

	return apiserver.Response(http.StatusNotImplemented, nil), errors.New("BookingsGet method not implemented")
}

// BookingsPost - Create a booking
func (s *BookingAPIService) BookingsPost(ctx context.Context, createBookingRequest apiserver.CreateBookingRequest) (apiserver.ImplResponse, error) {
	// TODO - update BookingsPost with the required logic for this service method.
	// Add api_booking_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(201, Booking{}) or use other options such as http.Ok ...
	// return Response(201, Booking{}), nil

	// TODO: Uncomment the next line to return response Response(400, {}) or use other options such as http.Ok ...
	// return Response(400, nil),nil

	return apiserver.Response(http.StatusNotImplemented, nil), errors.New("BookingsPost method not implemented")
}

// BookingsBookingIdDelete - Cancel a booking
func (s *BookingAPIService) BookingsBookingIdDelete(ctx context.Context, bookingId string) (apiserver.ImplResponse, error) {
	// TODO - update BookingsBookingIdDelete with the required logic for this service method.
	// Add api_booking_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(204, {}) or use other options such as http.Ok ...
	// return Response(204, nil),nil

	// TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	// return Response(404, nil),nil

	return apiserver.Response(http.StatusNotImplemented, nil), errors.New("BookingsBookingIdDelete method not implemented")
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
