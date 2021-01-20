select name, id, type, disk_size_bytes, disk_size_gb, incremental, create_option, sku_name, region, resource_group, subscription_id
from azure.azure_compute_snapshot
where name = '{{resourceName}}snapshot' and resource_group = '{{resourceName}}'