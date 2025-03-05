package azure

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/locks"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type managementLockInfo = struct {
	ManagementLockObject locks.ManagementLockObject
	Name                 *string
	ResourceGroup        *string
}

//// TABLE DEFINITION

func tableAzureManagementLock(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_management_lock",
		Description: "Azure Management Lock",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getManagementLock,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"LockNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listManagementLocks,
		},

		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies management lock.",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a lock uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagementLockObject.ID"),
			},
			{
				Name:        "type",
				Description: "The resource type of the lock.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagementLockObject.Type"),
			},
			{
				Name:        "lock_level",
				Description: "The level of the lock.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagementLockObject.ManagementLockProperties.Level").Transform(transform.ToString),
			},
			{
				Name:        "scope",
				Description: "Contains the scope of the lock.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getAzureManagementLockScope),
			},
			{
				Name:        "notes",
				Description: "Contains the notes about the lock.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagementLockObject.ManagementLockProperties.Notes"),
			},
			{
				Name:        "owners",
				Description: "A list of owners of the lock.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ManagementLockObject.ManagementLockProperties.Owners"),
			},
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ManagementLockObject.ID").Transform(idToAkas),
			},

			// Azure standard columns
			{
				Name:        "resource_group",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionResourceGroup,
				Transform:   transform.FromField("ResourceGroup").Transform(toLower),
			},
		}),
	}
}

//// LIST FUNCTION

func listManagementLocks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	locksClient := locks.NewManagementLocksClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	locksClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &locksClient, d.Connection)

	result, err := locksClient.ListAtSubscriptionLevel(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	for _, managementLock := range result.Values() {
		resourceGroup := &strings.Split(string(*managementLock.ID), "/")[4]
		d.StreamListItem(ctx, managementLockInfo{managementLock, managementLock.Name, resourceGroup})
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, managementLock := range result.Values() {
			resourceGroup := &strings.Split(string(*managementLock.ID), "/")[4]
			d.StreamListItem(ctx, managementLockInfo{managementLock, managementLock.Name, resourceGroup})
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getManagementLock(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getManagementLock")

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	name := d.EqualsQuals["name"].GetStringValue()
	resourceGroup := d.EqualsQuals["resource_group"].GetStringValue()

	locksClient := locks.NewManagementLocksClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	locksClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &locksClient, d.Connection)

	filter := fmt.Sprintf("name eq '%s'", name)
	op, err := locksClient.ListAtResourceGroupLevel(ctx, resourceGroup, filter)
	if err != nil {
		return nil, err
	}

	if op.Values() != nil && len(op.Values()) > 0 {
		return managementLockInfo{op.Values()[0], op.Values()[0].Name, &resourceGroup}, nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func getAzureManagementLockScope(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(managementLockInfo)
	return strings.Split(string(*data.ManagementLockObject.ID), "/providers/Microsoft.Authorization/locks/")[0], nil
}
