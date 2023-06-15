package msgraph

import (
	"context"
	"fmt"

	azcore "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	auth "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

type GraphHelper struct {
	credential      azcore.TokenCredential
	userClient      *msgraphsdk.GraphServiceClient
	graphUserScopes []string
}

func NewGraphHelper() *GraphHelper {
	g := &GraphHelper{}
	return g
}

func (g *GraphHelper) InitializeGraphForUserAuth(clientId, tenantId, clientSecret, username, password string) error {
	if username != "" {
		cred, err := azidentity.NewUsernamePasswordCredential(
			tenantId,
			clientId,
			username,
			password,
			nil,
		)
		if err != nil {
			return fmt.Errorf("creating the username/password credential: %v", err)
		}
		g.credential = cred
	} else {
		cred, err := azidentity.NewClientSecretCredential(
			tenantId,
			clientId,
			clientSecret,
			nil,
		)
		if err != nil {
			return fmt.Errorf("creating the device code credential: %v", err)
		}
		g.credential = cred
	}

	authProvider, err := auth.NewAzureIdentityAuthenticationProviderWithScopes(g.credential, g.graphUserScopes)
	if err != nil {
		return fmt.Errorf("Creating an auth provider: %v", err)
	}

	adapter, err := msgraphsdk.NewGraphRequestAdapter(authProvider)
	if err != nil {
		return fmt.Errorf("Creating a request adapter: %v", err)
	}

	client := msgraphsdk.NewGraphServiceClient(adapter)
	g.userClient = client

	return nil
}

func (g *GraphHelper) GetUserToken() (*string, error) {
	token, err := g.credential.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: g.graphUserScopes,
	})
	if err != nil {
		return nil, err
	}

	return &token.Token, nil
}

func (g *GraphHelper) GetUser() (models.Userable, error) {
	query := users.UserItemRequestBuilderGetQueryParameters{
		// Only request specific properties
		Select: []string{"displayName", "mail", "userPrincipalName"},
	}

	return g.userClient.Me().Get(context.Background(),
		&users.UserItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &query,
		})
}

func (g *GraphHelper) GetInbox() (models.MessageCollectionResponseable, error) {
	var topValue int32 = 25
	query := users.ItemMailFoldersItemMessagesRequestBuilderGetQueryParameters{
		// Only request specific properties
		Select: []string{"from", "isRead", "receivedDateTime", "subject"},
		// Get at most 25 results
		Top: &topValue,
		// Sort by received time, newest first
		Orderby: []string{"receivedDateTime DESC"},
	}

	return g.userClient.Me().MailFolders().
		ByMailFolderId("inbox").
		Messages().
		Get(context.Background(),
			&users.ItemMailFoldersItemMessagesRequestBuilderGetRequestConfiguration{
				QueryParameters: &query,
			})
}

func (g *GraphHelper) SendMail(subject *string, body *string, recipient *string) error {
	message := models.NewMessage()
	message.SetSubject(subject)

	messageBody := models.NewItemBody()
	messageBody.SetContent(body)
	contentType := models.TEXT_BODYTYPE
	messageBody.SetContentType(&contentType)
	message.SetBody(messageBody)

	toRecipient := models.NewRecipient()
	emailAddress := models.NewEmailAddress()
	emailAddress.SetAddress(recipient)
	toRecipient.SetEmailAddress(emailAddress)
	message.SetToRecipients([]models.Recipientable{
		toRecipient,
	})

	sendMailBody := users.NewItemSendMailPostRequestBody()
	sendMailBody.SetMessage(message)

	return g.userClient.Me().SendMail().Post(context.Background(), sendMailBody, nil)
}

type PhysicalAddress struct {
	City            *string `eliona:"city"`
	CountryOrRegion *string `eliona:"country_or_region"`
	PostalCode      *string `eliona:"postal_code"`
	State           *string `eliona:"state"`
	Street          *string `eliona:"street"`
}

type GeoCoordinates struct {
	Accuracy         *float64 `eliona:"accuracy"`
	Altitude         *float64 `eliona:"altitude"`
	AltitudeAccuracy *float64 `eliona:"altitude_accuracy"`
	Latitude         *float64 `eliona:"latitude"`
	Longitude        *float64 `eliona:"longitude"`
}

type Room struct {
	Address                PhysicalAddress `eliona:"address"`
	DisplayName            *string         `eliona:"display_name"`
	GeoCoordinates         GeoCoordinates  `eliona:"geo_coordinates"`
	Phone                  *string         `eliona:"phone"`
	AudioDeviceName        *string         `eliona:"audio_device_name"`
	BookingType            BookingType     `eliona:"booking_type"`
	Building               *string         `eliona:"building"`
	Capacity               *int32          `eliona:"capacity"`
	DisplayDeviceName      *string         `eliona:"display_device_name"`
	EmailAddress           *string         `eliona:"email_address"`
	FloorLabel             *string         `eliona:"floor_label"`
	FloorNumber            *int32          `eliona:"floor_number"`
	IsWheelChairAccessible *bool           `eliona:"is_wheel_chair_accessible"`
	Label                  *string         `eliona:"label"`
	Nickname               *string         `eliona:"nickname"`
	Tags                   []string        `eliona:"tags"`
	VideoDeviceName        *string         `eliona:"video_device_name"`
}

type BookingType int

const (
	UNKNOWN_BOOKINGTYPE BookingType = iota
	STANDARD_BOOKINGTYPE
	RESERVED_BOOKINGTYPE
)

func (i BookingType) String() string {
	return []string{"unknown", "standard", "reserved"}[i]
}

func mapBookingType(bt *models.BookingType) BookingType {
	if bt == nil {
		return UNKNOWN_BOOKINGTYPE
	}

	switch *bt {
	case models.STANDARD_BOOKINGTYPE:
		return STANDARD_BOOKINGTYPE
	case models.RESERVED_BOOKINGTYPE:
		return RESERVED_BOOKINGTYPE
	default:
		return UNKNOWN_BOOKINGTYPE
	}
}

func (g *GraphHelper) GetRooms() ([]Room, error) {
	r, err := g.userClient.Places().GraphRoom().Get(context.Background(), nil)
	if err != nil {
		printOdataError(err)
		return nil, fmt.Errorf("querying API: %+v", err)
	}

	pageIterator, err := msgraphcore.NewPageIterator[*models.Room](
		r, g.userClient.GetAdapter(), models.CreateRoomFromDiscriminatorValue,
	)
	if err != nil {
		return nil, fmt.Errorf("getting iterator: %v", err)
	}

	var rooms []Room
	if err := pageIterator.Iterate(context.Background(), func(msroom *models.Room) bool {
		if msroom == nil {
			return false
		}
		room := convertToRoom(*msroom)
		rooms = append(rooms, room)
		fmt.Printf("%+v:\n", room)
		// Return true to continue the iteration
		return true
	}); err != nil {
		return nil, fmt.Errorf("iterating: %v", err)
	}

	return rooms, nil
}

func convertToRoom(r models.Room) Room {
	address := r.GetAddress()
	geoCoordinates := r.GetGeoCoordinates()

	room := Room{
		DisplayName:            r.GetDisplayName(),
		Phone:                  r.GetPhone(),
		AudioDeviceName:        r.GetAudioDeviceName(),
		BookingType:            mapBookingType(r.GetBookingType()),
		Building:               r.GetBuilding(),
		Capacity:               r.GetCapacity(),
		DisplayDeviceName:      r.GetDisplayDeviceName(),
		EmailAddress:           r.GetEmailAddress(),
		FloorLabel:             r.GetFloorLabel(),
		FloorNumber:            r.GetFloorNumber(),
		IsWheelChairAccessible: r.GetIsWheelChairAccessible(),
		Label:                  r.GetLabel(),
		Nickname:               r.GetNickname(),
		Tags:                   r.GetTags(), // assuming it always returns a non-nil slice
		VideoDeviceName:        r.GetVideoDeviceName(),
	}

	if geoCoordinates != nil {
		room.GeoCoordinates = GeoCoordinates{
			Accuracy:         geoCoordinates.GetAccuracy(),
			Altitude:         geoCoordinates.GetAltitude(),
			AltitudeAccuracy: geoCoordinates.GetAltitudeAccuracy(),
			Latitude:         geoCoordinates.GetLatitude(),
			Longitude:        geoCoordinates.GetLongitude(),
		}
	}

	if address != nil {
		room.Address = PhysicalAddress{
			City:            address.GetCity(),
			CountryOrRegion: address.GetCountryOrRegion(),
			PostalCode:      address.GetPostalCode(),
			State:           address.GetState(),
			Street:          address.GetStreet(),
		}
	}
	return room
}

func printOdataError(err error) {
	switch err.(type) {
	case *odataerrors.ODataError:
		typed := err.(*odataerrors.ODataError)
		fmt.Printf("error: %v\n", typed.Error())
		if terr := typed.GetError(); terr != nil {
			fmt.Printf("code: %s\n", *terr.GetCode())
			fmt.Printf("msg: %s\n", *terr.GetMessage())
		}
	default:
		fmt.Printf("%T > error: %#v", err, err)
	}
}
