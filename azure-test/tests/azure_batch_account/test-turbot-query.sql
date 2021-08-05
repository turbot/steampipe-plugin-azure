select name, tags, title, akas
from azure.azure_batch_account
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';