select name, id, type, version, location, version, administrator_login, fully_qualified_domain_name, minimal_tls_version, public_network_access, sku_family, sku_name, sku_tier, sku_size, ssl_enforcement, backup_retention_days, geo_redundant_backup, resource_group, region, subscription_id
from azure.azure_postgresql_server
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
