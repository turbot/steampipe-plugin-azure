select object_id, object_type, user_type, akas, title
from azure.azure_ad_user
where object_id = '{{ output.object_id.value }}'
