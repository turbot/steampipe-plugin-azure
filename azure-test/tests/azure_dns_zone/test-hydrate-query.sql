select name,
    akas,
    tags,
    title
from azure.azure_dns_zone
where name = '{{resourceName}}.xyz'
