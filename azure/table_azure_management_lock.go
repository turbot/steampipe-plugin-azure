package azure

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2016-09-01/locks"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
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
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			ItemFromKey:       managementLockDataFromKey,
			Hydrate:           getManagementLock,
			ShouldIgnoreError: isNotFoundError([]string{"LockNotFound"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listManagementLocks,
		},

		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The friendly name that identifies management lock",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a lock uniquely",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagementLockObject.ID"),
			},
			{
				Name:        "type",
				Description: "The resource type of the lock",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagementLockObject.Type"),
			},
			{
				Name:        "lock_level",
				Description: "The level of the lock",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagementLockObject.ManagementLockProperties.Level").Transform(transform.ToString),
			},
			{
				Name:        "scope",
				Description: "Contains the scope of the lock",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getAzureManagementLockScope),
			},
			{
				Name:        "notes",
				Description: "Contains the notes about the lock",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagementLockObject.ManagementLockProperties.Notes"),
			},
			{
				Name:        "owners",
				Description: "A list of owners of the lock",
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
			{
				Name:        "resource_group",
				Type:        proto.ColumnType_STRING,
				Description: ColumnDescriptionResourceGroup,
				Transform:   transform.FromField("ResourceGroup").Transform(toLower),
			},
			{
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagementLockObject.ID").Transform(idToSubscriptionID),
			},
		},
	}
}

//// ITEM FROM KEY

func managementLockDataFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	name := quals["name"].GetStringValue()
	resourceGroup := quals["resource_group"].GetStringValue()
	item := &managementLockInfo{
		Name:          &name,
		ResourceGroup: &resourceGroup,
	}
	return item, nil
}

//// LIST FUNCTION

func listManagementLocks(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	locksClient := locks.NewManagementLocksClient(subscriptionID)
	locksClient.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		result, err := locksClient.ListAtSubscriptionLevel(ctx, subscriptionID)
		if err != nil {
			return nil, err
		}

		for _, managementLock := range result.Values() {
			resourceGroup := &strings.Split(string(*managementLock.ID), "/")[4]
			d.StreamListItem(ctx, managementLockInfo{managementLock, managementLock.Name, resourceGroup})
		}
		result.NextWithContext(context.Background())
		pagesLeft = result.NotDone()
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getManagementLock(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	managementLock := h.Item.(*managementLockInfo)

	session, err := GetNewSession(ctx, d.ConnectionManager, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	locksClient := locks.NewManagementLocksClient(subscriptionID)
	locksClient.Authorizer = session.Authorizer

	filter := fmt.Sprintf("name eq '%s'", *managementLock.Name)
	op, err := locksClient.ListAtResourceGroupLevel(ctx, *managementLock.ResourceGroup, filter)
	if err != nil {
		return nil, err
	}

	if op.Values() != nil && len(op.Values()) > 0 {
		return managementLockInfo{op.Values()[0], op.Values()[0].Name, managementLock.ResourceGroup}, nil
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func getAzureManagementLockScope(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(managementLockInfo)
	return strings.Split(string(*data.ManagementLockObject.ID), "/providers/Microsoft.Authorization/locks/")[0], nil
}
