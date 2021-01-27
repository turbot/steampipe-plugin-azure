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

func resourceInterfaceDescription(key string) string {
	switch key {
	case "akas":
		return "Array of globally unique identifier strings (also known as) for the resource."
	case "tags":
		return "A map of tags for the resource."
	case "title":
		return "Title of the resource."
	}
	return ""
}
