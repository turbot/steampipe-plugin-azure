select name, id, type, region, resource_group, subscription_id
from azure.azure_batch_account
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';