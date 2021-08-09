select name, id, type
from azure.azure_recovery_services_vault
where name = '{{ resourceName }}';