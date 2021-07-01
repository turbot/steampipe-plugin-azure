select name, tags, title, akas
from azure.azure_data_factory
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'