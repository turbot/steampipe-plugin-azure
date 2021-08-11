select name, tags, title, akas
from azure.azure_batch_account
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';