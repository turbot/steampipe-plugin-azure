select name, id, type, resource_group, storage_account_name
from azure.azure_storage_account_local_user
where id = '{{ output.resource_id.value }}'; 