select name, display_name, tenant_id
from azure.azure_tenant
where tenant_id = '{{ output.tenant_id.value }}'
