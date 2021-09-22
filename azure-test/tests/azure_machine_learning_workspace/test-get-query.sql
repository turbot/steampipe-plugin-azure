select name, id, type, tags, region, resource_group, subscription_id
from azure.azure_machine_learning_workspace
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';