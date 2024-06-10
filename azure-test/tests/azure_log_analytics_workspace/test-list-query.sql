select
  name,
  lower(id) as id
from
  azure.azure_log_analytics_workspace
where 
  name = '{{resourceName}}';

