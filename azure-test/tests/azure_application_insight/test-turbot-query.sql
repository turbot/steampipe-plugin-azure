select
  name,
  to_jsonb(string_to_array(STRING_AGG(lower(x), ','), ',')) as akas,
  tags,
  title
from
  azure.azure_application_insight,
  jsonb_array_elements_text(akas) x
where 
  name = '{{resourceName}}' 
  and resource_group = '{{resourceName}}'
group by
  name, tags, title;