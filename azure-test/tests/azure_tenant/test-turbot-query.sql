select name, title
from azure.azure_tenant
where tenant_id = '{{ output.tenant_id.value }}';