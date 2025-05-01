select name, id, type, resource_group, storage_account_name
from azure.azure_storage_account_local_user
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}' and storage_account_name = '{{ resourceName }}'; 