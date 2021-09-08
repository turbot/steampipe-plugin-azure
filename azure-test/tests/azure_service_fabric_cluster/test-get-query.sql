select name, id, type, region, resource_group, subscription_id, reliability_level, upgrade_mode, vm_image, management_endpoint
from azure.azure_service_fabric_cluster
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';
