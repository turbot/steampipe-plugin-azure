select name, akas, title
from azure.azure_storage_account_local_user
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}' and storage_account_name = '{{ resourceName }}'; 