select name, id, kind, region, client_affinity_enabled, enabled, https_only, reserved, resource_group, identity_type
from azure.azure_app_service_web_app
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
