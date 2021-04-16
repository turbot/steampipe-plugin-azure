select name, id, type, location, sku, identity, dns_prefix, enable_rbac, fqdn, max_agent_pools, node_resource_group, provisioning_state, addon_profiles, api_server_access_profile, power_state, title, tags, akas, region, resource_group, subscription_id
from azure.azure_kubernetes_cluster
where name = '{{resourceName}}' and resource_group = '{{resourceName}}';