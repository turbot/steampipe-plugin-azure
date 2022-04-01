select name, id, storage_account_name 
from azure_storage_share_file
where name = '{{resourceName}}' and resource_group = '{{resourceName}}' and storage_account_name = '{{resourceName}}';