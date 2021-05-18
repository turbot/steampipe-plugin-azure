select display_name, title
from azure.azure_subscription
where display_name = '{{ output.current_subscription_display_name.value }}';