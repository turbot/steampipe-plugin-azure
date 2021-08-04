select name, id, type, region
from azure.azure_lb
where name = 'dummy-test794612891' and resource_group = '{{resourceName}}'
