select name, akas, title
from azure.azure_log_alert
where name = 'dummy-{{resourceName}}' and resource_group = '{{resourceName}}'
