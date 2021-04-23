select id, name, identity
from azure.azure_app_service_web_app
where akas::text = '["{{output.resource_aka.value}}", "{{output.resource_aka_lower.value}}"]';