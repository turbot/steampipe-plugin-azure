---
title: "Steampipe Table: azure_security_center_sub_assessment - Query Azure Security Center Sub-Assessments using SQL"
description: "Allows users to query Azure Security Center Sub-Assessments, providing detailed security findings for each resource."
---

# Table: azure_security_center_sub_assessment - Query Azure Security Center Sub-Assessments using SQL

Azure Security Center is a unified infrastructure security management system that strengthens the security posture of your data centers, and provides advanced threat protection across your hybrid workloads in the cloud - whether they're in Azure or not. It provides security management and threat protection across your hybrid cloud workloads. It allows you to prevent, detect, and respond to threats with increased visibility.

## Table Usage Guide

The `azure_security_center_sub_assessment` table provides insights into the detailed security findings for each resource within Azure Security Center. As a security engineer, you can explore specific security assessment details through this table, including severity, status, and associated metadata. Utilize it to uncover information about security vulnerabilities and the remediation steps for each resource.

## Examples

### Basic info
Explore the various sub-assessments within Azure's Security Center. This allows you to understand and categorize different security elements by their unique identifiers, names, display names, types, and categories, providing a comprehensive overview of your security landscape.

```sql+postgres
select
  id,
  name,
  display_name,
  type,
  category
from
  azure_security_center_sub_assessment;
```

```sql+sqlite
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
Determine the areas in which security aspects are marked as unhealthy in Azure Security Center. This provides a way to identify instances where security measures may need to be improved or updated.

```sql+postgres
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

```sql+sqlite
select
  name,
  type,
  category,
  status
from
  azure_security_center_sub_assessment
where
  json_extract(status, '$.Code') = 'Unhealthy';
```

### List all container registry vulnerabilities with corresponding remedies
Identify potential security vulnerabilities in your container registry and uncover the specific remediation steps to mitigate them. This is crucial for maintaining robust security practices and ensuring system integrity.

```sql+postgres
select
  container_registry_vulnerability_properties,
  remediation,
  resource_details
from
  azure_security_center_sub_assessment
where
  container_registry_vulnerability_properties ->> 'AssessedResourceType' =  'ContainerRegistryVulnerability';
```

```sql+sqlite
select
  container_registry_vulnerability_properties,
  remediation,
  resource_details
from
  azure_security_center_sub_assessment
where
  json_extract(container_registry_vulnerability_properties, '$.AssessedResourceType') =  'ContainerRegistryVulnerability';
```

### List all server vulnerabilities with corresponding remedies
Determine the areas in which server vulnerabilities exist and discover the corresponding remedies. This is beneficial for maintaining server security and ensuring prompt remediation of any identified vulnerabilities.

```sql+postgres
select
  server_vulnerability_properties,
  remediation,
  resource_details
from
  azure_security_center_sub_assessment
where
  server_vulnerability_properties ->> 'AssessedResourceType' =  'ServerVulnerabilityAssessment';
```

```sql+sqlite
select
  server_vulnerability_properties,
  remediation,
  resource_details
from
  azure_security_center_sub_assessment
where
  json_extract(server_vulnerability_properties, '$.AssessedResourceType') =  'ServerVulnerabilityAssessment';
```

### List all sql server vulnerabilities with corresponding remedies
Explore vulnerabilities in your SQL server and ascertain appropriate remedies. This query is useful for maintaining security and addressing potential risks in your SQL server environment.

```sql+postgres
select
  sql_server_vulnerability_properties,
  remediation,
  resource_details
from
  azure_security_center_sub_assessment
where
  sql_server_vulnerability_properties ->> 'AssessedResourceType' =  'SqlServerVulnerability';
```

```sql+sqlite
select
  sql_server_vulnerability_properties,
  remediation,
  resource_details
from
  azure_security_center_sub_assessment
where
  json_extract(sql_server_vulnerability_properties, '$.AssessedResourceType') =  'SqlServerVulnerability';
```