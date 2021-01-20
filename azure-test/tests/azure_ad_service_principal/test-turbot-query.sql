select object_id, akas, title
from azure.azure_ad_service_principal
where object_id = '{{ output.object_id.value }}'
