select name, id, type
from azure.azure_logic_app_workflow
where id = '{{ output.resource_id.value }}';