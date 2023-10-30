package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION ////

func tableAzureAPIManagementBackend(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_api_management_backend",
		Description: "Azure API Management Backend",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"backend_id", "resource_group", "service_name"}),
			Hydrate:    getAPIManagementBackend,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound"}),
			},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listAPIManagements,
			Hydrate:       listAPIManagementBackends,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "A friendly name that identifies an API management backend.",
			},
			{
				Name:        "id",
				Description: "Contains ID to identify an API management backend uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "url",
				Description: "Runtime Url of the Backend.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BackendContractProperties.URL"),
			},
			{
				Name:        "type",
				Description: "Resource type for API Management resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "protocol",
				Description: "Backend communication protocol. Possible values include: 'BackendProtocolHTTP', 'BackendProtocolSoap'.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BackendContractProperties.Protocol"),
			},
			{
				Name:        "description",
				Description: "Backend Description.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BackendContractProperties.Description"),
			},
			{
				Name:        "resource_id",
				Description: "Backend Description.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("BackendContractProperties.ResourceID"),
			},
			{
				Name:        "properties",
				Description: "Backend Properties contract.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("BackendContractProperties.Properties"),
			},
			{
				Name:        "credentials",
				Description: "Backend Credentials Contract Properties.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("BackendContractProperties.Credentials"),
			},
			{
				Name:        "proxy",
				Description: "Backend Proxy Contract Properties.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("BackendContractProperties.Proxy"),
			},
			{
				Name:        "tls",
				Description: "Backend TLS Properties.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("BackendContractProperties.TLS"),
			},
			{
				Name:        "service_name",
				Description: "Name of the API management service.",
				Type:        proto.ColumnType_STRING,
			},
			// We have added this as an extra column because the get call takes only the last path of the id as the backend_id which we do not get from the API
			{
				Name:        "backend_id",
				Description: "The API management backend ID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(lastPathElement),
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

type BackendWithServiceName struct {
	apimanagement.BackendContract
	ServiceName string
}

//// LIST FUNCTION

func listAPIManagementBackends(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	serviceInfo := h.Item.(apimanagement.ServiceResource)
	serviceName := *serviceInfo.Name
	resourceGroup := strings.Split(*serviceInfo.ID, "/")[4]

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	apiManagementBackendClient := apimanagement.NewBackendClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	apiManagementBackendClient.Authorizer = session.Authorizer

	result, err := apiManagementBackendClient.ListByService(ctx, resourceGroup, serviceName, "", nil, nil)
	if err != nil {
		plugin.Logger(ctx).Error("listAPIManagementBackends", "list", err)
		return nil, err
	}
	for _, apiManagementBackend := range result.Values() {
		backendWithService := &BackendWithServiceName{
			apiManagementBackend,
			serviceName,
		}
		d.StreamListItem(ctx, backendWithService)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listAPIManagementBackends", "list_paging", err)
			return nil, err
		}

		for _, apiManagementBackend := range result.Values() {
			d.StreamListItem(ctx, apiManagementBackend)
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

func getAPIManagementBackend(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("getAPIManagementBackend")

	backendID := d.EqualsQualString("backend_id")
	serviceName := d.EqualsQualString("service_name")
	resourceGroup := d.EqualsQualString("resource_group")

	// resourceGroupName can't be empty
	// Error: pq: rpc error: code = Unknown desc = apimanagement.ServiceClient#Get: Invalid input: autorest/validation: validation failed: parameter=serviceName
	// constraint=MinLength value="" details: value length must be greater than or equal to 1
	if len(backendID) < 1 {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	apiManagementBackendClient := apimanagement.NewBackendClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	apiManagementBackendClient.Authorizer = session.Authorizer

	op, err := apiManagementBackendClient.Get(ctx, resourceGroup, serviceName, backendID)
	if err != nil {
		plugin.Logger(ctx).Error("getAPIManagementBackend", "get", err)
		return nil, err
	}

	return BackendWithServiceName{op, serviceName}, nil
}