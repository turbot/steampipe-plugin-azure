select name, id
from azure.azure_log_profile
where name = '{{ resourceName }}';
