select object_id, object_type, display_name, mail_enabled, security_enabled
from azure.azure_ad_group
where object_id = '{{ output.object_id.value }}'
