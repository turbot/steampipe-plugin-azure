package azure

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
)

//// TABLE DEFINITION

func tableAzureStorageContainer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_storage_container",
		Description: "Azure Storage Container",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group", "account_name"}),
			Hydrate:           getStorageContainer,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "ContainerNotFound"}),
		},
		List: &plugin.ListConfig{
			ParentHydrate: listStorageAccounts,
			Hydrate:       listStorageContainers,
		},

		Columns: []*plugin.Column{
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
				Transform:   transform.FromField("ContainerProperties.DeletedTime"),
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
				Transform:   transform.FromField("ContainerProperties.LastModifiedTime"),
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
				Transform:   transform.FromField("ContainerProperties.ImmutabilityPolicy"),
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
			{
				Name:        "subscription_id",
				Description: ColumnDescriptionSubscription,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(idToSubscriptionID),
			},
		},
	}
}

//// LIST FUNCTION

func listStorageContainers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get the details of storage account
	account := h.Item.(*storageAccountInfo)

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := storage.NewBlobContainersClient(subscriptionID)
	client.Authorizer = session.Authorizer

	pagesLeft := true
	for pagesLeft {
		containerList, err := client.List(ctx, *account.ResourceGroup, *account.Name, "", "", "")
		if err != nil {
			return nil, err
		}

		for _, container := range containerList.Values() {
			d.StreamLeafListItem(ctx, container)
		}
		containerList.NextWithContext(context.Background())
		pagesLeft = containerList.NotDone()
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

	client := storage.NewBlobContainersClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, accountName, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

//// TRANSFORM FUNCTIONS

func idToAccountName(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	id := types.SafeString(d.Value)
	accountName := strings.Split(id, "/")[8]
	return accountName, nil
}
