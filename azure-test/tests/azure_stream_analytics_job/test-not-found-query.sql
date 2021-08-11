select name, tags, title, akas
from azure.azure_stream_analytics_job
where name = 'dummy-{{ resourceName }}' and resource_group = '{{ resourceName }}';