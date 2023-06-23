/*
 * MS Graph App
 *
 * API to access and configure the MS Graph App
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apiserver

// WidgetData - Data for a widget
type WidgetData struct {

	// The internal Id of widget data
	Id *int32 `json:"id,omitempty"`

	// Position of the element in widget type
	ElementSequence *int32 `json:"elementSequence,omitempty"`

	// The master asset id of this widget
	AssetId *int32 `json:"assetId,omitempty"`

	// individual config parameters depending on category
	Data *map[string]interface{} `json:"data,omitempty"`
}

// AssertWidgetDataRequired checks if the required fields are not zero-ed
func AssertWidgetDataRequired(obj WidgetData) error {
	return nil
}

// AssertRecurseWidgetDataRequired recursively checks if required fields are not zero-ed in a nested slice.
// Accepts only nested slice of WidgetData (e.g. [][]WidgetData), otherwise ErrTypeAssertionError is thrown.
func AssertRecurseWidgetDataRequired(objSlice interface{}) error {
	return AssertRecurseInterfaceRequired(objSlice, func(obj interface{}) error {
		aWidgetData, ok := obj.(WidgetData)
		if !ok {
			return ErrTypeAssertionError
		}
		return AssertWidgetDataRequired(aWidgetData)
	})
}
