select name, id
from azure_storage_share_file
where name = 'dummy{{resourceName}}' and resource_group = '{{resourceName}}';