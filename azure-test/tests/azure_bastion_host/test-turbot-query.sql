select 
  name,
  to_jsonb(string_to_array(STRING_AGG(lower(x), ','), ',')) as akas,
  title
from 
  azure.azure_bastion_host,
  jsonb_array_elements_text(akas) x
where
  name = '{{resourceName}}' 
  and resource_group = '{{resourceName}}'
group by
  name, title;