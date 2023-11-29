---
title: "Steampipe Table: azure_security_center_sub_assessment - Query Azure Security Center Sub-Assessments using SQL"
description: "Allows users to query Azure Security Center Sub-Assessments"
---

# Table: azure_security_center_sub_assessment - Query Azure Security Center Sub-Assessments using SQL

Azure Security Center is a unified infrastructure security management system by Microsoft Azure that improves the security posture of your data centers. It provides advanced threat protection across your hybrid workloads in the cloud, whether they're in Azure or not. As part of this service, Sub-Assessments provide detailed security recommendations and potential vulnerabilities within your resources.

## Table Usage Guide

The 'azure_security_center_sub_assessment' table provides insights into Sub-Assessments within Azure Security Center. As a security professional, you can explore detailed security recommendations and potential vulnerabilities for your resources through this table. Utilize it to uncover information about each sub-assessment, such as its status, severity, and associated metadata. The schema presents a range of attributes of the sub-assessment for your analysis, like the resource ID, resource type, and associated recommendations. This can be particularly useful in identifying and mitigating potential security risks in your Azure environment.

## Examples

### Basic info
Explore which security assessments in your Azure Security Center have specific characteristics. This can help you identify potential risk areas and understand the security posture of your resources.

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
Explore which sub-assessments in Azure Security Center are marked as unhealthy. This can help you identify areas of your Azure environment that may require immediate attention or remediation.

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
Explore potential security vulnerabilities within your container registry and understand the corresponding solutions. This is useful for maintaining the security of your applications by identifying and addressing potential threats.

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
Explore server vulnerabilities and their corresponding remedies within the Azure Security Center. This is useful for identifying potential security issues and understanding how to address them.

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
Discover the segments that contain vulnerabilities in your SQL server and understand the corresponding remedies. This can help in ensuring your server's security by addressing these vulnerabilities promptly.

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