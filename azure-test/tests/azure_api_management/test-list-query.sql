select name, id
from azure.azure_api_management
where id = '{{ output.resource_id.value }}';
