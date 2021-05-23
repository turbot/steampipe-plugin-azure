select id, name
from azure.azure_security_center_subscription_pricing
where id = '{{ output.resource_id.value }}'
