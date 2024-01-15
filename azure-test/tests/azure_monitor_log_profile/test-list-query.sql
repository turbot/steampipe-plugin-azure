select
  name,
  id
from
  azure.azure_monitor_log_profile
where
  id = '{{ output.resource_id_lower.value }}';