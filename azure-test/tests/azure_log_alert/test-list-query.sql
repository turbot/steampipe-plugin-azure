select id, name
from azure.azure_log_alert
where name = '{{resourceName}}'
