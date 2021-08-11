select name, id, type, region, resource_group, subscription_id
from azure.azure_stream_analytics_job
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';