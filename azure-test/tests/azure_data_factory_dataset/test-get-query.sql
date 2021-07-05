select name, id, type, resource_group, subscription_id
from azure.azure_data_factory_dataset
where name = '{{resourceName}}' and resource_group = '{{resourceName}}' and factory_name = '{{resourceName}}'