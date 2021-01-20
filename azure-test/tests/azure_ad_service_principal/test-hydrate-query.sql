select object_id, object_type
from azure.azure_ad_service_principal
where object_id = '{{ output.object_id.value }}'
