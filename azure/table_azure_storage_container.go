package azure

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
)

//// TABLE DEFINITION

func tableAzureStorageContainer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_storage_container",
		Description: "Azure Storage Container",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group", "account_name"}),
			Hydrate:    getStorageContainer,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "ContainerNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listStorageAccounts,
			Hydrate:       listStorageContainers,
		},

		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the container.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a container uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "account_name",
				Description: "The friendly name that identifies the storage account.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(idToAccountName),
			},
			{
				Name:        "deleted",
				Description: "Indicates whether the blob container was deleted.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ContainerProperties.Deleted"),
			},
			{
				Name:        "public_access",
				Description: "Specifies whether data in the container may be accessed publicly and the level of access.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContainerProperties.PublicAccess").Transform(transform.ToString),
			},
			{
				Name:        "type",
				Description: "Specifies the type of the container.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_encryption_scope",
				Description: "Default the container to use specified encryption scope for all writes.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContainerProperties.DefaultEncryptionScope"),
			},
			{
				Name:        "deleted_time",
				Description: "Specifies the time when the container was deleted.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ContainerProperties.DeletedTime").Transform(convertDateToTime),
			},
			{
				Name:        "deny_encryption_scope_override",
				Description: "Indicates whether block override of encryption scope from the container default, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ContainerProperties.DenyEncryptionScopeOverride"),
			},
			{
				Name:        "has_immutability_policy",
				Description: "The hasImmutabilityPolicy public property is set to true by SRP if ImmutabilityPolicy has been created for this container. The hasImmutabilityPolicy public property is set to false by SRP if ImmutabilityPolicy has not been created for this container.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ContainerProperties.HasImmutabilityPolicy"),
			},
			{
				Name:        "has_legal_hold",
				Description: "The hasLegalHold public property is set to true by SRP if there are at least one existing tag. The hasLegalHold public property is set to false by SRP if all existing legal hold tags are cleared out. There can be a maximum of 1000 blob containers with hasLegalHold=true for a given account.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ContainerProperties.HasLegalHold"),
			},
			{
				Name:        "last_modified_time",
				Description: "Specifies the date and time the container was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ContainerProperties.LastModifiedTime").Transform(convertDateToTime),
			},
			{
				Name:        "lease_status",
				Description: "Specifies the lease status of the container.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContainerProperties.LeaseStatus").Transform(transform.ToString),
			},
			{
				Name:        "lease_state",
				Description: "Specifies the lease state of the container.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContainerProperties.LeaseState").Transform(transform.ToString),
			},
			{
				Name:        "lease_duration",
				Description: "Specifies whether the lease on a container is of infinite or fixed duration, only when the container is leased. Possible values are: 'Infinite', 'Fixed'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContainerProperties.LeaseDuration").Transform(transform.ToString),
			},
			{
				Name:        "remaining_retention_days",
				Description: "Remaining retention days for soft deleted blob container.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ContainerProperties.RemainingRetentionDays"),
			},
			{
				Name:        "version",
				Description: "The version of the deleted blob container.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ContainerProperties.Version"),
			},
			{
				Name:        "immutability_policy",
				Description: "The ImmutabilityPolicy property of the container.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getImmutabilityPolicy,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "legal_hold",
				Description: "The LegalHold property of the container.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerProperties.LegalHold"),
			},
			{
				Name:        "metadata",
				Description: "A name-value pair to associate with the container as metadata.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ContainerProperties.Metadata"),
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

			// Azure standard columns
			{
				Name:        "resource_group",
				Description: ColumnDescriptionResourceGroup,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID),
			},
		}),
	}
}

//// LIST FUNCTION

func listStorageContainers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of storage account
	account := h.Item.(*storageAccountInfo)

	// Blob is not supported for the account if storage type is FileStorage
	if account.Account.Kind == "FileStorage" {
		return nil, nil
	}

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := storage.NewBlobContainersClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(ctx, *account.ResourceGroup, *account.Name, "", "", "")
	if err != nil {
		return nil, err
	}
	for _, container := range result.Values() {
		d.StreamListItem(ctx, container)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, container := range result.Values() {
			d.StreamListItem(ctx, container)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getStorageContainer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getStorageContainer")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()
	accountName := d.KeyColumnQuals["account_name"].GetStringValue()

	client := storage.NewBlobContainersClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func getImmutabilityPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getImmutabilityPolicy")
	container := h.Item.(storage.ListContainerItem)

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	resourceGroup := strings.Split(*container.ID, "/")[4]
	accountName := strings.Split(*container.ID, "/")[8]

	client := storage.NewBlobContainersClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.GetImmutabilityPolicy(ctx, resourceGroup, accountName, *container.Name, "")
	if err != nil {
		return nil, err
	}

	// If we return the API response directly, the output only gives top level
	// contents of ImmutabilityPolicy
	ImmutabilityPolicy := make(map[string]interface{})
	if op.ID != nil {
		ImmutabilityPolicy["ID"] = op.ID
	}
	if op.Name != nil {
		ImmutabilityPolicy["Name"] = op.Name
	}
	if op.Type != nil {
		ImmutabilityPolicy["Type"] = op.Type
	}
	if op.ImmutabilityPolicyProperty != nil {
		ImmutabilityPolicy["AllowProtectedAppendWrites"] = op.ImmutabilityPolicyProperty.AllowProtectedAppendWrites
		ImmutabilityPolicy["ImmutabilityPeriodSinceCreationInDays"] = op.ImmutabilityPolicyProperty.ImmutabilityPeriodSinceCreationInDays
		ImmutabilityPolicy["State"] = op.ImmutabilityPolicyProperty.State
	}

	return ImmutabilityPolicy, nil
}

//// TRANSFORM FUNCTIONS

func idToAccountName(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	id := types.SafeString(d.Value)
	accountName := strings.Split(id, "/")[8]
	return accountName, nil
}
