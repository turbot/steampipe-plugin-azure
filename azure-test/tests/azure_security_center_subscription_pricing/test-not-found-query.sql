select name, akas, title
from azure.azure_security_center_subscription_pricing
where name = 'dummy-{{ output.resource_name.value }}';
