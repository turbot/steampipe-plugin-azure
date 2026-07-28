select id, name, type, query
from azure.azure_resource_graph
where query = 'ResourceContainers | where type == "microsoft.resources/subscriptions" | limit 1'
