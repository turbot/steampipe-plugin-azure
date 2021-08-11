select name, id, type, region, resource_group, subscription_id
from azure.azure_recovery_services_vault
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';