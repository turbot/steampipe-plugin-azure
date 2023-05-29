select
  name,
  kind,
  retention_in_days,
  region,
  resource_group
from
  azure.azure_application_insight
where 
  name = '{{resourceName}}' 
  and resource_group = '{{resourceName}}';
