package azure

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v1.0/security"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

//// TABLE DEFINITION

func tableAzureSecurityCenterSubAssessment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_security_center_sub_assessment",
		Description: "Azure Security Center Sub Assessment",
		List: &plugin.ListConfig{
			Hydrate: listSecurityCenterSubAssessments,
		},
		Columns: azureColumns([]*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The resource id.",
				Transform:   transform.FromGo(),
			},
			{
				Name:        "name",
				Description: "The resource name.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The resource type.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "assessment_name",
				Description: "The assessment name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(extractAssessmentName),
			},
			{
				Name:        "category",
				Description: "Category of the sub-assessment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SubAssessmentProperties.Category"),
			},
			{
				Name:        "description",
				Description: "Human readable description of the assessment status.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SubAssessmentProperties.Description"),
			},
			{
				Name:        "display_name",
				Description: "User friendly display name of the sub-assessment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SubAssessmentProperties.DisplayName"),
			},
			{
				Name:        "impact",
				Description: "Description of the impact of this sub-assessment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SubAssessmentProperties.Impact"),
			},
			{
				Name:        "remediation",
				Description: "Information on how to remediate this sub-assessment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SubAssessmentProperties.Remediation"),
			},
			{
				Name:        "time_generated",
				Description: "The date and time the sub-assessment was generated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SubAssessmentProperties.TimeGenerated"),
			},
			{
				Name:        "assessed_resource_type",
				Description: "Details of the sub-assessment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(extractAssessedResourceType),
			},
			{
				Name:        "container_registry_vulnerability_properties",
				Description: "ContainerRegistryVulnerabilityProperties details of the resource that was assessed.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractContainerRegistryVulnerabilityProperties),
			},
			{
				Name:        "resource_details",
				Description: "Details of the resource that was assessed.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractResourceDetails),
			},
			{
				Name:        "status",
				Description: "The status of the sub-assessment.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractSubAssessmentStatus),
			},
			{
				Name:        "server_vulnerability_properties",
				Description: "ServerVulnerabilityProperties details of the resource that was assessed.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractServerVulnerabilityProperties),
			},
			{
				Name:        "sql_server_vulnerability_properties",
				Description: "SQLServerVulnerabilityProperties details of the resource that was assessed.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(extractSQLServerVulnerabilityProperties),
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
				Transform:   transform.FromField("ID").Transform(extractResourceGroupFromSubAssessmentID),
			},
		}),
	}
}

//// LIST FUNCTION

func listSecurityCenterSubAssessments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Info("azure_security_center_sub_assessment.listSecurityCenterSubAssessments")
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		logger.Error("azure_security_center_sub_assessment.listSecurityCenterSubAssessments", "connection_error", err)
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	subAssessmentClient := security.NewSubAssessmentsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID, "")
	subAssessmentClient.Authorizer = session.Authorizer

	result, err := subAssessmentClient.ListAll(ctx, "subscriptions/"+subscriptionID)
	if err != nil {
		logger.Error("azure_security_center_sub_assessment.listSecurityCenterSubAssessments", "query_error", err)
		return err, nil
	}

	for _, subAssessments := range result.Values() {
		d.StreamListItem(ctx, subAssessments)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			logger.Error("azure_security_center_sub_assessment.listSecurityCenterSubAssessments", "query_error", err)
			return err, nil
		}
		for _, subAssessments := range result.Values() {
			d.StreamListItem(ctx, subAssessments)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}

//// TRANSFORM FUNCTIONS

func extractAssessmentName(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	subAssessment := d.HydrateItem.(security.SubAssessment)

	r := regexp.MustCompile(`\bassessments\b`)
	splitStr := r.Split(*subAssessment.ID, len(*subAssessment.ID))[1]
	assessmentName := strings.Split(splitStr, "/")[1]
	return assessmentName, nil
}

func extractSQLServerVulnerabilityProperties(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	subAssessment := d.HydrateItem.(security.SubAssessment)
	additionalData := subAssessment.AdditionalData
	if additionalData == nil {
		return nil, nil
	}
	SQLServerVulnerabilityProperties, flag := additionalData.AsSQLServerVulnerabilityProperties()

	objectMap := make(map[string]interface{})
	if flag {
		if SQLServerVulnerabilityProperties.Type != nil {
			objectMap["Type"] = SQLServerVulnerabilityProperties.Type
		}
		if SQLServerVulnerabilityProperties.Query != nil {
			objectMap["Query"] = SQLServerVulnerabilityProperties.Query
		}
		if SQLServerVulnerabilityProperties.AssessedResourceType != "" {
			objectMap["AssessedResourceType"] = SQLServerVulnerabilityProperties.AssessedResourceType
		}
		jsonStr, err := json.Marshal(objectMap)
		if err != nil {
			return nil, nil
		}
		return string(jsonStr), nil
	}

	return nil, nil
}

func extractContainerRegistryVulnerabilityProperties(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	subAssessment := d.HydrateItem.(security.SubAssessment)
	additionalData := subAssessment.AdditionalData
	if additionalData == nil {
		return nil, nil
	}
	containerRegistryVulnerabilityProperties, flag := additionalData.AsContainerRegistryVulnerabilityProperties()

	objectMap := make(map[string]interface{})
	if flag {
		if containerRegistryVulnerabilityProperties.Type != nil {
			objectMap["Type"] = *containerRegistryVulnerabilityProperties.Type
		}
		if containerRegistryVulnerabilityProperties.Cvss != nil {
			objectMap["Cvss"] = extractCVSS(containerRegistryVulnerabilityProperties.Cvss)
		}
		if containerRegistryVulnerabilityProperties.Patchable != nil {
			objectMap["Patchable"] = *containerRegistryVulnerabilityProperties.Patchable
		}
		if containerRegistryVulnerabilityProperties.Cve != nil {
			objectMap["Cve"] = extractCVE(containerRegistryVulnerabilityProperties.Cve)
		}
		if containerRegistryVulnerabilityProperties.PublishedTime != nil {
			objectMap["PublishedTime"] = *containerRegistryVulnerabilityProperties.PublishedTime
		}
		if containerRegistryVulnerabilityProperties.VendorReferences != nil {
			objectMap["VendorReferences"] = extractVendorReferences(containerRegistryVulnerabilityProperties.VendorReferences)
		}
		if containerRegistryVulnerabilityProperties.RepositoryName != nil {
			objectMap["RepositoryName"] = *containerRegistryVulnerabilityProperties.RepositoryName
		}
		if containerRegistryVulnerabilityProperties.ImageDigest != nil {
			objectMap["ImageDigest"] = *containerRegistryVulnerabilityProperties.ImageDigest
		}
		if containerRegistryVulnerabilityProperties.AssessedResourceType != "" {
			objectMap["AssessedResourceType"] = containerRegistryVulnerabilityProperties.AssessedResourceType
		}
		jsonStr, err := json.Marshal(objectMap)
		if err != nil {
			return nil, nil
		}
		return string(jsonStr), nil
	}
	return nil, nil
}

func extractServerVulnerabilityProperties(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	subAssessment := d.HydrateItem.(security.SubAssessment)
	additionalData := subAssessment.AdditionalData
	if additionalData == nil {
		return nil, nil
	}
	serverVulnerabilityProperties, flag := additionalData.AsServerVulnerabilityProperties()

	objectMap := make(map[string]interface{})
	if flag {
		if serverVulnerabilityProperties.Type != nil {
			objectMap["Type"] = *serverVulnerabilityProperties.Type
		}
		if serverVulnerabilityProperties.Cvss != nil {
			objectMap["Cvss"] = extractCVSS(serverVulnerabilityProperties.Cvss)
		}
		if serverVulnerabilityProperties.Patchable != nil {
			objectMap["Patchable"] = serverVulnerabilityProperties.Patchable
		}
		if serverVulnerabilityProperties.Cve != nil {
			objectMap["Cve"] = extractCVE(serverVulnerabilityProperties.Cve)
		}
		if serverVulnerabilityProperties.PublishedTime != nil {
			objectMap["PublishedTime"] = serverVulnerabilityProperties.PublishedTime
		}
		if serverVulnerabilityProperties.VendorReferences != nil {
			objectMap["VendorReferences"] = extractVendorReferences(serverVulnerabilityProperties.VendorReferences)
		}
		if serverVulnerabilityProperties.Threat != nil {
			objectMap["Threat"] = serverVulnerabilityProperties.Threat
		}
		if serverVulnerabilityProperties.AssessedResourceType != "" {
			objectMap["AssessedResourceType"] = serverVulnerabilityProperties.AssessedResourceType
		}
		jsonStr, err := json.Marshal(objectMap)
		if err != nil {
			return nil, nil
		}
		return string(jsonStr), nil
	}

	return nil, nil
}

func extractAssessedResourceType(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	subAssessment := d.HydrateItem.(security.SubAssessment)
	additional := subAssessment.AdditionalData
	if additional == nil {
		return nil, nil
	}
	additionalData, _ := additional.AsAdditionalData()
	if additionalData == nil {
		return nil, nil
	}
	return additionalData.AssessedResourceType, nil
}

func extractCVE(cve *[]security.CVE) []map[string]interface{} {
	var cveop []map[string]interface{}
	for _, i := range *cve {
		objectMap := make(map[string]interface{})
		if i.Title != nil {
			objectMap["Title"] = *i.Title
		}
		if i.Link != nil {
			objectMap["Link"] = *i.Link
		}
		cveop = append(cveop, objectMap)
	}
	return cveop
}

func extractVendorReferences(vendorReferences *[]security.VendorReference) []map[string]interface{} {
	var vendorReferencesop []map[string]interface{}
	for _, i := range *vendorReferences {
		objectMap := make(map[string]interface{})
		if i.Title != nil {
			objectMap["Title"] = *i.Title
		}
		if i.Link != nil {
			objectMap["Link"] = *i.Link
		}
		vendorReferencesop = append(vendorReferencesop, objectMap)
	}
	return vendorReferencesop
}

func extractCVSS(CVSS map[string]*security.CVSS) map[string]interface{} {
	objectMap := make(map[string]interface{})
	for key, value := range CVSS {
		if value != nil && value.Base != nil {
			objectMap[key] = *value.Base
		}
	}
	return objectMap
}

func extractResourceDetails(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	subAssessment := d.HydrateItem.(security.SubAssessment)
	resourceDetails := subAssessment.SubAssessmentProperties.ResourceDetails
	if resourceDetails == nil {
		return nil, nil
	}
	azureResourceDetails, flag := resourceDetails.AsAzureResourceDetails()
	if flag {
		return extractAzureResourceDetails(azureResourceDetails), nil
	}
	onPremiseResourceDetails, flag := resourceDetails.AsOnPremiseResourceDetails()
	if flag {
		return extractOnPremiseResourceDetails(onPremiseResourceDetails), nil
	}
	onPremiseSQLResourceDetails, flag := resourceDetails.AsOnPremiseSQLResourceDetails()
	if flag {
		return extractOnPremiseSQLResourceDetails(onPremiseSQLResourceDetails), nil
	}
	resourceDetail, flag := resourceDetails.AsResourceDetails()
	if flag {
		return extractResourceDetail(resourceDetail), nil
	}
	return nil, nil
}

func extractAzureResourceDetails(azureResourceDetails *security.AzureResourceDetails) interface{} {
	objectMap := make(map[string]interface{})
	if azureResourceDetails.ID != nil {
		objectMap["ID"] = *azureResourceDetails.ID
	}
	if azureResourceDetails.Source != "" {
		objectMap["Source"] = azureResourceDetails.Source
	}
	jsonStr, err := json.Marshal(objectMap)
	if err != nil {
		return nil
	}
	return string(jsonStr)
}

func extractOnPremiseResourceDetails(onPremiseResourceDetails *security.OnPremiseResourceDetails) interface{} {
	objectMap := make(map[string]interface{})
	if onPremiseResourceDetails != nil {
		objectMap["MachineName"] = *onPremiseResourceDetails.MachineName
	}
	if onPremiseResourceDetails.Source != "" {
		objectMap["Source"] = onPremiseResourceDetails.Source
	}
	if onPremiseResourceDetails.SourceComputerID != nil {
		objectMap["SourceComputerID"] = *onPremiseResourceDetails.SourceComputerID
	}
	if onPremiseResourceDetails.WorkspaceID != nil {
		objectMap["WorkspaceID"] = *onPremiseResourceDetails.WorkspaceID
	}
	if onPremiseResourceDetails.Vmuuid != nil {
		objectMap["Vmuuid"] = *onPremiseResourceDetails.Vmuuid
	}
	jsonStr, err := json.Marshal(objectMap)
	if err != nil {
		return nil
	}
	return string(jsonStr)
}

func extractOnPremiseSQLResourceDetails(onPremiseSQLResourceDetails *security.OnPremiseSQLResourceDetails) interface{} {
	objectMap := make(map[string]interface{})
	if onPremiseSQLResourceDetails != nil {
		objectMap["MachineName"] = *onPremiseSQLResourceDetails.MachineName
	}
	if onPremiseSQLResourceDetails.Source != "" {
		objectMap["Source"] = onPremiseSQLResourceDetails.Source
	}
	if onPremiseSQLResourceDetails.SourceComputerID != nil {
		objectMap["SourceComputerID"] = *onPremiseSQLResourceDetails.SourceComputerID
	}
	if onPremiseSQLResourceDetails.WorkspaceID != nil {
		objectMap["WorkspaceID"] = *onPremiseSQLResourceDetails.WorkspaceID
	}
	if onPremiseSQLResourceDetails.Vmuuid != nil {
		objectMap["Vmuuid"] = *onPremiseSQLResourceDetails.Vmuuid
	}
	if onPremiseSQLResourceDetails.ServerName != nil {
		objectMap["ServerName"] = *onPremiseSQLResourceDetails.ServerName
	}
	if onPremiseSQLResourceDetails.DatabaseName != nil {
		objectMap["DatabaseName"] = *onPremiseSQLResourceDetails.DatabaseName
	}
	jsonStr, err := json.Marshal(objectMap)
	if err != nil {
		return nil
	}
	return string(jsonStr)
}

func extractResourceDetail(resourceDetail *security.ResourceDetails) interface{} {
	objectMap := make(map[string]interface{})
	if resourceDetail.Source != "" {
		objectMap["Source"] = resourceDetail.Source
	}
	jsonStr, err := json.Marshal(objectMap)
	if err != nil {
		return nil
	}
	return string(jsonStr)
}

func extractSubAssessmentStatus(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	subAssessment := d.HydrateItem.(security.SubAssessment)
	subAssessmentStatus := subAssessment.Status
	objectMap := make(map[string]interface{})
	if subAssessmentStatus.Cause != nil {
		objectMap["Cause"] = *subAssessmentStatus.Cause
	}
	if subAssessmentStatus.Code != "" {
		objectMap["Code"] = subAssessmentStatus.Code
	}
	if subAssessmentStatus.Description != nil {
		objectMap["Description"] = *subAssessmentStatus.Description
	}
	if subAssessmentStatus.Severity != "" {
		objectMap["Severity"] = subAssessmentStatus.Severity
	}
	jsonStr, err := json.Marshal(objectMap)
	if err != nil {
		return nil, err
	}
	return string(jsonStr), nil
}

func extractResourceGroupFromSubAssessmentID(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	id := types.SafeString(d.Value)

	// Common resource properties
	if !strings.Contains(id, "resourceGroups") {
		return nil, nil
	}
	splitID := strings.Split(id, "/")
	resourceGroup := splitID[4]
	resourceGroup = strings.ToLower(resourceGroup)
	return resourceGroup, nil
}
