select name, id, type, sku_name, sku_tier, disk_size_gb, disk_size_bytes, encryption_settings_collection_enabled, encryption_type, region, resource_group, subscription_id
from azure.azure_compute_disk
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'