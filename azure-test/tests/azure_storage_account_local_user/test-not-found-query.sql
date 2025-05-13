select name, id, type, resource_group
from azure.azure_storage_account_local_user
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}' and storage_account_name = '{{ resourceName }}'; 