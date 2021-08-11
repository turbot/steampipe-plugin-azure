select name, title, akas
from azure.azure_recovery_services_vault
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';