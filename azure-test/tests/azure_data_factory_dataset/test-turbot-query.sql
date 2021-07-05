select name, title, akas
from azure.azure_data_factory_dataset
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'