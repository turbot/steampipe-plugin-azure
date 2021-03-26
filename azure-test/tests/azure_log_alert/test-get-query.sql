select name, id, region, type
from azure.azure_log_alert
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
