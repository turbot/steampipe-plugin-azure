select name, id, type, tags, region, resource_group, subscription_id
from azure.azure_data_factory
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';