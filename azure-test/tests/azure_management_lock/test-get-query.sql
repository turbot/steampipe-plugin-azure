select name, id, lock_level, type, resource_group, scope, notes
from azure.azure_management_lock
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
