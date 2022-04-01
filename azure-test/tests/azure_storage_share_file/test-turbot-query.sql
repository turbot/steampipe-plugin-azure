select title, akas
from azure_storage_share_file
where name = '{{resourceName}}' and resource_group = '{{resourceName}}';