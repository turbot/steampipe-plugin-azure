select id
from azure.azure_resource_graph
where query = 'Resources | where id == "nonexistent-resource-id-xyz"'
