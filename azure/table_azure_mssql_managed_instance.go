package azure

import (
	"context"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v5.0/sql"
)

//// TABLE DEFINITION

func tableAzureMSSQLManagedInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_mssql_managed_instance",
		Description: "Azure Microsoft SQL Managed Instance",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:           getMSSQLManagedInstance,
			ShouldIgnoreError: isNotFoundError([]string{"ResourceNotFound", "ResourceGroupNotFound", "404", "InvalidApiVersionParameter"}),
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
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID
	client := sql.NewManagedInstancesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(ctx, "")
	if err != nil {
		plugin.Logger(ctx).Error("listMSSQLManagedInstances", "list", err)
		return nil, err
	}
	for _, managedInstance := range result.Values() {
		d.StreamListItem(ctx, managedInstance)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listMSSQLManagedInstances", "list_paging", err)
			return nil, err
		}
		for _, managedInstance := range result.Values() {
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
	plugin.Logger(ctx).Trace("getMSSQLManagedInstance")

	name := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroup := d.KeyColumnQuals["resource_group"].GetStringValue()

	// Return nil, of no input provided
	if name == "" || resourceGroup == "" {
		return nil, nil
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := sql.NewManagedInstancesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		plugin.Logger(ctx).Error("getMSSQLManagedInstance", "get", err)
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}

func listMSSQLManagedInstanceEncryptionProtectors(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listMSSQLManagedInstanceEncryptionProtectors")

	managedInstance := h.Item.(sql.ManagedInstance)
	resourceGroup := strings.Split(string(*managedInstance.ID), "/")[4]
	managedInstanceName := *managedInstance.Name

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := sql.NewManagedInstanceEncryptionProtectorsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.ListByInstance(ctx, resourceGroup, managedInstanceName)
	if err != nil {
		plugin.Logger(ctx).Error("listMSSQLManagedInstanceEncryptionProtectors", "list", err)
		return nil, err
	}

	var managedInstanceEncryptionProtectors []map[string]interface{}

	for _, i := range op.Values() {
		managedInstanceEncryptionProtectors = append(managedInstanceEncryptionProtectors, extractMSSQLManagedInstanceEncryptionProtector(i))
	}

	for op.NotDone() {
		err = op.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listMSSQLManagedInstanceEncryptionProtectors", "list_paging", err)
			return nil, err
		}
		for _, i := range op.Values() {
			managedInstanceEncryptionProtectors = append(managedInstanceEncryptionProtectors, extractMSSQLManagedInstanceEncryptionProtector(i))
		}
	}

	return managedInstanceEncryptionProtectors, nil
}

func listMSSQLManagedInstanceVulnerabilityAssessments(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listMSSQLManagedInstanceVulnerabilityAssessments")

	managedInstance := h.Item.(sql.ManagedInstance)
	resourceGroup := strings.Split(string(*managedInstance.ID), "/")[4]
	managedInstanceName := *managedInstance.Name

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := sql.NewManagedInstanceVulnerabilityAssessmentsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.ListByInstance(ctx, resourceGroup, managedInstanceName)
	if err != nil {
		plugin.Logger(ctx).Error("listMSSQLManagedInstanceVulnerabilityAssessments", "list", err)
		return nil, err
	}

	var managedInstanceVulnerabilityAssessments []map[string]interface{}

	for _, i := range op.Values() {
		managedInstanceVulnerabilityAssessments = append(managedInstanceVulnerabilityAssessments, extractMSSQLManagedInstanceVulnerabilityAssessment(i))
	}

	for op.NotDone() {
		err = op.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listMSSQLManagedInstanceVulnerabilityAssessments", "list_paging", err)
			return nil, err
		}
		for _, i := range op.Values() {
			managedInstanceVulnerabilityAssessments = append(managedInstanceVulnerabilityAssessments, extractMSSQLManagedInstanceVulnerabilityAssessment(i))
		}
	}

	return managedInstanceVulnerabilityAssessments, nil
}

func listMSSQLManagedInstanceSecurityAlertPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listMSSQLManagedInstanceSecurityAlertPolicies")

	managedInstance := h.Item.(sql.ManagedInstance)
	resourceGroup := strings.Split(string(*managedInstance.ID), "/")[4]
	managedInstanceName := *managedInstance.Name

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := sql.NewManagedServerSecurityAlertPoliciesClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	client.Authorizer = session.Authorizer

	op, err := client.ListByInstance(ctx, resourceGroup, managedInstanceName)
	if err != nil {
		plugin.Logger(ctx).Error("listMSSQLManagedInstanceSecurityAlertPolicies", "list", err)
		return nil, err
	}

	var managedInstanceSecurityAlertPolicies []map[string]interface{}

	for _, i := range op.Values() {
		managedInstanceSecurityAlertPolicies = append(managedInstanceSecurityAlertPolicies, extractMSSQLManagedInstanceSecurityAlertPolicy(i))
	}

	for op.NotDone() {
		err = op.NextWithContext(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("listMSSQLManagedInstanceSecurityAlertPolicies", "list_paging", err)
			return nil, err
		}
		for _, i := range op.Values() {
			managedInstanceSecurityAlertPolicies = append(managedInstanceSecurityAlertPolicies, extractMSSQLManagedInstanceSecurityAlertPolicy(i))
		}
	}

	return managedInstanceSecurityAlertPolicies, nil
}

// If we return the API response directly, the output will not provide
// all the properties of SecurityAlertPolicies
func extractMSSQLManagedInstanceSecurityAlertPolicy(i sql.ManagedServerSecurityAlertPolicy) map[string]interface{} {
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
	if i.SecurityAlertsPolicyProperties != nil {
		if len(i.SecurityAlertsPolicyProperties.State) > 0 {
			managedInstanceSecurityAlertPolicy["state"] = i.SecurityAlertsPolicyProperties.State
		}
		if i.SecurityAlertsPolicyProperties.DisabledAlerts != nil {
			managedInstanceSecurityAlertPolicy["disabledAlerts"] = i.SecurityAlertsPolicyProperties.DisabledAlerts
		}
		if i.SecurityAlertsPolicyProperties.EmailAddresses != nil {
			managedInstanceSecurityAlertPolicy["emailAddresses"] = i.SecurityAlertsPolicyProperties.EmailAddresses
		}
		if i.SecurityAlertsPolicyProperties.EmailAccountAdmins != nil {
			managedInstanceSecurityAlertPolicy["emailAccountAdmins"] = i.SecurityAlertsPolicyProperties.EmailAccountAdmins
		}
		if i.SecurityAlertsPolicyProperties.StorageEndpoint != nil {
			managedInstanceSecurityAlertPolicy["storageEndpoint"] = i.SecurityAlertsPolicyProperties.StorageEndpoint
		}
		if i.SecurityAlertsPolicyProperties.StorageAccountAccessKey != nil {
			managedInstanceSecurityAlertPolicy["storageAccountAccessKey"] = i.SecurityAlertsPolicyProperties.StorageAccountAccessKey
		}
		if i.SecurityAlertsPolicyProperties.RetentionDays != nil {
			managedInstanceSecurityAlertPolicy["retentionDays"] = i.SecurityAlertsPolicyProperties.RetentionDays
		}
		if i.SecurityAlertsPolicyProperties.CreationTime != nil {
			managedInstanceSecurityAlertPolicy["creationTime"] = i.SecurityAlertsPolicyProperties.CreationTime
		}
	}
	return managedInstanceSecurityAlertPolicy
}

// If we return the API response directly, the output will not provide
// all the properties of ManagedInstanceVulnerabilityAssessment
func extractMSSQLManagedInstanceVulnerabilityAssessment(i sql.ManagedInstanceVulnerabilityAssessment) map[string]interface{} {
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
	if i.ManagedInstanceVulnerabilityAssessmentProperties.RecurringScans != nil {
		managedInstanceVulnerabilityAssessment["recurringScans"] = i.ManagedInstanceVulnerabilityAssessmentProperties.RecurringScans
	}
	if i.ManagedInstanceVulnerabilityAssessmentProperties.StorageAccountAccessKey != nil {
		managedInstanceVulnerabilityAssessment["storageAccountAccessKey"] = *i.ManagedInstanceVulnerabilityAssessmentProperties.StorageAccountAccessKey
	}
	if i.ManagedInstanceVulnerabilityAssessmentProperties.StorageContainerPath != nil {
		managedInstanceVulnerabilityAssessment["storageContainerPath"] = *i.ManagedInstanceVulnerabilityAssessmentProperties.StorageContainerPath
	}
	if i.ManagedInstanceVulnerabilityAssessmentProperties.StorageContainerSasKey != nil {
		managedInstanceVulnerabilityAssessment["storageContainerSasKey"] = *i.ManagedInstanceVulnerabilityAssessmentProperties.StorageContainerSasKey
	}
	return managedInstanceVulnerabilityAssessment
}

// If we return the API response directly, the output will not provide
// all the properties of ManagedInstanceEncryptionProtector
func extractMSSQLManagedInstanceEncryptionProtector(i sql.ManagedInstanceEncryptionProtector) map[string]interface{} {
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
	if i.ManagedInstanceEncryptionProtectorProperties.AutoRotationEnabled != nil {
		managedInstanceEncryptionProtector["autoRotationEnabled"] = i.ManagedInstanceEncryptionProtectorProperties.AutoRotationEnabled
	}
	if i.ManagedInstanceEncryptionProtectorProperties.ServerKeyName != nil {
		managedInstanceEncryptionProtector["serverKeyName"] = i.ManagedInstanceEncryptionProtectorProperties.ServerKeyName
	}
	if len(i.ManagedInstanceEncryptionProtectorProperties.ServerKeyType) > 0 {
		managedInstanceEncryptionProtector["serverKeyType"] = i.ManagedInstanceEncryptionProtectorProperties.ServerKeyType
	}
	if i.ManagedInstanceEncryptionProtectorProperties.Thumbprint != nil {
		managedInstanceEncryptionProtector["thumbprint"] = i.ManagedInstanceEncryptionProtectorProperties.Thumbprint
	}
	if i.ManagedInstanceEncryptionProtectorProperties.URI != nil {
		managedInstanceEncryptionProtector["uri"] = i.ManagedInstanceEncryptionProtectorProperties.URI
	}
	return managedInstanceEncryptionProtector
}
