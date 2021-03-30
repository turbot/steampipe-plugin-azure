select name, type, version, location, administrator_login, backup_retention_days, fully_qualified_domain_name, geo_redundant_backup, infrastructure_encryption, minimal_tls_version, public_network_access, sku_name, sku_capacity, sku_family, sku_tier, ssl_enforcement, storage_profile_storage_auto_grow, resource_group, storage_profile_storage_mb, region, subscription_id
from azure.azure_mysql_server
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
