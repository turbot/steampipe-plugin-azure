select name, id
from azure.azure_lb_probe
where name = 'dummy-test-azure-lb-probe' and resource_group = '{{resourceName}}'
