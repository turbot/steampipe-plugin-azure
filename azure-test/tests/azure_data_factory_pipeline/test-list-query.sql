select name, id, resource_group, subscription_id
from azure.azure_data_factory_pipeline
where name = '{{ resourceName }}';