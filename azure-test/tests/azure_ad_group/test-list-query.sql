select object_id, display_name
from azure.azure_ad_group
where object_id = '{{ output.object_id.value }}'
