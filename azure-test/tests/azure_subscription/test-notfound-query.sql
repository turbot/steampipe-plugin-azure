select display_name, title
from azure.azure_subscription
where display_name = 'dummy-{{ resourceName }}';