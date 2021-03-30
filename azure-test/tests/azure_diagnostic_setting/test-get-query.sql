select name, id,type
from azure.azure_diagnostic_setting
where name = '{{resourceName}}'
