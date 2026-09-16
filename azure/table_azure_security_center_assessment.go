package azure

import (
	"context"
	"net/http"
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
				&plugin.KeyColumn{
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
				Transform:   transform.FromField("AssessmentPropertiesResponse.ResourceDetails").Transform(extractAssessmentResourceDetails),
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
				Transform:   transform.FromField("AssessmentPropertiesResponse.Links"),
			},
			{
				Name:        "metadata",
				Description: "Describes properties of an assessment metadata. By default this is not populated, unless it's specified in expand.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("AssessmentPropertiesResponse.Metadata"),
			},
			{
				Name:        "resource_name",
				Description: "Name of the resource that was assessed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AssessmentPropertiesResponse.ResourceDetails").Transform(extractAssessmentResourceName),
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

	var result security.AssessmentListPage
	expand := d.EqualsQualString("expand")
	if expand == "" {
		result, err = assessmentClient.List(ctx, "subscriptions/"+subscriptionID)
		if err != nil {
			logger.Error("azure_security_center_assessment.listSecurityCenterAssessments", "query_error", err)
			return nil, err
		}
	} else {
		result, err = listSecurityCenterAssessmentsExpanded(ctx, &assessmentClient, "subscriptions/"+subscriptionID, expand)
		if err != nil {
			logger.Error("azure_security_center_assessment.listSecurityCenterAssessments", "query_error", err)
			return nil, err
		}
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

// listSecurityCenterAssessmentsExpanded lists assessments with the "$expand" query parameter
// added to the request, so that "links"/"metadata" are populated in the same response instead of
// requiring a separate Get per row. The generated SDK's List method doesn't expose an "expand"
// parameter (only Get does), so the request/response are built manually, reusing the client's
// preparer/sender/responder plumbing (which also takes care of authorization).
func listSecurityCenterAssessmentsExpanded(ctx context.Context, assessmentClient *security.AssessmentsClient, scope string, expand string) (security.AssessmentListPage, error) {
	req, err := assessmentClient.ListPreparer(ctx, scope)
	if err != nil {
		return security.AssessmentListPage{}, err
	}
	q := req.URL.Query()
	q.Set("$expand", expand)
	req.URL.RawQuery = q.Encode()

	resp, err := assessmentClient.ListSender(req)
	if err != nil {
		return security.AssessmentListPage{}, err
	}
	firstPage, err := assessmentClient.ListResponder(resp)
	if err != nil {
		return security.AssessmentListPage{}, err
	}

	return security.NewAssessmentListPage(firstPage, func(ctx context.Context, lastResults security.AssessmentList) (security.AssessmentList, error) {
		if lastResults.NextLink == nil || *lastResults.NextLink == "" {
			return security.AssessmentList{}, nil
		}
		nextReq, err := http.NewRequestWithContext(ctx, http.MethodGet, *lastResults.NextLink, nil)
		if err != nil {
			return security.AssessmentList{}, err
		}
		nextResp, err := assessmentClient.ListSender(nextReq)
		if err != nil {
			return security.AssessmentList{}, err
		}
		return assessmentClient.ListResponder(nextResp)
	}), nil
}

func getSecurityCenterAssessment(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		logger.Error("azure_security_center_assessment.getSecurityCenterAssessment", "connection_error", err)
		return nil, err
	}

	assessmentID := d.EqualsQualString("id")
	if assessmentID == "" {
		return nil, nil
	}
	name := getLastPathElement(assessmentID)

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

func extractAssessmentResourceDetails(_ context.Context, d *transform.TransformData) (interface{}, error) {
	resourceDetails, ok := d.Value.(security.BasicResourceDetails)
	if !ok || resourceDetails == nil {
		return nil, nil
	}

	if azureResourceDetails, flag := resourceDetails.AsAzureResourceDetails(); flag {
		return extractAzureResourceDetails(azureResourceDetails), nil
	}
	if onPremiseResourceDetails, flag := resourceDetails.AsOnPremiseResourceDetails(); flag {
		return extractOnPremiseResourceDetails(onPremiseResourceDetails), nil
	}
	if onPremiseSQLResourceDetails, flag := resourceDetails.AsOnPremiseSQLResourceDetails(); flag {
		return extractOnPremiseSQLResourceDetails(onPremiseSQLResourceDetails), nil
	}
	if resourceDetail, flag := resourceDetails.AsResourceDetails(); flag {
		return extractResourceDetail(resourceDetail), nil
	}

	return nil, nil
}

// getSecurityCenterAssessmentExpand builds the ExpandEnum value to send to the Security Center
// API from the "expand" qual, e.g. expand = 'links' or expand = 'links,metadata'.
func getSecurityCenterAssessmentExpand(d *plugin.QueryData) security.ExpandEnum {
	return security.ExpandEnum(d.EqualsQualString("expand"))
}

func extractAssessmentResourceName(_ context.Context, d *transform.TransformData) (interface{}, error) {
	resourceDetails, ok := d.Value.(security.BasicResourceDetails)
	if !ok || resourceDetails == nil {
		return nil, nil
	}

	if azureResourceDetails, flag := resourceDetails.AsAzureResourceDetails(); flag && azureResourceDetails.ID != nil {
		return getLastPathElement(*azureResourceDetails.ID), nil
	}
	if onPremiseResourceDetails, flag := resourceDetails.AsOnPremiseResourceDetails(); flag && onPremiseResourceDetails.MachineName != nil {
		return *onPremiseResourceDetails.MachineName, nil
	}
	if onPremiseSQLResourceDetails, flag := resourceDetails.AsOnPremiseSQLResourceDetails(); flag && onPremiseSQLResourceDetails.MachineName != nil {
		return *onPremiseSQLResourceDetails.MachineName, nil
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

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
