package msgraph

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"template/apiserver"

	azcore "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	auth "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/microsoftgraph/msgraph-sdk-go/users"

	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/log"
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
	City            *string
	CountryOrRegion *string
	PostalCode      *string
	State           *string
	Street          *string
}

type GeoCoordinates struct {
	Accuracy         *float64
	Altitude         *float64
	AltitudeAccuracy *float64
	Latitude         *float64
	Longitude        *float64
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

type Room struct {
	Address                PhysicalAddress `eliona:"address" subtype:"info"`
	DisplayName            *string         `eliona:"display_name" subtype:"info"`
	Nickname               *string         `eliona:"nickname" subtype:"info"`
	Label                  *string         `eliona:"label" subtype:"info"`
	GeoCoordinates         GeoCoordinates  `eliona:"geo_coordinates" subtype:"info"`
	Phone                  *string         `eliona:"phone" subtype:"info"`
	EmailAddress           *string         `eliona:"email_address" subtype:"info"`
	BookingType            BookingType     `eliona:"booking_type" subtype:"info"`
	Building               *string         `eliona:"building" subtype:"info"`
	Capacity               *int32          `eliona:"capacity" subtype:"info"`
	FloorLabel             *string         `eliona:"floor_label" subtype:"info"`
	FloorNumber            *int32          `eliona:"floor_number" subtype:"info"`
	IsWheelChairAccessible *bool           `eliona:"is_wheel_chair_accessible" subtype:"info"`
	Tags                   []string        `eliona:"tags" subtype:"info"`
	DisplayDeviceName      *string         `eliona:"display_device_name" subtype:"info"`
	AudioDeviceName        *string         `eliona:"audio_device_name" subtype:"info"`
	VideoDeviceName        *string         `eliona:"video_device_name" subtype:"info"`
}

func (room *Room) AdheresToFilter(config apiserver.Configuration) (bool, error) {
	f := apiFilterToCommonFilter(config.AssetFilter)
	fp, err := structToMap(room)
	if err != nil {
		return false, fmt.Errorf("converting struct to map: %v", err)
	}
	adheres, err := common.Filter(f, fp)
	if err != nil {
		return false, err
	}
	return adheres, nil
}

func (g *GraphHelper) GetRooms(config apiserver.Configuration) ([]Room, error) {
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
		adheres, err := room.AdheresToFilter(config)
		if err != nil {
			log.Error("ms-graph", "checking if room adheres to a filter: %v", err)
			return false
		}
		if !adheres {
			log.Debug("ms-graph", "Room %s skipped.", room.EmailAddress)
			return true
		}
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

// To be moved to go-utils.

func structToMap(input interface{}) (map[string]string, error) {
	if input == nil {
		return nil, fmt.Errorf("input is nil")
	}

	inputValue := reflect.ValueOf(input)
	inputType := reflect.TypeOf(input)

	if inputValue.Kind() == reflect.Ptr {
		inputValue = inputValue.Elem()
		inputType = inputType.Elem()
	}

	if inputValue.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not a struct")
	}

	output := make(map[string]string)
	for i := 0; i < inputValue.NumField(); i++ {
		fieldType := inputType.Field(i)

		fieldTag, err := parseElionaTag(fieldType)
		if err != nil {
			return nil, err
		}

		if !fieldTag.Filterable {
			continue
		}

		fieldValue := inputValue.Field(i)
		output[fieldTag.ParamName] = fieldValue.String()
	}

	return output, nil
}

type SubType string

const (
	Status SubType = "status"
	Info   SubType = "info"
	Input  SubType = "input"
	Output SubType = "output"
)

type FieldTag struct {
	ParamName  string
	SubType    SubType
	Filterable bool
}

func parseElionaTag(field reflect.StructField) (*FieldTag, error) {
	elionaTag := field.Tag.Get("eliona")
	subtypeTag := field.Tag.Get("subtype")

	elionaTagParts := strings.Split(elionaTag, ",")
	if len(elionaTagParts) < 1 {
		return nil, fmt.Errorf("invalid eliona tag on field %s", field.Name)
	}

	paramName := elionaTagParts[0]
	filterable := len(elionaTagParts) > 1 && elionaTagParts[1] == "filterable"

	var subType SubType
	if subtypeTag != "" {
		subType = SubType(subtypeTag)
		switch subType {
		case Status, Info, Input, Output:
			// valid subtype
		default:
			return nil, fmt.Errorf("invalid subtype in eliona tag on field %s", field.Name)
		}
	}

	return &FieldTag{
		ParamName:  paramName,
		SubType:    subType,
		Filterable: filterable,
	}, nil
}

// ^^ To be moved to go-utils.

func apiFilterToCommonFilter(input [][]apiserver.FilterRule) [][]common.FilterRule {
	result := make([][]common.FilterRule, len(input))
	for i := 0; i < len(input); i++ {
		result[i] = make([]common.FilterRule, len(input[i]))
		for j := 0; j < len(input[i]); j++ {
			result[i][j] = common.FilterRule{
				Parameter: input[i][j].Parameter,
				Regex:     input[i][j].Regex,
			}
		}
	}
	return result
}
