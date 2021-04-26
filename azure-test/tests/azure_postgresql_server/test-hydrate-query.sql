select 
  name, 
  rule ->> 'Name' as rule_name,
  rule ->> 'Type' as rule_type,
  rule -> 'FirewallRuleProperties' ->> 'endIpAddress' as end_ip_address,
  rule -> 'FirewallRuleProperties' ->> 'startIpAddress' as start_ip_address, 
  configurations ->> 'Name' as configuration_name,
  configurations -> 'ConfigurationProperties' ->> 'value' as configuration_value, 
  server_admin ->> 'Name' as server_admin_login_name
from 
  azure_postgresql_server,
  jsonb_array_elements(server_configurations) as configurations,
  jsonb_array_elements(firewall_rules) as rule,
  jsonb_array_elements(server_administrators) as server_admin
where 
  name = '{{ resourceName }}' 
  and resource_group = '{{ resourceName }}'
  and configurations ->> 'Name' = 'log_checkpoints';