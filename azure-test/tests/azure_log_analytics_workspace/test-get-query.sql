select
  name,
  retention_in_days,
  region,
  resource_group
from
  azure.azure_log_analytics_workspace
where 
  name = '{{resourceName}}' 
  and resource_group = '{{resourceName}}';
