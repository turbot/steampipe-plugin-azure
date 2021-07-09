select name, title, akas
from azure.azure_data_factory_dataset
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';