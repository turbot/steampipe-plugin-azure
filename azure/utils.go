package azure

import (
	"context"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TRANSFORM FUNCTION ////

func idToSubscriptionID(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	id := types.SafeString(d.Value)
	subscriptionid := strings.Split(id, "/")[2]
	return subscriptionid, nil
}

func idToAkas(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	id := types.SafeString(d.Value)
	akas := []string{"azure://" + id, "azure://" + strings.ToLower(id)}
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

func convertDateToTime(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	dateValue := d.Value.(*date.Time)

	if dateValue != nil {
		// convert from *date.Time to *date.Time
		timeValue := dateValue.ToTime().Format(time.RFC3339)

		return timeValue, nil
	}

	return nil, nil
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
	ColumnDescriptionAkas          = "Array of globally unique identifier strings (also known as) for the resource."
	ColumnDescriptionRegion        = "The Azure region/location in which the resource is located."
	ColumnDescriptionResourceGroup = "The resource group which holds this resource."
	ColumnDescriptionSubscription  = "The Azure Subscription ID in which the resource is located."
	ColumnDescriptionTags          = "A map of tags for the resource."
	ColumnDescriptionTitle         = "Title of the resource."
)

// convert string to lower case
func toLower(_ context.Context, d *transform.TransformData) (interface{}, error) {
	valStr := types.SafeString(d.Value)
	return strings.ToLower(valStr), nil
}
