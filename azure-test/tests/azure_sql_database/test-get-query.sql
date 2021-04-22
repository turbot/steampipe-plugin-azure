select name, id, server_name, status, type, containment_state, default_secondary_location, earliest_restore_date, edition, elastic_pool_name, location,max_size_bytes, zone_redundant, requested_service_objective_name, service_level_objective,transparent_data_encryption, title, tags, akas, region, resource_group, subscription_id
from
  azure.azure_sql_database
where
  name = '{{ resourceName }}'
  and resource_group = '{{ resourceName }}';