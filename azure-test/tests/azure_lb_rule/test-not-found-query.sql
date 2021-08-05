select name, id
from azure.azure_lb_rule
where name = 'dummy-test-azure-lb-rule' and resource_group = '{{resourceName}}'
