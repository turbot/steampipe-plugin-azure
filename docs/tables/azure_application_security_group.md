---
title: "Steampipe Table: azure_application_security_group - Query Azure Application Security Groups using SQL"
description: "Allows users to query Azure Application Security Groups, providing insights into security configuration and potential network vulnerabilities."
---

# Table: azure_application_security_group - Query Azure Application Security Groups using SQL

An Azure Application Security Group is a logical representation of an application in Azure. It allows for the grouping of servers based on applications for security and isolation of network traffic. This provides a more natural way to apply and manage security policies based on applications rather than explicit IP addresses or subnets.

## Table Usage Guide

The `azure_application_security_group` table provides insights into Application Security Groups within Azure. As a security analyst, explore application-specific details through this table, including security configurations, associated network interfaces, and potential vulnerabilities. Utilize it to uncover information about applications, such as those with weak security settings, the relationships between applications and network interfaces, and the verification of security policies.

## Examples

### Basic info
Explore which applications are grouped together in Azure, and determine the areas in which these groups are deployed. This can aid in understanding the organization and distribution of your applications across different regions.

```sql+postgres
select
  name,
  region,
  resource_group
from
  azure_application_security_group;
```

```sql+sqlite
select
  name,
  region,
  resource_group
from
  azure_application_security_group;
```


### List of application security group without application tag key
Identify instances where Azure application security groups lack the 'application' tag key. This can help streamline organization and management of security groups.

```sql+postgres
select
  name,
  tags
from
  azure_application_security_group
where
  not tags :: JSONB ? 'application';
```

```sql+sqlite
select
  name,
  tags
from
  azure_application_security_group
where
  json_extract(tags, '$.application') is null;
```