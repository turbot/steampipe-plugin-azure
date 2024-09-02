select name, id, type, region, resource_group, subscription_id
from azure.azure_web_application_firewall_policy
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
