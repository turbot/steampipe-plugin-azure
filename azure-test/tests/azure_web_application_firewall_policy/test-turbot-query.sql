select name, akas, title
from azure.azure_web_application_firewall_policy
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
