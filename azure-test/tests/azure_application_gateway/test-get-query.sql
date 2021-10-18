select name, id, type, region, resource_group, subscription_id
from azure.azure_application_gateway
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
