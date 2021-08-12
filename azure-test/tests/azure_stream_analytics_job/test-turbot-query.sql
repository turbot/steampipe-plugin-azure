select name, title, akas
from azure.azure_stream_analytics_job
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';