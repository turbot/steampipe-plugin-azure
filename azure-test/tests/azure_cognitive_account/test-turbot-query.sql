select name, akas, tags, title
from azure.azure_cognitive_account
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
