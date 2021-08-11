package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2021-02-01/containerservice"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAzureKubernetesCluster(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "azure_kubernetes_cluster",
		Description: "Azure Kubernetes Cluster",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"name", "resource_group"}),
			Hydrate:    getKubernetesCluster,
		},
		List: &plugin.ListConfig{
			Hydrate: listKubernetesClusters,
		},
		Columns: []*plugin.Column{
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the cluster.",
			},
			{
				Name:        "id",
				Description: "The ID of the cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "type",
				Description: "The type of the cluster.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "location",
				Description: "The location where the cluster is created.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "azure_portal_fqdn",
				Description: "FQDN for the master pool which used by proxy config.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedClusterProperties.AzurePortalFQDN"),
			},
			{
				Name:        "disk_encryption_set_id",
				Description: "ResourceId of the disk encryption set to use for enabling encryption at rest.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedClusterProperties.DiskEncryptionSetID"),
			},
			{
				Name:        "dns_prefix",
				Description: "DNS prefix specified when creating the managed cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedClusterProperties.DNSPrefix"),
			},
			{
				Name:        "enable_pod_security_policy",
				Description: "Whether to enable Kubernetes pod security policy (preview).",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ManagedClusterProperties.EnablePodSecurityPolicy"),
			},
			{
				Name:        "enable_rbac",
				Description: "Whether to enable Kubernetes Role-Based Access Control.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromField("ManagedClusterProperties.EnableRBAC"),
			},
			{
				Name:        "fqdn",
				Description: "FQDN for the master pool.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedClusterProperties.Fqdn"),
			},
			{
				Name:        "fqdn_subdomain",
				Description: "FQDN subdomain specified when creating private cluster with custom private dns zone.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedClusterProperties.FqdnSubdomain"),
			},
			{
				Name:        "kubernetes_version",
				Description: "Version of Kubernetes specified when creating the managed cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedClusterProperties.KubernetesVersion"),
			},
			{
				Name:        "max_agent_pools",
				Description: "The max number of agent pools for the managed cluster.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("ManagedClusterProperties.MaxAgentPools"),
			},
			{
				Name:        "node_resource_group",
				Description: "Name of the resource group containing agent pool nodes.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedClusterProperties.NodeResourceGroup"),
			},
			{
				Name:        "private_fqdn",
				Description: "FQDN of private cluster.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedClusterProperties.PrivateFQDN"),
			},
			{
				Name:        "provisioning_state",
				Description: "The current deployment or provisioning state.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ManagedClusterProperties.ProvisioningState"),
			},
			{
				Name:        "aad_profile",
				Description: "Profile of Azure Active Directory configuration.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ManagedClusterProperties.AadProfile"),
			},
			{
				Name:        "addon_profiles",
				Description: "Profile of managed cluster add-on.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ManagedClusterProperties.AddonProfiles"),
			},
			{
				Name:        "agent_pool_profiles",
				Description: "Properties of the agent pool.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ManagedClusterProperties.AgentPoolProfiles"),
			},
			{
				Name:        "api_server_access_profile",
				Description: "Access profile for managed cluster API server.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ManagedClusterProperties.APIServerAccessProfile"),
			},
			{
				Name:        "auto_scaler_profile",
				Description: "Parameters to be applied to the cluster-autoscaler when enabled.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ManagedClusterProperties.AutoScalerProfile"),
			},
			{
				Name:        "auto_upgrade_profile",
				Description: "Profile of auto upgrade configuration.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ManagedClusterProperties.AutoUpgradeProfile"),
			},
			{
				Name:        "identity",
				Description: "The identity of the managed cluster, if configured.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "identity_profile",
				Description: "Identities associated with the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ManagedClusterProperties.IdentityProfile"),
			},
			{
				Name:        "linux_profile",
				Description: "Profile for Linux VMs in the container service cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ManagedClusterProperties.LinuxProfile"),
			},
			{
				Name:        "network_profile",
				Description: "Profile of network configuration.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ManagedClusterProperties.NetworkProfile"),
			},
			{
				Name:        "pod_identity_profile",
				Description: "Profile of managed cluster pod identity.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ManagedClusterProperties.PodIdentityProfile"),
			},
			{
				Name:        "power_state",
				Description: "Represents the Power State of the cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ManagedClusterProperties.PowerState"),
			},
			{
				Name:        "service_principal_profile",
				Description: "Information about a service principal identity for the cluster to use for manipulating Azure APIs.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ManagedClusterProperties.ServicePrincipalProfile"),
			},
			{
				Name:        "sku",
				Description: "The managed cluster SKU.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "windows_profile",
				Description: "Profile for Windows VMs in the container service cluster.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("ManagedClusterProperties.WindowsProfile"),
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

func listKubernetesClusters(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := containerservice.NewManagedClustersClient(subscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, cluster := range result.Values() {
		d.StreamListItem(ctx, cluster)
	}

	for result.NotDone() {
		err = result.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}

		for _, cluster := range result.Values() {
			d.StreamListItem(ctx, cluster)
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getKubernetesCluster(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getKubernetesCluster")

	session, err := GetNewSession(ctx, d, "MANAGEMENT")
	if err != nil {
		return nil, err
	}
	subscriptionID := session.SubscriptionID

	client := containerservice.NewManagedClustersClient(subscriptionID)
	client.Authorizer = session.Authorizer

	resourceName := d.KeyColumnQuals["name"].GetStringValue()
	resourceGroupName := d.KeyColumnQuals["resource_group"].GetStringValue()

	op, err := client.Get(ctx, resourceGroupName, resourceName)
	if err != nil {
		return nil, err
	}

	// In some cases resource does not give any notFound error
	// instead of notFound error, it returns empty data
	if op.ID != nil {
		return op, nil
	}

	return nil, nil
}
