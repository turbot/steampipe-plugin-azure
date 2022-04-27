package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v1.0/security"
	"github.com/turbot/steampipe-plugin-sdk/v2/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v2/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/v2/plugin"
)

//// TABLE DEFINITION

func tableAzureSecurityCenterSubAssessment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_security_center_sub_assessment",
		Description: "Azure Security Center Sub Assessment",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "assessment_name"}),
			Hydrate:    getSecurityCenterSubAssessment,
		},
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
				Description: "Assessment name.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(getAssessmentName),
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
				Name:        "status",
				Description: "The status of the sub-assessment.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SubAssessmentProperties.Status"),
			},
			{
				Name:        "time_generated",
				Description: "The date and time the sub-assessment was generated.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SubAssessmentProperties.TimeGenerated"),
			},
			{
				Name:        "additional_data",
				Description: "Details of the sub-assessment.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SubAssessmentProperties.AdditionalData"),
			},
			{
				Name:        "resource_details",
				Description: "Details of the resource that was assessed.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("SubAssessmentProperties.ResourceDetails"),
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
		}),
	}
}

//// LIST FUNCTION

func listSecurityCenterSubAssessments(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Error("azure_security_center_sub_assessment.listSecurityCenterSubAssessments")
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

//// HYDRATE FUNCTIONS

func getSecurityCenterSubAssessment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Error("azure_security_center_sub_assessment.getSecurityCenterSubAssessment")
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		logger.Error("azure_security_center_sub_assessment.getSecurityCenterSubAssessment", "connection_error", err)
		return nil, err
	}

	name := d.KeyColumnQuals["name"].GetStringValue()
	assessmentName := d.KeyColumnQuals["assessment_name"].GetStringValue()
	subscriptionID := session.SubscriptionID
	subAssessmentClient := security.NewSubAssessmentsClientWithBaseURI(session.ResourceManagerEndpoint, subscriptionID, "")
	subAssessmentClient.Authorizer = session.Authorizer

	subAssessment, err := subAssessmentClient.Get(ctx, "subscriptions/"+subscriptionID, assessmentName, name)
	if err != nil {
		logger.Error("azure_security_center_sub_assessment.getSecurityCenterSubAssessment", "query_error", err)
		return err, nil
	}
	return subAssessment, nil
}

//// TRANSFORM FUNCTIONS

func getAssessmentName(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	subAssessment := d.HydrateItem.(security.SubAssessment)
	assessmentName := strings.Split(string(*subAssessment.ID), "/")[6]
	return assessmentName, nil
}
