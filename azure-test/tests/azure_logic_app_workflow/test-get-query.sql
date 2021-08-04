select name, id, type, region, resource_group, subscription_id
from azure.azure_logic_app_workflow
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';