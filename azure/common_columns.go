package azure

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/memoize"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// column definitions for the common columns
func commonColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "cloud_environment",
			Type:        proto.ColumnType_STRING,
			Hydrate:     getCloudEnvironment,
			Description: ColumnDescriptionCloudEnvironment,
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

// if the caching is required other than per connection, build a cache key for the call and use it in Memoize.
var getSubscriptionIDMemoized = plugin.HydrateFunc(getSubscriptionIDUncached).Memoize(memoize.WithCacheKeyFunction(getSubscriptionIDCacheKey))

// declare a wrapper hydrate function to call the memoized function
// - this is required when a memoized function is used for a column definition
func getSubscriptionID(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return getSubscriptionIDMemoized(ctx, d, h)
}

// Build a cache key for the call to getSubscriptionIDCacheKey.
func getSubscriptionIDCacheKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	key := "getSubscriptionID"
	return key, nil
}

func getSubscriptionIDUncached(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

// if the caching is required other than per connection, build a cache key for the call and use it in Memoize.
var getCloudEnvironmentMemoized = plugin.HydrateFunc(getCloudEnvironmentUncached).Memoize(memoize.WithCacheKeyFunction(getCloudEnvironmentCacheKey))

// declare a wrapper hydrate function to call the memoized function
// - this is required when a memoized function is used for a column definition
func getCloudEnvironment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return getCloudEnvironmentMemoized(ctx, d, h)
}

// Build a cache key for the call to getCloudEnvironment.
func getCloudEnvironmentCacheKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	key := "getCloudEnvironment"
	return key, nil
}

func getCloudEnvironmentUncached(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCloudEnvironment")
	cacheKey := "getCloudEnvironment"

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
