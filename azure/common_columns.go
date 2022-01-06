package azure

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

// column definitions for the common columns
func commonColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "cloud_environment",
			Type:        proto.ColumnType_STRING,
			Hydrate:     getEnvironmentName,
			Description: ColumnDescriptionEnvironmentName,
			Transform:   transform.FromValue(),
		},
		{
			Name:        "subscription_id",
			Type:        proto.ColumnType_STRING,
			Hydrate:     getSubscriptionID,
			Description: ColumnDescriptionSubscription,
			Transform:   transform.FromValue(),
		},
	}
}

// append the common azure columns onto the column list
func azureColumns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, commonColumns()...)
}

func getSubscriptionID(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSubscriptionID")
	cacheKey := "getSubscriptionID"

	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(string), nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	// cache subscription id for the session
	d.ConnectionManager.Cache.Set(cacheKey, session.SubscriptionID)

	return session.SubscriptionID, nil
}

func getEnvironmentName(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getEnvironmentName")
	cacheKey := "getEnvironmentName"

	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(string), nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}

	// cache environment name for the session
	d.ConnectionManager.Cache.Set(cacheKey, session.CloudEnvironment)

	return session.CloudEnvironment, nil
}
