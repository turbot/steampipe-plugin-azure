package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/containerregistry/mgmt/2020-11-01-preview/containerregistry"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureContainerRegistry(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_container_registry",
		Description: "Azure Container Registry",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getContainerRegistry,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceGroupNotFound", "ResourceNotFound", "Invalid input", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listContainerRegistries,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "The unique id identifying the resource in subscription.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The provisioning state of the container registry at the time the operation was called. Valid values are: 'Creating', 'Updating', 'Deleting', 'Succeeded', 'Failed', 'Canceled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RegistryProperties.ProvisioningState").Transform(transform.ToString),
			},
			{
				Name:        "creation_date",
				Description: "The creation date of the container registry.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("RegistryProperties.CreationDate").Transform(convertDateToTime),
			},
			{
				Name:        "admin_user_enabled",
				Description: "Indicates whether the admin user is enabled, or not.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("RegistryProperties.AdminUserEnabled"),
			},
			{
				Name:        "data_endpoint_enabled",
				Description: "Enable a single data endpoint per region for serving data.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("RegistryProperties.DataEndpointEnabled"),
			},
			{
				Name:        "login_server",
				Description: "The URL that can be used to log into the container registry.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RegistryProperties.LoginServer"),
			},
			{
				Name:        "network_rule_bypass_options",
				Description: "Indicates whether to allow trusted Azure services to access a network restricted registry. Valid values are: 'AzureServices', 'None'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RegistryProperties.NetworkRuleBypassOptions").Transform(transform.ToString),
			},
			{
				Name:        "public_network_access",
				Description: "Indicates whether or not public network access is allowed for the container registry. Valid values are: 'Enabled', 'Disabled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RegistryProperties.PublicNetworkAccess").Transform(transform.ToString),
			},
			{
				Name:        "sku_name",
				Description: "The SKU name of the container registry. Required for registry creation. Valid values are: 'Classic', 'Basic', 'Standard', 'Premium'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name").Transform(transform.ToString),
			},
			{
				Name:        "sku_tier",
				Description: "The SKU tier based on the SKU name. Valid values are: 'Classic', 'Basic', 'Standard', 'Premium'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Tier").Transform(transform.ToString),
			},
			{
				Name:        "status",
				Description: "The current status of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RegistryProperties.Status.DisplayStatus"),
			},
			{
				Name:        "status_message",
				Description: "The detailed message for the status, including alerts and error messages.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RegistryProperties.Status.Message"),
			},
			{
				Name:        "status_timestamp",
				Description: "The timestamp when the status was changed to the current value.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("RegistryProperties.Status.Timestamp").Transform(convertDateToTime),
			},
			{
				Name:        "storage_account_id",
				Description: "The resource ID of the storage account. Only applicable to Classic SKU.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RegistryProperties.StorageAccount.ID"),
			},
			{
				Name:        "zone_redundancy",
				Description: "Indicates whether or not zone redundancy is enabled for this container registry. Valid values are: 'Enabled', 'Disabled'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("RegistryProperties.ZoneRedundancy").Transform(transform.ToString),
			},
			{
				Name:        "data_endpoint_host_names",
				Description: "A list of host names that will serve data when dataEndpointEnabled is true.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RegistryProperties.DataEndpointHostNames"),
			},
			{
				Name:        "encryption",
				Description: "The encryption settings of container registry.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RegistryProperties.Encryption"),
			},
			{
				Name:        "identity",
				Description: "The identity of the container registry.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "login_credentials",
				Description: "The login credentials for the specified container registry.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listContainerRegistryLoginCredentials,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "network_rule_set",
				Description: "The network rule set for a container registry.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RegistryProperties.NetworkRuleSet"),
			},
			{
				Name:        "policies",
				Description: "The policies for a container registry.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RegistryProperties.Policies"),
			},
			{
				Name:        "private_endpoint_connections",
				Description: "A list of private endpoint connections for a container registry.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("RegistryProperties.PrivateEndpointConnections"),
			},
			{
				Name:        "system_data",
				Description: "Metadata pertaining to creation and last modification of the resource.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "usages",
				Description: "Specifies the quota usages for the specified container registry.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listContainerRegistryUsages,
				Transform:   transform.FromValue(),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ID").Transform(idToAkas),
			},

			// Azure standard columns
			{
				Name:        "region",
				Description: ColumnDescriptionRegion,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Location").Transform(toLower),
			},
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

func listContainerRegistries(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listContainerRegistries")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := containerregistry.NewRegistriesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, registry := range result.Values() {
		d.StreamListItem(ctx, registry)
		// This will return zero if context has been cancelled (i.e due to manual cancellation) or
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

		for _, registry := range result.Values() {
			d.StreamListItem(ctx, registry)
			// This will return zero if context has been cancelled (i.e due to manual cancellation) or
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getContainerRegistry(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getContainerRegistry")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Return nil, if no input provided
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := containerregistry.NewRegistriesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}

func listContainerRegistryLoginCredentials(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listContainerRegistryLoginCredentials")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := containerregistry.NewRegistriesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	data := h.Item.(containerregistry.Registry)
	resourceGroup := strings.Split(*data.ID, "/")[4]

	op, err := client.ListCredentials(ctx, resourceGroup, *data.Name)
	if err != nil {
		if strings.Contains(err.Error(), "UnAuthorizedForCredentialOperations") {
			return nil, nil
		}
		return nil, err
	}

	return op, nil
}

func listContainerRegistryUsages(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listContainerRegistryUsages")

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := containerregistry.NewRegistriesClient(subscriptionID)
	client.Authorizer = session.Authorizer

	data := h.Item.(containerregistry.Registry)
	resourceGroup := strings.Split(*data.ID, "/")[4]

	op, err := client.ListUsages(ctx, resourceGroup, *data.Name)
	if err != nil {
		return nil, err
	}

	return op, nil
}
