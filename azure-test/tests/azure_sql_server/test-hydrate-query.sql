select name, server_audit_policy, server_security_alert_policy, server_azure_ad_administrator, server_vulnerability_assessment, firewall_rules, encryption_protector
from azure.azure_sql_server
where name = '{{ resourceName }}' and resource_group = '{{ resourceName }}';