select name, id, type, region
from azure.azure_web_application_firewall_policy
where name = 'dummy-test{{ resourceName }}' and resource_group = '{{ resourceName }}';
