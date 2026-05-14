package azure

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/security/mgmt/security"
	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
)

//// TABLE DEFINITION

func tableAzureSecurityCenterAssessment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_security_center_assessment",
		Description: "Azure Security Center Assessment",
		List: &plugin.ListConfig{
			Hydrate: listSecurityCenterAssessments,
			Tags: map[string]string{
				"service": "Microsoft.Security",
				"action":  "assessments/read",
			},
			KeyColumns: plugin.KeyColumnSlice{
				{
					Name:      "expand",
					Operators: []string{"="},
					Require:   plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getSecurityCenterAssessment,
			Tags: map[string]string{
				"service": "Microsoft.Security",
				"action":  "assessments/read",
			},
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
				Name:        "display_name",
				Description: "User friendly display name of the assessment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssessmentPropertiesResponse.DisplayName"),
			},
			{
				Name:        "status_cause",
				Description: "Status cause of the assessment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssessmentPropertiesResponse.Status.Cause"),
			},
			{
				Name:        "status_code",
				Description: "Status code of the assessment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssessmentPropertiesResponse.Status.Code"),
			},
			{
				Name:        "status_description",
				Description: "Status description of the assessment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssessmentPropertiesResponse.Status.Description"),
			},
			{
				Name:        "status_change_date",
				Description: "Status change date of the assessment.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("AssessmentPropertiesResponse.Status.StatusChangeDate").Transform(convertDateToTime),
			},
			{
				Name:        "status_first_evaluation_date",
				Description: "The timestamp of the first evaluation of the assessment.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("AssessmentPropertiesResponse.Status.FirstEvaluationDate").Transform(convertDateToTime),
			},
			{
				Name:        "resource_details",
				Description: "Details of the resource that was assessed.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSecurityCenterAssessmentProperties,
				Transform:   transform.FromValue().Transform(getSecurityCenterAssessmentResourceDetailsFromProperties),
			},
			{
				Name:        "additional_data",
				Description: "Additional data regarding the assessment.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AssessmentPropertiesResponse.AdditionalData"),
			},
			{
				Name:        "expand",
				Description: "May be used to expand the 'links' or 'metadata' of the assessment. By default, these fields are not included when listing or getting assessments.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("expand"),
			},
			{
				Name:        "links",
				Description: "Links relevant to the assessment. By default this is not populated, unless it's specified in expand.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSecurityCenterAssessmentProperties,
				Transform:   transform.FromValue().Transform(getSecurityCenterAssessmentLinksFromProperties),
			},
			{
				Name:        "metadata",
				Description: "Describes properties of an assessment metadata. By default this is not populated, unless it's specified in expand.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getSecurityCenterAssessmentProperties,
				Transform:   transform.FromValue().Transform(getSecurityCenterAssessmentMetadataFromProperties),
			},
			{
				Name:        "resource_name",
				Description: "Name of the resource that was assessed.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getSecurityCenterAssessmentProperties,
				Transform:   transform.FromValue().Transform(getSecurityCenterAssessmentResourceNameFromProperties),
			},
			{
				Name:        "resource_id",
				Description: "The full ID of the assessed resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID").Transform(extractAssessedResourceIDFromAssessmentID),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssessmentPropertiesResponse.DisplayName"),
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

func listSecurityCenterAssessments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		logger.Error("azure_security_center_assessment.listSecurityCenterAssessments", "connection_error", err)
		return nil, err
	}

	subscriptionID := session.SubscriptionID
	assessmentClient := security.NewAssessmentsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID)
	assessmentClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &assessmentClient, d.Connection)

	result, err := assessmentClient.List(ctx, "subscriptions/"+subscriptionID)
	if err != nil {
		logger.Error("azure_security_center_assessment.listSecurityCenterAssessments", "query_error", err)
		return nil, err
	}

	for _, assessments := range result.Values() {
		d.StreamListItem(ctx, assessments)
		// Check if context has been cancelled or if the limit has been hit (if specified)
		// if there is a limit, it will return the number of rows required to reach this limit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	for result.NotDone() {
		// Wait for rate limiting
		d.WaitForListRateLimit(ctx)

		err = result.NextWithContext(ctx)
		if err != nil {
			logger.Error("azure_security_center_assessment.listSecurityCenterAssessments", "query_error", err)
			return nil, err
		}
		for _, assessments := range result.Values() {
			d.StreamListItem(ctx, assessments)
			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, nil
}

func getSecurityCenterAssessment(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		logger.Error("azure_security_center_assessment.getSecurityCenterAssessment", "connection_error", err)
		return nil, err
	}

	assessmentID := d.EqualsQualString("id")
	name := d.EqualsQualString("name")
	if assessmentID == "" || name == "" {
		return nil, nil
	}

	resourceID := assessmentID
	securityProviderIndex := strings.Index(strings.ToLower(assessmentID), "/providers/microsoft.security/assessments/")
	if securityProviderIndex > 0 {
		resourceID = assessmentID[:securityProviderIndex]
	}

	assessmentClient := security.NewAssessmentsClientWithBaseURI(session.ResourceManagerEndpoint, session.SubscriptionID)
	assessmentClient.Authorizer = session.Authorizer

	// Apply Retry rule
	ApplyRetryRules(ctx, &assessmentClient, d.Connection)

	assessment, err := assessmentClient.Get(ctx, resourceID, name, getSecurityCenterAssessmentExpand(d))
	if err != nil {
		logger.Error("azure_security_center_assessment.getSecurityCenterAssessment", "query_error", err)
		return nil, err
	}

	return assessment, nil
}

func getSecurityCenterAssessmentResourceDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	propertiesRaw, err := getSecurityCenterAssessmentProperties(ctx, d, h)
	if err != nil {
		return nil, err
	}
	properties, ok := propertiesRaw.(map[string]interface{})
	if !ok || properties == nil {
		return nil, nil
	}
	if resourceDetails, ok := properties["resourceDetails"]; ok {
		return resourceDetails, nil
	}
	return nil, nil
}

func getSecurityCenterAssessmentResourceDetailsFromProperties(_ context.Context, d *transform.TransformData) (interface{}, error) {
	properties, ok := d.Value.(map[string]interface{})
	if !ok || properties == nil {
		return nil, nil
	}

	if resourceDetails, ok := properties["resourceDetails"]; ok {
		return resourceDetails, nil
	}
	if resourceDetails, ok := properties["ResourceDetails"]; ok {
		return resourceDetails, nil
	}

	return nil, nil
}

// getSecurityCenterAssessmentLinksFromProperties extracts the "links" attribute from the raw
// assessment properties. The API only populates this attribute when it's requested via the
// "expand" qual (e.g. expand = 'links' or expand = 'links,metadata').
func getSecurityCenterAssessmentLinksFromProperties(_ context.Context, d *transform.TransformData) (interface{}, error) {
	properties, ok := d.Value.(map[string]interface{})
	if !ok || properties == nil {
		return nil, nil
	}

	if links, ok := properties["links"]; ok {
		return links, nil
	}
	if links, ok := properties["Links"]; ok {
		return links, nil
	}

	return nil, nil
}

// getSecurityCenterAssessmentMetadataFromProperties extracts the "metadata" attribute from the raw
// assessment properties. The API only populates this attribute when it's requested via the
// "expand" qual (e.g. expand = 'metadata' or expand = 'links,metadata').
func getSecurityCenterAssessmentMetadataFromProperties(_ context.Context, d *transform.TransformData) (interface{}, error) {
	properties, ok := d.Value.(map[string]interface{})
	if !ok || properties == nil {
		return nil, nil
	}

	if metadata, ok := properties["metadata"]; ok {
		return metadata, nil
	}
	if metadata, ok := properties["Metadata"]; ok {
		return metadata, nil
	}

	return nil, nil
}

// getSecurityCenterAssessmentExpand builds the ExpandEnum value to send to the Security Center
// API from the "expand" qual, e.g. expand = 'links' or expand = 'links,metadata'.
func getSecurityCenterAssessmentExpand(d *plugin.QueryData) security.ExpandEnum {
	return security.ExpandEnum(d.EqualsQualString("expand"))
}

func getSecurityCenterAssessmentResourceNameFromProperties(_ context.Context, d *transform.TransformData) (interface{}, error) {
	properties, ok := d.Value.(map[string]interface{})
	if !ok || properties == nil {
		return nil, nil
	}

	resourceDetails, ok := properties["resourceDetails"].(map[string]interface{})
	if !ok {
		resourceDetails, _ = properties["ResourceDetails"].(map[string]interface{})
	}
	if resourceDetails == nil {
		return nil, nil
	}

	if resourceName, ok := resourceDetails["resourceName"].(string); ok && resourceName != "" {
		return resourceName, nil
	}
	if resourceName, ok := resourceDetails["ResourceName"].(string); ok && resourceName != "" {
		return resourceName, nil
	}
	if id, ok := resourceDetails["id"].(string); ok && id != "" {
		splitID := strings.Split(strings.TrimSuffix(id, "/"), "/")
		return splitID[len(splitID)-1], nil
	}
	if id, ok := resourceDetails["ID"].(string); ok && id != "" {
		splitID := strings.Split(strings.TrimSuffix(id, "/"), "/")
		return splitID[len(splitID)-1], nil
	}
	if machineName, ok := resourceDetails["machineName"].(string); ok && machineName != "" {
		return machineName, nil
	}
	if machineName, ok := resourceDetails["MachineName"].(string); ok && machineName != "" {
		return machineName, nil
	}

	return nil, nil
}

func getSecurityCenterAssessmentProperties(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	assessmentID, assessmentName := getAssessmentIdentity(h.Item)
	if assessmentID == "" {
		assessmentID = d.EqualsQualString("id")
	}
	if assessmentName == "" {
		assessmentName = d.EqualsQualString("name")
	}
	if assessmentID == "" || assessmentName == "" {
		return nil, nil
	}

	resourceID := assessmentID
	securityProviderIndex := strings.Index(strings.ToLower(assessmentID), "/providers/microsoft.security/assessments/")
	if securityProviderIndex > 0 {
		resourceID = assessmentID[:securityProviderIndex]
	}

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		logger.Error("azure_security_center_assessment.getSecurityCenterAssessmentProperties", "connection_error", err)
		return nil, err
	}

	assessmentClient := security.NewAssessmentsClientWithBaseURI(session.ResourceManagerEndpoint, session.SubscriptionID)
	assessmentClient.Authorizer = session.Authorizer

	ApplyRetryRules(ctx, &assessmentClient, d.Connection)

	req, err := assessmentClient.GetPreparer(ctx, resourceID, assessmentName, getSecurityCenterAssessmentExpand(d))
	if err != nil {
		logger.Error("azure_security_center_assessment.getSecurityCenterAssessmentProperties", "prepare_error", err)
		return nil, err
	}

	resp, err := assessmentClient.GetSender(req)
	if err != nil {
		logger.Error("azure_security_center_assessment.getSecurityCenterAssessmentProperties", "query_error", err)
		return nil, err
	}
	defer resp.Body.Close()

	var rawAssessment map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&rawAssessment); err != nil {
		if err != io.EOF {
			logger.Error("azure_security_center_assessment.getSecurityCenterAssessmentProperties", "decode_error", err)
			return nil, err
		}
		return nil, nil
	}

	properties, ok := rawAssessment["properties"].(map[string]interface{})
	if !ok {
		return nil, nil
	}
	return properties, nil
}

//// TRANSFORM FUNCTIONS

func getAssessmentIdentity(item interface{}) (string, string) {
	switch assessment := item.(type) {
	case security.AssessmentResponse:
		return types.SafeString(assessment.ID), types.SafeString(assessment.Name)
	case security.Assessment:
		return types.SafeString(assessment.ID), types.SafeString(assessment.Name)
	default:
		return "", ""
	}
}

func extractAssessedResourceIDFromAssessmentID(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	assessmentID := types.SafeString(d.Value)
	if assessmentID == "" {
		return nil, nil
	}

	securityProviderIndex := strings.Index(strings.ToLower(assessmentID), "/providers/microsoft.security/assessments/")
	if securityProviderIndex > 0 {
		return assessmentID[:securityProviderIndex], nil
	}

	return assessmentID, nil
}
