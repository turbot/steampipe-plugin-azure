select name, id, type
from azure.azure_stream_analytics_job
where name = '{{ resourceName }}';