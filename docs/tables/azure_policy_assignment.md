---
title: "Steampipe Table: azure_policy_assignment - Query Azure Policy Assignments using SQL"
description: "Allows users to query Policy Assignments in Azure, specifically the policy assignment details, providing insights into compliance status and policy enforcement."
---

# Table: azure_policy_assignment - Query Azure Policy Assignments using SQL

A Policy Assignment in Azure is a security tool that enables operators to apply a policy definition to a resource or a set of resources. The assignment is the process of binding a policy definition to a specific scope. This scope could range from a management group to a resource group.

## Table Usage Guide

The `azure_policy_assignment` table provides insights into Policy Assignments within Azure Policy. As a Security Analyst, explore specific details through this table, including policy definitions, scopes, and compliance statuses. Utilize it to uncover information about policy assignments, such as those associated with specific resources, the scope of these assignments, and their compliance status.

## Examples

### Basic info
Explore the policies assigned within your Azure environment to ensure adherence to your organization's governance and compliance requirements. This can help identify any instances where policies may not be correctly applied, potentially exposing your environment to risks.

```sql+postgres
select
  id,
  policy_definition_id,
  name,
  type
from
  azure_policy_assignment;
```

```sql+sqlite
select
  id,
  policy_definition_id,
  name,
  type
from
  azure_policy_assignment;
```

### Get SQL auditing and threat detection monitoring status for the subscription
Explore the status of SQL auditing and threat detection monitoring for your subscription. This query helps you assess whether these important security measures are active, promoting better risk management and data protection.

```sql+postgres
select
  id,
  policy_definition_id,
  display_name,
  parameters -> 'sqlAuditingMonitoringEffect' -> 'value' as sqlAuditingMonitoringEffect
from
  azure_policy_assignment;
```

```sql+sqlite
select
  id,
  policy_definition_id,
  display_name,
  json_extract(json_extract(parameters, '$.sqlAuditingMonitoringEffect'), '$.value') as sqlAuditingMonitoringEffect
from
  azure_policy_assignment;
```

### Get SQL encryption monitoring status for the subscription
Explore the status of SQL encryption monitoring for your subscription. This can help in maintaining the security of your data by keeping an eye on the encryption status.

```sql+postgres
select
  id,
  policy_definition_id,
  display_name,
  parameters -> 'sqlEncryptionMonitoringEffect' -> 'value' as sqlEncryptionMonitoringEffect
from
  azure_policy_assignment;
```

```sql+sqlite
select
  id,
  policy_definition_id,
  display_name,
  json_extract(json_extract(parameters, '$.sqlEncryptionMonitoringEffect'), '$.value') as sqlEncryptionMonitoringEffect
from
  azure_policy_assignment;
```