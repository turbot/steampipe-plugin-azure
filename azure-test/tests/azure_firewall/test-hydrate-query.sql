select name, sku_name, sku_tier, threat_intel_mode, akas, tags, title
from azure.azure_firewall
where name = '{{resourceName}}' and resource_group = '{{resourceName}}'
