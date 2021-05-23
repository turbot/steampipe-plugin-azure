select name, id, type
from azure.azure_security_center_subscription_pricing
where name = '{{ output.resource_name.value }}';
