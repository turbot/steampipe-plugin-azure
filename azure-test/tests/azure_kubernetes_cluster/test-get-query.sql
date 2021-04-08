select name, id, type, location, identity, managed_cluster_properties -> 'enableRBAC' as rbac, sku, title, tags, akas, region, resource_group, subscription_id
from azure.azure_kubernetes_cluster
where name = '{{resourceName}}' and resource_group = '{{resourceName}}';