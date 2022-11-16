package azure

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	sql "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
)

//// TABLE DEFINITION

func tableAzureMSSQLManagedInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_mssql_managed_instance",
		Description: "Azure Microsoft SQL Managed Instance",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getMSSQLManagedInstance,
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404", "InvalidApiVersionParameter"}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listMSSQLManagedInstances,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The friendly name that identifies the managed instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "id",
				Description: "Contains ID to identify a managed instance uniquely.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The resource type of the managed instance.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The state of the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.State"),
			},
			{
				Name:        "administrator_login",
				Description: "Administrator username for the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.AdministratorLogin"),
			},
			{
				Name:        "administrator_login_password",
				Description: "Administrator password for the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.AdministratorLoginPassword"),
			},
			{
				Name:        "collation",
				Description: "Collation of the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.Collation"),
			},
			{
				Name:        "dns_zone",
				Description: "The Dns zone that the managed instance is in.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.DNSZone"),
			},
			{
				Name:        "dns_zone_partner",
				Description: "The resource id of another managed instance whose DNS zone this managed instance will share after creation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.DNSZonePartner"),
			},
			{
				Name:        "fully_qualified_domain_name",
				Description: "The fully qualified domain name of the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.FullyQualifiedDomainName"),
			},
			{
				Name:        "instance_pool_id",
				Description: "The Id of the instance pool this managed server belongs to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.InstancePoolID"),
			},
			{
				Name:        "license_type",
				Description: "The license type of the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.LicenseType"),
			},
			{
				Name:        "maintenance_configuration_id",
				Description: "Specifies maintenance configuration id to apply to this managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.MaintenanceConfigurationID"),
			},
			{
				Name:        "managed_instance_create_mode",
				Description: "Specifies the mode of database creation.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.ManagedInstanceCreateMode"),
			},
			{
				Name:        "minimal_tls_version",
				Description: "Minimal TLS version of the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.MinimalTLSVersion"),
			},
			{
				Name:        "proxy_override",
				Description: "Connection type used for connecting to the instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.ProxyOverride"),
			},
			{
				Name:        "public_data_endpoint_enabled",
				Description: "Whether or not the public data endpoint is enabled.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ManagedInstanceProperties.PublicDataEndpointEnabled"),
			},
			{
				Name:        "restore_point_in_time",
				Description: "Specifies the point in time of the source database that will be restored to create the new database.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("ManagedInstanceProperties.RestorePointInTime").Transform(convertDateToTime),
			},
			{
				Name:        "source_managed_instance_id",
				Description: "The resource identifier of the source managed instance associated with create operation of this instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.SourceManagedInstanceID"),
			},
			{
				Name:        "storage_size_in_gb",
				Description: "The managed instance storage size in GB.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ManagedInstanceProperties.StorageSizeInGB"),
			},
			{
				Name:        "subnet_id",
				Description: "Subnet resource ID for the managed instance.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.SubnetID"),
			},
			{
				Name:        "timezone_id",
				Description: "Id of the timezone.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedInstanceProperties.TimezoneID"),
			},
			{
				Name:        "v_cores",
				Description: "The number of vcores of the managed instance.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ManagedInstanceProperties.VCores"),
			},
			{
				Name:        "encryption_protectors",
				Description: "The managed instance encryption protectors.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listMSSQLManagedInstanceEncryptionProtectors,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "identity",
				Description: "The azure active directory identity of the managed instance.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "security_alert_policies",
				Description: "The security alert policies of the managed instance.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listMSSQLManagedInstanceSecurityAlertPolicies,
				Transform:   transform.FromValue(),
			},
			{
				Name:        "sku",
				Description: "Managed instance SKU.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "vulnerability_assessments",
				Description: "The managed instance vulnerability assessments.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     listMSSQLManagedInstanceVulnerabilityAssessments,
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
		}),
	}
}

//// LIST FUNCTION

func listMSSQLManagedInstances(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		plugin.Logger(ctx).Error("azure_mssql_managed_instance.listMSSQLManagedInstances", "session error", err)
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sql.NewManagedInstancesClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mssql_managed_instance.listMSSQLManagedInstances", "client error", err)
	}

	pager := client.NewListPager(nil)
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_mssql_managed_instance.listMSSQLManagedInstances", "api error", err)
		}
		for _, managedInstance := range nextResult.Value {
			d.StreamListItem(ctx, managedInstance)
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

func getMSSQLManagedInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mssql_managed_instance.getMSSQLManagedInstance", "credential error", err)
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sql.NewManagedInstancesClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mssql_managed_instance.getMSSQLManagedInstance", "client error", err)
	}

	op, err := client.Get(ctx, resourceGroup, name, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mssql_managed_instance.getMSSQLManagedInstance", "api error", err)
		return nil, err
	}

	return op, nil
}

func listMSSQLManagedInstanceEncryptionProtectors(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	managedInstance := h.Item.(*sql.ManagedInstance)
	resourceGroup := strings.Split(string(*managedInstance.ID), "/")[4]
	managedInstanceName := *managedInstance.Name

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mssql_managed_instance.listMSSQLManagedInstanceEncryptionProtectors", "credential error", err)
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sql.NewManagedInstanceEncryptionProtectorsClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mssql_managed_instance.listMSSQLManagedInstanceEncryptionProtectors", "client error", err)
	}

	op := client.NewListByInstancePager(resourceGroup, managedInstanceName, nil)
	if err != nil {
		return nil, err
	}

	var managedInstanceEncryptionProtectors []map[string]interface{}
	for op.More() {
		nextResult, err := op.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_mssql_managed_instance.listMSSQLManagedInstanceEncryptionProtectors", "api page error", err)
		}

		for _, i := range nextResult.Value {
			managedInstanceEncryptionProtectors = append(managedInstanceEncryptionProtectors, extractMSSQLManagedInstanceEncryptionProtector(i))
		}
	}

	return managedInstanceEncryptionProtectors, nil
}

func listMSSQLManagedInstanceVulnerabilityAssessments(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	managedInstance := h.Item.(*sql.ManagedInstance)
	resourceGroup := strings.Split(string(*managedInstance.ID), "/")[4]
	managedInstanceName := *managedInstance.Name

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mssql_managed_instance.listMSSQLManagedInstanceVulnerabilityAssessments", "credential error", err)
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sql.NewManagedInstanceVulnerabilityAssessmentsClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mssql_managed_instance.listMSSQLManagedInstanceVulnerabilityAssessments", "client error", err)
	}

	op := client.NewListByInstancePager(resourceGroup, managedInstanceName, nil)
	if err != nil {
		return nil, err
	}

	var managedInstanceVulnerabilityAssessments []map[string]interface{}

	for op.More() {
		nextResult, err := op.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_mssql_managed_instance.listMSSQLManagedInstanceEncryptionProtectors", "api page error", err)
		}

		for _, i := range nextResult.Value {
			managedInstanceVulnerabilityAssessments = append(managedInstanceVulnerabilityAssessments, extractMSSQLManagedInstanceVulnerabilityAssessment(i))
		}
	}

	return managedInstanceVulnerabilityAssessments, nil
}

func listMSSQLManagedInstanceSecurityAlertPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	managedInstance := h.Item.(*sql.ManagedInstance)
	resourceGroup := strings.Split(string(*managedInstance.ID), "/")[4]
	managedInstanceName := *managedInstance.Name

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mssql_managed_instance.listMSSQLManagedInstanceSecurityAlertPolicies", "credential error", err)
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client, err := sql.NewManagedServerSecurityAlertPoliciesClient(subscriptionID, cred, nil)
	if err != nil {
		plugin.Logger(ctx).Error("azure_mssql_managed_instance.listMSSQLManagedInstanceSecurityAlertPolicies", "client error", err)
	}

	op := client.NewListByInstancePager(resourceGroup, managedInstanceName, nil)

	var managedInstanceSecurityAlertPolicies []map[string]interface{}

	for op.More() {
		nextResult, err := op.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("azure_mssql_managed_instance.listMSSQLManagedInstanceSecurityAlertPolicies", "api page error", err)
		}

		for _, i := range nextResult.Value {
			managedInstanceSecurityAlertPolicies = append(managedInstanceSecurityAlertPolicies, extractMSSQLManagedInstanceSecurityAlertPolicy(i))
		}
	}

	return managedInstanceSecurityAlertPolicies, nil
}

// If we return the API response directly, the output will not provide
// all the properties of SecurityAlertPolicies
func extractMSSQLManagedInstanceSecurityAlertPolicy(i *sql.ManagedServerSecurityAlertPolicy) map[string]interface{} {
	managedInstanceSecurityAlertPolicy := make(map[string]interface{})
	if i.ID != nil {
		managedInstanceSecurityAlertPolicy["id"] = *i.ID
	}
	if i.Name != nil {
		managedInstanceSecurityAlertPolicy["name"] = *i.Name
	}
	if i.Type != nil {
		managedInstanceSecurityAlertPolicy["type"] = *i.Type
	}
	if i.SystemData != nil {
		managedInstanceSecurityAlertPolicy["systemData"] = i.SystemData
	}
	if i.Properties != nil {
		if len(*i.Properties.State) > 0 {
			managedInstanceSecurityAlertPolicy["state"] = *i.Properties.State
		}
		if i.Properties.DisabledAlerts != nil {
			managedInstanceSecurityAlertPolicy["disabledAlerts"] = i.Properties.DisabledAlerts
		}
		if i.Properties.EmailAddresses != nil {
			managedInstanceSecurityAlertPolicy["emailAddresses"] = i.Properties.EmailAddresses
		}
		if i.Properties.EmailAccountAdmins != nil {
			managedInstanceSecurityAlertPolicy["emailAccountAdmins"] = i.Properties.EmailAccountAdmins
		}
		if i.Properties.StorageEndpoint != nil {
			managedInstanceSecurityAlertPolicy["storageEndpoint"] = i.Properties.StorageEndpoint
		}
		if i.Properties.StorageAccountAccessKey != nil {
			managedInstanceSecurityAlertPolicy["storageAccountAccessKey"] = i.Properties.StorageAccountAccessKey
		}
		if i.Properties.RetentionDays != nil {
			managedInstanceSecurityAlertPolicy["retentionDays"] = i.Properties.RetentionDays
		}
		if i.Properties.CreationTime != nil {
			managedInstanceSecurityAlertPolicy["creationTime"] = i.Properties.CreationTime
		}
	}
	return managedInstanceSecurityAlertPolicy
}

// If we return the API response directly, the output will not provide
// all the properties of ManagedInstanceVulnerabilityAssessment
func extractMSSQLManagedInstanceVulnerabilityAssessment(i *sql.ManagedInstanceVulnerabilityAssessment) map[string]interface{} {
	managedInstanceVulnerabilityAssessment := make(map[string]interface{})
	if i.ID != nil {
		managedInstanceVulnerabilityAssessment["id"] = *i.ID
	}
	if i.Name != nil {
		managedInstanceVulnerabilityAssessment["name"] = *i.Name
	}
	if i.Type != nil {
		managedInstanceVulnerabilityAssessment["type"] = *i.Type
	}
	if i.Properties != nil {

		if i.Properties.RecurringScans != nil {
			managedInstanceVulnerabilityAssessment["recurringScans"] = i.Properties.RecurringScans
		}
		if i.Properties.StorageAccountAccessKey != nil {
			managedInstanceVulnerabilityAssessment["storageAccountAccessKey"] = *i.Properties.StorageAccountAccessKey
		}
		if i.Properties.StorageContainerPath != nil {
			managedInstanceVulnerabilityAssessment["storageContainerPath"] = *i.Properties.StorageContainerPath
		}
		if i.Properties.StorageContainerSasKey != nil {
			managedInstanceVulnerabilityAssessment["storageContainerSasKey"] = *i.Properties.StorageContainerSasKey
		}
	}
	return managedInstanceVulnerabilityAssessment
}

// If we return the API response directly, the output will not provide
// all the properties of ManagedInstanceEncryptionProtector
func extractMSSQLManagedInstanceEncryptionProtector(i *sql.ManagedInstanceEncryptionProtector) map[string]interface{} {
	managedInstanceEncryptionProtector := make(map[string]interface{})
	if i.ID != nil {
		managedInstanceEncryptionProtector["id"] = *i.ID
	}
	if i.Name != nil {
		managedInstanceEncryptionProtector["name"] = *i.Name
	}
	if i.Type != nil {
		managedInstanceEncryptionProtector["type"] = *i.Type
	}
	if i.Kind != nil {
		managedInstanceEncryptionProtector["kind"] = *i.Kind
	}
	if i.Properties != nil {

		if i.Properties.AutoRotationEnabled != nil {
			managedInstanceEncryptionProtector["autoRotationEnabled"] = i.Properties.AutoRotationEnabled
		}
		if i.Properties.ServerKeyName != nil {
			managedInstanceEncryptionProtector["serverKeyName"] = i.Properties.ServerKeyName
		}
		if len(*i.Properties.ServerKeyType) > 0 {
			managedInstanceEncryptionProtector["serverKeyType"] = i.Properties.ServerKeyType
		}
		if i.Properties.Thumbprint != nil {
			managedInstanceEncryptionProtector["thumbprint"] = i.Properties.Thumbprint
		}
		if i.Properties.URI != nil {
			managedInstanceEncryptionProtector["uri"] = i.Properties.URI
		}
	}
	return managedInstanceEncryptionProtector
}
