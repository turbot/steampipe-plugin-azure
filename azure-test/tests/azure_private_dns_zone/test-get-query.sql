select name,
    id,
    region,
    type
from azure.azure_private_dns_zone
where name = '{{resourceName}}.local'
