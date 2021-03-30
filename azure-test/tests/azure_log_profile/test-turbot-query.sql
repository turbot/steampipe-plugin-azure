select name, akas, title
from azure.azure_log_profile
where name = '{{resourceName}}'
