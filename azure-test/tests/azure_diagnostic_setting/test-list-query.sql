select id, name
from azure.azure_diagnostic_setting
where title = '{{resourceName}}'
