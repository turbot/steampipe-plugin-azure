select name, id, type
from azure.azure_batch_account
where name = '{{ resourceName }}';