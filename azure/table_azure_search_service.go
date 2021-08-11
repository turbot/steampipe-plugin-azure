package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/search/mgmt/2020-08-01/search"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAzureSearchService(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_search_service",
		Description: "Azure Search Service",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getSearchService,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listSearchServices,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the resource.",
			},
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "Fully qualified resource ID for the resource.",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "provisioning_state",
				Type:        proto.ColumnType_STRING,
				Description: "The state of the last provisioning operation performed on the search service. Provisioning is an intermediate state that occurs while service capacity is being established. After capacity is set up, provisioningState changes to either 'succeeded' or 'failed'. Client applications can poll provisioning status (the recommended polling interval is from 30 seconds to one minute) by using the Get Search Service operation to see when an operation is completed. If you are using the free service, this value tends to come back as 'succeeded' directly in the call to Create search service. This is because the free service uses capacity that is already set up. Possible values include: 'Succeeded', 'Provisioning', 'Failed'.",
				Transform:   transform.FromField("ServiceProperties.ProvisioningState"),
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "The status of the search service. Possible values include: 'running': The search service is running and no provisioning operations are underway. 'provisioning': The search service is being provisioned or scaled up or down. 'deleting': The search service is being deleted. 'degraded': The search service is degraded. This can occur when the underlying search units are not healthy. The search service is most likely operational, but performance might be slow and some requests might be dropped. 'disabled': The search service is disabled. In this state, the service will reject all API requests. 'error': The search service is in an error state. If your service is in the degraded, disabled, or error states, it means the Azure Cognitive Search team is actively investigating the underlying issue. Dedicated services in these states are still chargeable based on the number of search units provisioned. Possible values include: 'ServiceStatusRunning', 'ServiceStatusProvisioning', 'ServiceStatusDeleting', 'ServiceStatusDegraded', 'ServiceStatusDisabled', 'ServiceStatusError'",
				Transform:   transform.FromField("ServiceProperties.Status"),
			},
			{
				Name:        "status_details",
				Type:        proto.ColumnType_STRING,
				Description: "The status of the search service. Possible values include: 'runniThe details of the search service status.",
				Transform:   transform.FromField("ServiceProperties.StatusDetails"),
			},
			{
				Name:        "type",
				Type:        proto.ColumnType_STRING,
				Description: "The type of the resource. E.g. 'Microsoft.Compute/virtualMachines' or 'Microsoft.Storage/storageAccounts'.",
			},
			{
				Name:        "hosting_mode",
				Type:        proto.ColumnType_STRING,
				Description: "Applicable only for the standard3 SKU. You can set this property to enable up to 3 high density partitions that allow up to 1000 indexes, which is much higher than the maximum indexes allowed for any other SKU. For the standard3 SKU, the value is either 'default' or 'highDensity'. For all other SKUs, this value must be 'default'. Possible values include: 'Default', 'HighDensity'.",
				Transform:   transform.FromField("ServiceProperties.HostingMode"),
			},
			{
				Name:        "partition_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of partitions in the search service; if specified, it can be 1, 2, 3, 4, 6, or 12. Values greater than 1 are only valid for standard SKUs. For 'standard3' services with hostingMode set to 'highDensity', the allowed values are between 1 and 3.",
				Transform:   transform.FromField("ServiceProperties.PartitionCount"),
			},
			{
				Name:        "public_network_access",
				Type:        proto.ColumnType_STRING,
				Description: "This value can be set to 'enabled' to avoid breaking changes on existing customer resources and templates. If set to 'disabled', traffic over public interface is not allowed, and private endpoint connections would be the exclusive access method. Possible values include: 'Enabled', 'Disabled'.",
				Transform:   transform.FromField("ServiceProperties.PublicNetworkAccess"),
			},
			{
				Name:        "replica_count",
				Type:        proto.ColumnType_INT,
				Description: "The number of replicas in the search service. If specified, it must be a value between 1 and 12 inclusive for standard SKUs or between 1 and 3 inclusive for basic SKU.",
				Transform:   transform.FromField("ServiceProperties.ReplicaCount"),
			},
			{
				Name:        "sku_name",
				Type:        proto.ColumnType_STRING,
				Description: "The SKU of the Search Service, which determines price tier and capacity limits. This property is required when creating a new Search Service.",
				Transform:   transform.FromField("Sku.Name"),
			},
			{
				Name:        "identity",
				Type:        proto.ColumnType_JSON,
				Description: "The identity of the resource.",
				Transform:   transform.FromField("Identity"),
			},		
			{
				Name:        "network_rule_set",
				Type:        proto.ColumnType_JSON,
				Description: "Network specific rules that determine how the Azure Cognitive Search service may be reached.",
				Transform:   transform.FromField("ServiceProperties.NetworkRuleSet"),
			},
			{
				Name:        "private_endpoint_connections",
				Type:        proto.ColumnType_JSON,
				Description: "The list of private endpoint connections to the Azure Cognitive Search service.",
				Transform:   transform.FromField("ServiceProperties.PrivateEndpointConnections"),
			},
			{
				Name:        "shared_private_link_resources",
				Type:        proto.ColumnType_JSON,
				Description: "The list of shared private link resources managed by the Azure Cognitive Search service.",
				Transform:   transform.FromField("ServiceProperties.SharedPrivateLinkResources"),
			},
			{
				Name:        "tags_src",
				Type:        proto.ColumnType_JSON,
				Description: "The resource tags.",
				Transform:   transform.FromField("Tags"),
			},

			// Azure standard columns
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
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromID).Transform(toLower),
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

func listSearchServices(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, nil
	}
	subscriptionID := session.SubscriptionID

	searchClient := search.NewServicesClient(subscriptionID)
	searchClient.Authorizer = session.Authorizer

	result, err := searchClient.ListBySubscription(ctx, nil)
	if err != nil {
		return nil, err
	}
	for _, service := range result.Values() {
		d.StreamListItem(ctx, service)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
		for _, service := range result.Values() {
			d.StreamListItem(ctx, service)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getSearchService(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getSearchService")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	searchClient := search.NewServicesClient(subscriptionID)
	searchClient.Authorizer = session.Authorizer

	op, err := searchClient.Get(ctx, resourceGroup, name, nil)
	if err != nil {
		return nil, err
	}

	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}
