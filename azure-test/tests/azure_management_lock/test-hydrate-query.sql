select name, lock_level, akas, title
from azure.azure_management_lock
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
