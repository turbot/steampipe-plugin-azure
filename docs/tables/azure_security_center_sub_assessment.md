# Table: azure_security_center_sub_assessment

Azure security center sub assessment details for the subscription.

## Examples

### Basic info

```sql
select
  id,
  name,
  display_name,
  type,
  category
from
  azure_security_center_sub_assessment;
```

### List unhealthy sub assessment details

```sql
select
  name,
  type,
  category,
  status
from
  azure_security_center_sub_assessment
where
  status ->> 'Code' = 'Unhealthy';
```

### List all container registry vulnerabilities with corresponding remedies

```sql
select
  container_registry_vulnerability_properties,
  remediation,
  resource_details
from
  azure_security_center_sub_assessment
where
  container_registry_vulnerability_properties ->> 'AssessedResourceType' =  'ContainerRegistryVulnerability';
```

### List all server vulnerabilities with corresponding remedies

```sql
select
  server_vulnerability_properties,
  remediation,
  resource_details
from
  azure_security_center_sub_assessment
where
  server_vulnerability_properties ->> 'AssessedResourceType' =  'ServerVulnerabilityAssessment';
```

### List all sql server vulnerabilities with corresponding remedies

```sql
select
  sql_server_vulnerability_properties,
  remediation,
  resource_details
from
  azure_security_center_sub_assessment
where
  sql_server_vulnerability_properties ->> 'AssessedResourceType' =  'SqlServerVulnerability';
```