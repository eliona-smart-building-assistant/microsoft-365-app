/*
 * Microsoft 365 App
 *
 * API to access and configure the Microsoft 365 App
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package apiserver

// GetDashboardTemplateByName200ResponseWidgetsInner - A widget on a frontend dashboard
type GetDashboardTemplateByName200ResponseWidgetsInner struct {

	// The internal Id of widget
	Id int32 `json:"id,omitempty"`

	// The name for the type of this widget
	WidgetTypeName string `json:"widgetTypeName"`

	// Detailed configuration depending on the widget type
	Details map[string]interface{} `json:"details,omitempty"`

	// The master asset id of this widget
	AssetId int32 `json:"assetId,omitempty"`

	// Placement order on dashboard; if not set the index in array is taken
	Sequence int32 `json:"sequence,omitempty"`

	Data []GetDashboardTemplateByName200ResponseWidgetsInnerDataInner `json:"data,omitempty"`
}

// AssertGetDashboardTemplateByName200ResponseWidgetsInnerRequired checks if the required fields are not zero-ed
func AssertGetDashboardTemplateByName200ResponseWidgetsInnerRequired(obj GetDashboardTemplateByName200ResponseWidgetsInner) error {
	elements := map[string]interface{}{
		"widgetTypeName": obj.WidgetTypeName,
	}
	for name, el := range elements {
		if isZero := IsZeroValue(el); isZero {
			return &RequiredError{Field: name}
		}
	}

	for _, el := range obj.Data {
		if err := AssertGetDashboardTemplateByName200ResponseWidgetsInnerDataInnerRequired(el); err != nil {
			return err
		}
	}
	return nil
}

// AssertGetDashboardTemplateByName200ResponseWidgetsInnerConstraints checks if the values respects the defined constraints
func AssertGetDashboardTemplateByName200ResponseWidgetsInnerConstraints(obj GetDashboardTemplateByName200ResponseWidgetsInner) error {
	return nil
}
