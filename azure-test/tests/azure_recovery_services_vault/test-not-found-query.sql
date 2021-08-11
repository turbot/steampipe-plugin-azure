select name, tags, title, akas
from azure.azure_recovery_services_vault
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';