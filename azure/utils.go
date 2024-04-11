package azure

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TRANSFORM FUNCTIONS

func idToAkas(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	id := types.SafeString(d.Value)
	akas := []string{"azure://" + id, "azure://" + strings.ToLower(id)}
	occured := map[string]bool{}
	result := []string{}
	for i := range akas {
		if !occured[akas[i]] {
			occured[akas[i]] = true
			result = append(result, akas[i])
		}
	}
	akas = result
	return akas, nil
}

func extractResourceGroupFromID(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	id := types.SafeString(d.Value)

	// Common resource properties
	splitID := strings.Split(id, "/")
	resourceGroup := splitID[4]
	resourceGroup = strings.ToLower(resourceGroup)
	return resourceGroup, nil
}

func lastPathElement(_ context.Context, d *transform.TransformData) (interface{}, error) {
	return getLastPathElement(types.SafeString(d.Value)), nil
}

func getLastPathElement(path string) string {
	if path == "" {
		return ""
	}

	pathItems := strings.Split(path, "/")
	return pathItems[len(pathItems)-1]
}

func convertDateToTime(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}
	dateValue := d.Value.(*date.Time)

	if dateValue != nil {
		// convert from *date.Time to *date.Time
		timeValue := dateValue.ToTime().Format(time.RFC3339)

		return timeValue, nil
	}

	return nil, nil
}

func structToMap(val reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		// Check if field is a struct and not a zero value
		if fieldValue.Kind() == reflect.Struct && !fieldValue.IsZero() {
			result[field.Name] = structToMap(fieldValue)
		} else if !fieldValue.IsZero() {
			result[field.Name] = fieldValue.Interface()
		} else {
			result[field.Name] = nil
		}
	}

	return result
}

func convertDateUnixToTime(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	dateValue := d.Value.(*date.UnixTime)
	if dateValue != nil {
		// convert from *date.Time to *date.Time
		timeValue := dateValue.Duration().Milliseconds()

		epochTime, err := types.ToInt64(timeValue)
		if err != nil {
			return nil, err
		}
		if epochTime == 0 {
			return nil, nil
		}
		timeIn := time.Unix(0, epochTime*int64(time.Millisecond))
		timestampRFC3339Format := timeIn.Format(time.RFC3339)
		return timestampRFC3339Format, nil
	}

	return nil, nil
}

// Constants for Standard Column Descriptions
const (
	ColumnDescriptionAkas             = "Array of globally unique identifier strings (also known as) for the resource."
	ColumnDescriptionCloudEnvironment = "The Azure Cloud Environment."
	ColumnDescriptionRegion           = "The Azure region/location in which the resource is located."
	ColumnDescriptionResourceGroup    = "The resource group which holds this resource."
	ColumnDescriptionSubscription     = "The Azure Subscription ID in which the resource is located."
	ColumnDescriptionTags             = "A map of tags for the resource."
	ColumnDescriptionTitle            = "Title of the resource."
)

// convert string to lower case
func toLower(_ context.Context, d *transform.TransformData) (interface{}, error) {
	valStr := types.SafeString(d.Value)
	return strings.ToLower(valStr), nil
}

// Remove spaces in between the location (i.e. 'east us' will be formatted to 'eastus')
func formatRegion(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	valStr := types.SafeString(d.Value)
	region := strings.ReplaceAll(valStr, " ", "")
	return region, nil
}
