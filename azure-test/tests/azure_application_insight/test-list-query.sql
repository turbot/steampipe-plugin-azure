select
  name,
  lower(id) as id
from
  azure.azure_application_insight
where 
  name = '{{resourceName}}';

