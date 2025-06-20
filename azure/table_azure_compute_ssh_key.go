package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION ////

func tableAzureComputeSshKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_compute_ssh_key",
		Description: "Azure Compute SSH Key",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getAzureComputeSshKey,
			Tags: map[string]string{
				"service": "Microsoft.Compute",
				"action":  "sshPublicKeys/read",
			},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "404"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listAzureComputeSshKeys,
			Tags: map[string]string{
				"service": "Microsoft.Compute",
				"action":  "sshPublicKeys/read",
			},
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "id",
				Description: "The unique ID identifying the resource in subscription.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "name",
				Description: "Name of the SSH key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The type of the resource in Azure.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "public_key",
				Description: "SSH public key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SSHPublicKeyResourceProperties.PublicKey"),
			},

			// Azure standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Location").Transform(toLower),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "resource_group",
				Description: ColumnDescriptionResourceGroup,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
			},

			// Steampipe standard columns
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
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},
		}),
	}
}

//// LIST FUNCTION ////

func listAzureComputeSshKeys(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := compute.NewSSHPublicKeysClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

	// Apply rate limiting
	d.WaitForListRateLimit(ctx)

	result, err := client.ListBySubscription(ctx)
	if err != nil {
		return nil, err
	}

	for _, key := range result.Values() {
		d.StreamListItem(ctx, key)
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

		for _, key := range result.Values() {
			d.StreamListItem(ctx, key)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTION ////

func getAzureComputeSshKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAzureComputeSshKey")

	name := d.EqualsQualString("name")
	resourceGroup := d.EqualsQualString("resource_group")

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_compute_ssh_key.getAzureComputeSshKey", "client_error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := compute.NewSSHPublicKeysClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &client, d.Connection)

	op, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		plugin.Logger(ctx).Error("azure_compute_ssh_key.getAzureComputeSshKey", "query_error", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}
