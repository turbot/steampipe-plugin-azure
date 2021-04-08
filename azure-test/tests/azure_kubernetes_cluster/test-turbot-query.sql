select name, akas, title, tags
from azure.azure_kubernetes_cluster
where name = '{{resourceName}}' and resource_group = '{{resourceName}}';