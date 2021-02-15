select id, name, location
from azure.azure_subnet
where name = '{{resourceName}}'
