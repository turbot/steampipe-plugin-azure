select name, tags, title, akas
from azure.azure_data_factory
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';