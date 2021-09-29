package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/2020-09-01/monitor/mgmt/insights"
	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
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
				Description: "A friendly name that identifies an API management service.",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify an API management service uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "provisioning_state",
				Description: "The current provisioning state of the API management service. Possible values include: 'Created', 'Activating', 'Succeeded', 'Updating', 'Failed', 'Stopped', 'Terminating', 'TerminationFailed', 'Deleted'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.ProvisioningState"),
			},
			{
				Name:        "type",
				Description: "Type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "created_at_utc",
				Description: "Creation UTC date of the API management service.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ServiceProperties.CreatedAtUtc").Transform(convertDateToTime),
			},
			{
				Name:        "developer_portal_url",
				Description: "Developer Portal endpoint URL of the API management service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.DeveloperPortalURL"),
			},
			{
				Name:        "disable_gateway",
				Description: "Property only valid for an API management service deployed in multiple locations. This can be used to disable the gateway in master region.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ServiceProperties.DisableGateway"),
				Default:     false,
			},
			{
				Name:        "enable_client_certificate",
				Description: "Property only meant to be used for Consumption SKU Service. This enforces a client certificate to be presented on each request to the gateway. This also enables the ability to authenticate the certificate in the policy on the gateway.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ServiceProperties.EnableClientCertificate"),
				Default:     false,
			},
			{
				Name:        "etag",
				Description: "An unique read-only string that changes whenever the resource is updated.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "gateway_regional_url",
				Description: "Gateway URL of the API management service in the default region.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.GatewayRegionalURL"),
			},
			{
				Name:        "gateway_url",
				Description: "Gateway URL of the API management service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.GatewayURL"),
			},
			{
				Name:        "identity_principal_id",
				Description: "The principal id of the identity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Identity.PrincipalID").Transform(transform.ToString),
			},
			{
				Name:        "identity_tenant_id",
				Description: "The client tenant id of the identity.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Identity.TenantID").Transform(transform.ToString),
			},
			{
				Name:        "identity_type",
				Description: "The type of identity used for the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Identity.Type").Transform(transform.ToString),
			},
			{
				Name:        "management_api_url",
				Description: "Management API endpoint URL of the API management service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.ManagementAPIURL"),
			},
			{
				Name:        "notification_sender_email",
				Description: "Email address from which the notification will be sent.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.NotificationSenderEmail"),
			},
			{
				Name:        "portal_url",
				Description: "Publisher portal endpoint URL of the API management service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.PortalURL"),
			},
			{
				Name:        "publisher_email",
				Description: "Publisher email of the API management service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.PublisherEmail"),
			},
			{
				Name:        "publisher_name",
				Description: "Publisher name of the API management service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.PublisherName"),
			},
			{
				Name:        "restore",
				Description: "Undelete API management service if it was previously soft-deleted.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ServiceProperties.Restore"),
				Default:     false,
			},
			{
				Name:        "scm_url",
				Description: "SCM endpoint URL of the API management service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.ScmURL"),
			},
			{
				Name:        "sku_capacity",
				Description: "Capacity of the SKU (number of deployed units of the SKU)",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Sku.Capacity"),
			},
			{
				Name:        "sku_name",
				Description: "Name of the Sku",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Sku.Name").Transform(transform.ToString),
			},
			{
				Name:        "target_provisioning_state",
				Description: "The provisioning state of the API management service, which is targeted by the long running operation started on the service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.TargetProvisioningState"),
			},
			{
				Name:        "virtual_network_configuration_subnet_name",
				Description: "The name of the subnet.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.VirtualNetworkConfiguration.Subnetname"),
			},
			{
				Name:        "virtual_network_configuration_subnet_resource_id",
				Description: "The full resource ID of a subnet in a virtual network to deploy the API Management service in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.VirtualNetworkConfiguration.SubnetResourceID"),
			},
			{
				Name:        "virtual_network_configuration_id",
				Description: "The virtual network ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.VirtualNetworkConfiguration.Vnetid"),
			},
			{
				Name:        "virtual_network_type",
				Description: "The type of VPN in which API management service needs to be configured in. None (Default Value) means the API management service is not part of any Virtual Network, External means the API management deployment is set up inside a Virtual Network having an Internet Facing Endpoint, and Internal means that API management deployment is setup inside a Virtual Network having an Intranet Facing Endpoint only. Possible values include: 'None', 'External', 'Internal'",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServiceProperties.VirtualNetworkType"),
			},
			{
				Name:        "additional_locations",
				Description: "Additional datacenter locations of the API management service.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServiceProperties.AdditionalLocations"),
			},
			{
				Name:        "api_version_constraint",
				Description: "Control plane APIs version constraint for the API management service.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServiceProperties.APIVersionConstraint"),
			},
			{
				Name:        "certificates",
				Description: "List of certificates that need to be installed in the API management service.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServiceProperties.Certificates"),
			},
			{
				Name:        "custom_properties",
				Description: "Custom properties of the API management service.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServiceProperties.CustomProperties"),
			},
			{
				Name:        "diagnostic_settings",
				Description: "A list of active diagnostic settings for the API management service.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listAPIManagementDiagnosticSettings,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "host_name_configurations",
				Description: "Custom hostname configuration of the API management service.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServiceProperties.HostnameConfigurations"),
			},
			{
				Name:        "identity_user_assigned_identities",
				Description: "The list of user identities associated with the resource.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Identity.UserAssignedIdentities"),
			},
			{
				Name:        "private_ip_addresses",
				Description: "Private static load balanced IP addresses of the API management service in primary region which is deployed in an internal virtual network. Available only for 'Basic', 'Standard', 'Premium' and 'Isolated' SKU.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServiceProperties.PrivateIPAddresses"),
			},
			{
				Name:        "public_ip_addresses",
				Description: "Public static load balanced IP addresses of the API management service in primary region. Available only for 'Basic', 'Standard', 'Premium' and 'Isolated' SKU.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ServiceProperties.PublicIPAddresses"),
			},
			{
				Name:        "zones",
				Description: "A list of availability zones denoting where the resource needs to come from.",
				Type:        proto.ColumnType_JSON,
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
				Transform:   transform.FromField("Location").Transform(formatRegion).Transform(toLower),
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
		plugin.Logger(ctx).Error("listAPIManagements", "list", err)
		return nil, err
	}
	for _, apiManagement := range result.Values() {
		d.StreamListItem(ctx, apiManagement)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listAPIManagements", "list_paging", err)
			return nil, err
		}

		for _, apiManagement := range result.Values() {
			d.StreamListItem(ctx, apiManagement)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

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
		plugin.Logger(ctx).Error("getAPIManagement", "get", err)
		return nil, err
	}

	return op, nil
}

func listAPIManagementDiagnosticSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listAPIManagementDiagnosticSettings")
	id := *h.Item.(apimanagement.ServiceResource).ID

	// Create session
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := insights.NewDiagnosticSettingsClient(subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.List(ctx, id)
	if err != nil {
		plugin.Logger(ctx).Error("listAPIManagementDiagnosticSettings", "list", err)
		return nil, err
	}

	// If we return the API response directly, the output does not provide
	// all the contents of DiagnosticSettings
	var diagnosticSettings []map[string]interface{}
	for _, i := range *op.Value {
		objectMap := make(map[string]interface{})
		if i.ID != nil {
			objectMap["id"] = i.ID
		}
		if i.Name != nil {
			objectMap["name"] = i.Name
		}
		if i.Type != nil {
			objectMap["type"] = i.Type
		}
		if i.DiagnosticSettings != nil {
			objectMap["properties"] = i.DiagnosticSettings
		}
		diagnosticSettings = append(diagnosticSettings, objectMap)
	}
	return diagnosticSettings, nil
}
