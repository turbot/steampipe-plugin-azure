select subscription_id, display_name
from azure.azure_subscription
where subscription_id = '{{ output.subscription_id.value }}'
