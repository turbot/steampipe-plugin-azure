select name,
    id,
    region,
    type
from azure.azure_dns_zone
where name = '{{resourceName}}.xyz'
