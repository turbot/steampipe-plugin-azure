select name, id, type, region
from azure.azure_cognitive_account
where name = 'dummy-test{{ resourceName }}' and resource_group = '{{ resourceName }}';
