package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION ////

func tableAzureAPIManagement(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_api_management",
		Description: "Azure API Management Service",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getAPIManagement,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "InvalidApiVersionParameter", "ResourceGroupNotFound"}),
		},
		List: &plugin.ListConfig{
			Hydrate: listAPIManagements,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "A friendly name that identifies an api management",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify an api management uniquely",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "Type of the resource",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "provisioning_state",
				Description: "The current provisioning state of the API Management service",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.ProvisioningState"),
			},
			{
				Name:        "created_at_utc",
				Description: "Creation UTC date of the API Management service",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ServiceProperties.CreatedAtUtc").Transform(convertDateToTime),
			},
			{
				Name:        "gateway_regional_url",
				Description: "Gateway URL of the API Management service in the Default Region",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.GatewayRegionalURL"),
			},
			{
				Name:        "gateway_url",
				Description: "Gateway URL of the API Management service",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.GatewayURL"),
			},
			{
				Name:        "identity_principal_id",
				Description: "The principal id of the identity",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Identity.PrincipalID").Transform(transform.ToString),
			},
			{
				Name:        "identity_tenant_id",
				Description: "The client tenant id of the identity",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Identity.TenantID").Transform(transform.ToString),
			},
			{
				Name:        "identity_type",
				Description: "The type of identity used for the resource",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Identity.Type").Transform(transform.ToString),
			},
			{
				Name:        "management_api_url",
				Description: "Management API endpoint URL of the API Management service",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.ManagementAPIURL"),
			},
			{
				Name:        "notification_sender_email",
				Description: "Email address from which the notification will be sent",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.NotificationSenderEmail"),
			},
			{
				Name:        "portal_url",
				Description: "Publisher portal endpoint Url of the API Management service",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.PortalURL"),
			},
			{
				Name:        "publisher_email",
				Description: "Email address of the publisher of the API Management service",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.PublisherEmail"),
			},
			{
				Name:        "publisher_name",
				Description: "Name of the publisher of the API Management service",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.PublisherName"),
			},
			{
				Name:        "sku_name",
				Description: "Name of the Sku",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name").Transform(transform.ToString),
			},
			{
				Name:        "sku_capacity",
				Description: "Capacity of the SKU (number of deployed units of the SKU)",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Sku.Capacity"),
			},
			{
				Name:        "target_provisioning_state",
				Description: "The provisioning state of the API Management service, which is targeted by the long running operation started on the service",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.TargetProvisioningState"),
			},
			{
				Name:        "virtual_network_configuration_id",
				Description: "Contains the virtual network ID",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.VirtualNetworkConfiguration.Vnetid"),
			},
			{
				Name:        "virtual_network_configuration_subnet_name",
				Description: "Contains the name of the subnet",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.VirtualNetworkConfiguration.Subnetname"),
			},
			{
				Name:        "virtual_network_configuration_subnet_resource_id",
				Description: "The full resource ID of a subnet in a virtual network to deploy the API Management service in",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.VirtualNetworkConfiguration.SubnetResourceID"),
			},
			{
				Name:        "additional_locations",
				Description: "Additional datacenter locations of the API Management service",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServiceProperties.AdditionalLocations"),
			},
			{
				Name:        "host_name_configurations",
				Description: "A list of custom hostname configuration of the API Management service",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServiceProperties.HostnameConfigurations"),
			},
			{
				Name:        "private_ip_addresses",
				Description: "A list of private Static Load Balanced IP addresses of the API Management service in Primary region which is deployed in an Internal Virtual Network",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServiceProperties.PrivateIPAddresses"),
			},
			{
				Name:        "public_ip_addresses",
				Description: "A list of public Static Load Balanced IP addresses of the API Management service in Primary region",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServiceProperties.PublicIPAddresses"),
			},

			// Standard columns
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

//// FETCH FUNCTIONS ////

func listAPIManagements(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	apiManagementClient := apimanagement.NewServiceClient(subscriptionID)
	apiManagementClient.Authorizer = session.Authorizer

	result, err := apiManagementClient.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, apiManagement := range result.Values() {
		d.StreamListItem(ctx, apiManagement)
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

		for _, apiManagement := range result.Values() {
			d.StreamListItem(ctx, apiManagement)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS ////

func getAPIManagement(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAPIManagement")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// resourceGroupName can't be empty
	// Error: pq: rpc error: code = Unknown desc = apimanagement.ServiceClient#Get: Invalid input: autorest/validation: validation failed: parameter=serviceName
	// constraint=MinLength value="" details: value length must be greater than or equal to 1
	if len(name) < 1 {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	apiManagementClient := apimanagement.NewServiceClient(subscriptionID)
	apiManagementClient.Authorizer = session.Authorizer

	op, err := apiManagementClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, err
	}

	return op, nil
}
