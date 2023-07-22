package msgraph

import (
	"context"
	"fmt"
	"microsoft-365/apiserver"
	"time"

	azcore "github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	auth "github.com/microsoft/kiota-authentication-azure-go"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
	"github.com/microsoftgraph/msgraph-sdk-go/users"

	"github.com/eliona-smart-building-assistant/go-eliona/utils"
	"github.com/eliona-smart-building-assistant/go-utils/common"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

type GraphHelper struct {
	credential      azcore.TokenCredential
	userClient      *msgraphsdk.GraphServiceClient
	graphUserScopes []string
	isDelegated     bool
}

func NewGraphHelper() *GraphHelper {
	g := &GraphHelper{}
	return g
}

func (g *GraphHelper) InitializeGraph(clientId, tenantId, clientSecret, username, password string) error {
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
		g.isDelegated = true
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
		g.isDelegated = false
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

type GraphAsset interface {
	getEmailAddress() *string
	setOnSchedule(*string)
}

type Room struct {
	Address                PhysicalAddress `eliona:"address,filterable" subtype:"info"`
	DisplayName            *string         `eliona:"display_name,filterable" subtype:"info"`
	Nickname               *string         `eliona:"nickname,filterable" subtype:"info"`
	Label                  *string         `eliona:"label,filterable" subtype:"info"`
	GeoCoordinates         GeoCoordinates  `eliona:"geo_coordinates,filterable" subtype:"info"`
	Phone                  *string         `eliona:"phone,filterable" subtype:"info"`
	EmailAddress           *string         `eliona:"email_address,filterable" subtype:"info"`
	BookingType            BookingType     `eliona:"booking_type,filterable" subtype:"info"`
	Building               *string         `eliona:"building,filterable" subtype:"info"`
	Capacity               *int32          `eliona:"capacity,filterable" subtype:"info"`
	FloorLabel             *string         `eliona:"floor_label,filterable" subtype:"info"`
	FloorNumber            *int32          `eliona:"floor_number,filterable" subtype:"info"`
	IsWheelChairAccessible *bool           `eliona:"is_wheel_chair_accessible,filterable" subtype:"info"`
	Tags                   []string        `eliona:"tags,filterable" subtype:"info"`
	DisplayDeviceName      *string         `eliona:"display_device_name,filterable" subtype:"info"`
	AudioDeviceName        *string         `eliona:"audio_device_name,filterable" subtype:"info"`
	VideoDeviceName        *string         `eliona:"video_device_name,filterable" subtype:"info"`
	OnSchedule             *string         `eliona:"on_schedule" subtype:"input"`
}

func (room Room) AssetType() string {
	return "microsoft_365_room"
}

func (room Room) Id() string {
	return room.AssetType() + "_" + *room.EmailAddress
}

func (room *Room) AdheresToFilter(config apiserver.Configuration) (bool, error) {
	f := apiFilterToCommonFilter(config.AssetFilter)
	fp, err := utils.StructToMap(room)
	if err != nil {
		return false, fmt.Errorf("converting struct to map: %v", err)
	}
	adheres, err := common.Filter(f, fp)
	if err != nil {
		return false, err
	}
	return adheres, nil
}

func (r *Room) getEmailAddress() *string {
	return r.EmailAddress
}

func (r *Room) setOnSchedule(s *string) {
	r.OnSchedule = s
}

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

func (g *GraphHelper) GetRooms(config apiserver.Configuration) ([]Room, error) {
	r, err := g.userClient.Places().GraphRoom().Get(context.Background(), nil)
	if err != nil {
		printOdataError(err)
		return nil, fmt.Errorf("querying rooms API: %+v", err)
	}

	pageIterator, err := msgraphcore.NewPageIterator[*models.Room](
		r, g.userClient.GetAdapter(), models.CreateRoomFromDiscriminatorValue,
	)
	if err != nil {
		return nil, fmt.Errorf("getting room iterator: %v", err)
	}

	rooms := make(map[string]*Room)
	if err := pageIterator.Iterate(context.Background(), func(msroom *models.Room) bool {
		if msroom == nil {
			return false
		}
		room := convertToRoom(*msroom)
		adheres, err := room.AdheresToFilter(config)
		if err != nil {
			log.Error("microsoft-365", "checking if room adheres to a filter: %v", err)
			return false
		}
		if !adheres {
			log.Debug("microsoft-365", "Room %s skipped.", *room.EmailAddress)
			return true
		}
		rooms[*room.EmailAddress] = &room
		// Return true to continue the iteration
		return true
	}); err != nil {
		return nil, fmt.Errorf("iterating rooms: %v", err)
	}
	if len(rooms) == 0 {
		return []Room{}, nil
	}

	rooms, err = fetchSchedules(g, rooms)
	if err != nil {
		return nil, fmt.Errorf("fetching schedules: %v", err)
	}

	var roomsSlice []Room
	for _, room := range rooms {
		roomsSlice = append(roomsSlice, *room)
	}
	return roomsSlice, nil
}

type Equipment struct {
	EmailAddress *string `eliona:"email_address,filterable" subtype:"info"`
	DisplayName  *string `eliona:"display_name,filterable" subtype:"info"`
	OnSchedule   *string `eliona:"on_schedule" subtype:"input"`
}

func (equipment Equipment) AssetType() string {
	return "microsoft_365_equipment"
}

func (equipment Equipment) Id() string {
	return equipment.AssetType() + "_" + *equipment.EmailAddress
}

func (equipment *Equipment) AdheresToFilter(config apiserver.Configuration) (bool, error) {
	f := apiFilterToCommonFilter(config.AssetFilter)
	fp, err := utils.StructToMap(equipment)
	if err != nil {
		return false, fmt.Errorf("converting struct to map: %v", err)
	}
	adheres, err := common.Filter(f, fp)
	if err != nil {
		return false, err
	}
	return adheres, nil
}

func (e *Equipment) getEmailAddress() *string {
	return e.EmailAddress
}

func (e *Equipment) setOnSchedule(s *string) {
	e.OnSchedule = s
}

func (g *GraphHelper) GetEquipment(config apiserver.Configuration) ([]Equipment, error) {
	// It would be wonderful if this filter worked. For some reason, mailboxSettings
	// can be accessed only user by user. See
	// https://feedbackportal.microsoft.com/feedback/idea/65df37d3-8f21-ee11-a81c-002248510ddf
	// for more info.
	//
	// f := "mailboxSettings/userPurpose eq 'equipment'"
	// requestParameters := &users.UsersRequestBuilderGetQueryParameters{
	// 	Select: []string{"id", "displayName", "mail", "mailboxSettings"},
	// 	Filter: &f,
	// }
	r, err := g.userClient.Users().Get(context.Background(), nil)
	if err != nil {
		printOdataError(err)
		return nil, fmt.Errorf("querying users API: %+v", err)
	}

	pageIterator, err := msgraphcore.NewPageIterator[*models.User](
		r, g.userClient.GetAdapter(), models.CreateUserFromDiscriminatorValue,
	)
	if err != nil {
		return nil, fmt.Errorf("getting users iterator: %v", err)
	}

	equipment := make(map[string]*Equipment)
	if err := pageIterator.Iterate(context.Background(), func(msuser *models.User) bool {
		if msuser == nil {
			return false
		}
		name := *msuser.GetUserPrincipalName()
		r, err := g.userClient.Users().ByUserId(name).MailboxSettings().Get(context.Background(), nil)
		if err != nil {
			printOdataError(err)
			log.Error("microsoft-365", "querying users API: %v", err)
			return true
		}
		purpose := r.GetUserPurpose()
		if *purpose != models.EQUIPMENT_USERPURPOSE {
			return true
		}

		e := convertToEquipment(*msuser)

		adheres, err := e.AdheresToFilter(config)
		if err != nil {
			log.Error("microsoft-365", "checking if equipment adheres to a filter: %v", err)
			return false
		}
		if !adheres {
			log.Debug("microsoft-365", "Room %s skipped.", *e.EmailAddress)
			return true
		}
		equipment[*e.EmailAddress] = &e
		// Return true to continue the iteration
		return true
	}); err != nil {
		return nil, fmt.Errorf("iterating equipment: %v", err)
	}
	if len(equipment) == 0 {
		return []Equipment{}, nil
	}

	equipment, err = fetchSchedules(g, equipment)
	if err != nil {
		return nil, fmt.Errorf("fetching schedules: %v", err)
	}

	var equipmentSlice []Equipment
	for _, e := range equipment {
		equipmentSlice = append(equipmentSlice, *e)
	}
	return equipmentSlice, nil
}

func fetchSchedules[T GraphAsset](g *GraphHelper, rooms map[string]T) (map[string]T, error) {
	var addressList []string
	for i := range rooms {
		addressList = append(addressList, i)
	}
	timeZone := "W. Europe Standard Time"
	headers := abstractions.NewRequestHeaders()
	// If not specified, values returned are in UTC.
	headers.Add("Prefer", fmt.Sprintf("outlook.timezone=\"%s\"", timeZone))

	configuration := &users.ItemCalendarGetScheduleRequestBuilderPostRequestConfiguration{
		Headers: headers,
	}
	requestBody := users.NewItemCalendarGetSchedulePostRequestBody()
	requestBody.SetSchedules(addressList)

	startTime := models.NewDateTimeTimeZone()
	t1 := time.Now()
	// The docs say "2006-01-02T15:04:05" is the correct format, but RFC3339
	// is acceptable as well (but undocumented).
	ts1 := t1.Format("2006-01-02T15:04:05")
	startTime.SetDateTime(&ts1)
	startTime.SetTimeZone(&timeZone)
	requestBody.SetStartTime(startTime)

	endTime := models.NewDateTimeTimeZone()

	t2 := t1.Add(time.Hour)

	ts2 := t2.Format("2006-01-02T15:04:05")
	endTime.SetDateTime(&ts2)
	endTime.SetTimeZone(&timeZone)
	requestBody.SetEndTime(endTime)

	availabilityViewInterval := int32(30)
	requestBody.SetAvailabilityViewInterval(&availabilityViewInterval)

	// POST https://graph.microsoft.com/v1.0/me/calendar/getSchedule
	// Prefer: outlook.timezone="W. Europe Standard Time"
	// Content-Type: application/json
	//
	// {
	//     "schedules": ["small.table@z0vmd.onmicrosoft.com", "silent.room@z0vmd.onmicrosoft.com"],
	//     "startTime": {
	//         "dateTime": "2019-03-15T09:00:00",
	//         "timeZone": "W. Europe Standard Time"
	//     },
	//     "endTime": {
	//         "dateTime": "2019-03-16T18:00:00",
	//         "timeZone": "W. Europe Standard Time"
	//     },
	//     "availabilityViewInterval": 60
	// }

	// ¯\_(ツ)_/¯
	//
	// While it is possible to get schedules of all the "users" in one query, the GetSchedule()
	// endpoint is accessible only through some user entity. It does not matter which, though,
	// so we can select one at random.
	//
	// Using Me() entity would be more elegant, but that entity is accessible only using
	// delegated permissions.
	//
	// Note: with delegated permissions, only the Me() endpoint is accessible so to test it in
	// Graph Explorer, use the query above this comment.
	var r users.ItemCalendarGetScheduleResponseable
	if g.isDelegated {
		var err error
		r, err = g.userClient.Me().Calendar().GetSchedule().Post(context.Background(), requestBody, configuration)
		if err != nil {
			printOdataError(err)
			return nil, fmt.Errorf("querying calendar API via delegated permission: %+v", err)
		}
	} else {
		randomAddress := addressList[0]
		var err error
		r, err = g.userClient.Users().ByUserId(randomAddress).Calendar().GetSchedule().Post(context.Background(), requestBody, configuration)
		if err != nil {
			printOdataError(err)
			return nil, fmt.Errorf("querying calendar API via app permission: %+v", err)
		}
	}

	pageIterator, err := msgraphcore.NewPageIterator[*models.ScheduleInformation](
		r, g.userClient.GetAdapter(), models.CreateScheduleInformationFromDiscriminatorValue,
	)
	if err != nil {
		return nil, fmt.Errorf("getting schedule iterator: %v", err)
	}

	if err := pageIterator.Iterate(context.Background(), func(schedule *models.ScheduleInformation) bool {
		if schedule == nil {
			return false
		}
		sID := schedule.GetScheduleId()
		if sID == nil {
			log.Debug("microsoft-365", "Empty schedule ID")
			return true
		}
		scheduleID := *sID
		fmt.Println(scheduleID)

		room := rooms[scheduleID]
		scheduleItems := schedule.GetScheduleItems()
		if scheduleItems == nil || len(scheduleItems) == 0 {
			room.setOnSchedule(nil)
		} else {
			d := getScheduleItemableDescription(scheduleItems[0])
			room.setOnSchedule(&d)
		}
		rooms[scheduleID] = room

		// Return true to continue the iteration.
		return true
	}); err != nil {
		return nil, fmt.Errorf("iterating schedules: %v", err)
	}
	return rooms, nil
}

func getScheduleItemableDescription(i models.ScheduleItemable) string {
	var result string
	if s := i.GetStatus(); s != nil {
		result = s.String()
	}
	if s := i.GetSubject(); s != nil {
		result = *s
	}
	return result
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

func convertToEquipment(u models.User) Equipment {
	return Equipment{
		DisplayName:  u.GetDisplayName(),
		EmailAddress: u.GetUserPrincipalName(),
	}
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
		fmt.Printf("%T > error: %#v\n", err, err)
	}
}
