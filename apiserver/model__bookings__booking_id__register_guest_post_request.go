/*
 * Microsoft 365 App
 *
 * API to access and configure the Microsoft 365 App
 *
 * API version: 1.1.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apiserver

type BookingsBookingIdRegisterGuestPostRequest struct {

	// The Eliona user to send the message to.
	NotificationRecipient string `json:"notificationRecipient,omitempty"`

	// The content of the message to be sent to the organizer.
	MessageEn string `json:"messageEn"`

	// The content of the message to be sent to the organizer.
	MessageDe string `json:"messageDe,omitempty"`

	// The content of the message to be sent to the organizer.
	MessageFr string `json:"messageFr,omitempty"`

	// The content of the message to be sent to the organizer.
	MessageIt string `json:"messageIt,omitempty"`
}

// AssertBookingsBookingIdRegisterGuestPostRequestRequired checks if the required fields are not zero-ed
func AssertBookingsBookingIdRegisterGuestPostRequestRequired(obj BookingsBookingIdRegisterGuestPostRequest) error {
	elements := map[string]interface{}{
		"messageEn": obj.MessageEn,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	return nil
}

// AssertBookingsBookingIdRegisterGuestPostRequestConstraints checks if the values respects the defined constraints
func AssertBookingsBookingIdRegisterGuestPostRequestConstraints(obj BookingsBookingIdRegisterGuestPostRequest) error {
	return nil
}
