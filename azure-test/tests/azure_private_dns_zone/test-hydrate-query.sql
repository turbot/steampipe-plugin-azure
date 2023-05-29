select name,
    akas,
    tags,
    title
from azure.azure_private_dns_zone
where name = '{{resourceName}}.local'
