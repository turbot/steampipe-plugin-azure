select object_id, account_enabled, display_name, object_type, app_role_assignment_required
from azure.azure_ad_service_principal
where object_id = 'test-service-principal12234455352412'
