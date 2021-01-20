select name, id, region, type, kind, consistency_policy_max_interval, consistency_policy_max_staleness_prefix, default_consistency_level, database_account_offer_type, disable_key_based_metadata_write_access, document_endpoint, enable_automatic_failover, enable_multiple_write_locations, is_virtual_network_filter_enabled, resource_group
from azure.azure_cosmosdb_account
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
