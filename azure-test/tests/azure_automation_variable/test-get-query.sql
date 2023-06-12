select name, id, type, resource_group, subscription_id
from azure.azure_automation_variable
where account_name = '{{ resourceName }}' and name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';