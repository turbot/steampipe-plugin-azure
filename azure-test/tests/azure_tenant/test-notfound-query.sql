select display_name, title
from azure.azure_tenant
where display_name = 'dummy-{{ resourceName }}';