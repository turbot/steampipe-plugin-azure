select name, id, type, resource_group
from azure.azure_kubernetes_cluster
where akas::text = '["{{output.resource_aka.value}}", "{{output.resource_aka_lower.value}}"]';