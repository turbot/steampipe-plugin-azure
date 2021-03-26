select name, akas, title
from azure.azure_log_alert
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
