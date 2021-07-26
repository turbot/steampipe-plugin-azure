select name, id, region, type, is_auto_inflate_enabled, kafka_enabled, resource_group
from azure.azure_eventhub_namespace
where name = '{{resourceName}}' and resource_group = '{{resourceName}}';
