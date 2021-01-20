select name, id, type, region
from azure.azure_network_security_group
where name = 'dummy-test794612891' and resource_group = '{{resourceName}}'
