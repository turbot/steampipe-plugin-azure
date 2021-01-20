select object_id, object_type, akas, title
from azure.azure_ad_group
where object_id = '{{ output.object_id.value }}'
