select object_id, user_principal_name, display_name, object_type, user_type, account_enabled, mail_nickname, given_name
from azure.azure_ad_user
where object_id = '{{ output.object_id.value }}'
