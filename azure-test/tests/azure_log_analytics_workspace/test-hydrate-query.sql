select
  name,
  to_jsonb(string_to_array(STRING_AGG(lower(x), ','), ',')) as akas
from
  azure.azure_log_analytics_workspace,
  jsonb_array_elements_text(akas) x
where 
  name = '{{resourceName}}' 
  and resource_group = '{{resourceName}}'
group by 
  name;