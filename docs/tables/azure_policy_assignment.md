---
title: "Steampipe Table: azure_policy_assignment - Query Azure Policy Assignments using SQL"
description: "Allows users to query Azure Policy Assignments."
---

# Table: azure_policy_assignment - Query Azure Policy Assignments using SQL

Azure Policy is a service in Azure that you use to create, assign, and manage policies. These policies enforce different rules and effects over your resources, so those resources stay compliant with your corporate standards and service level agreements. Azure Policy meets this need by evaluating your resources for non-compliance with assigned policies.

## Table Usage Guide

The 'azure_policy_assignment' table provides insights into policy assignments within Azure Policy. As a DevOps engineer, explore policy-specific details through this table, including policy definition, scope, and associated metadata. Utilize it to uncover information about policy assignments, such as those with specific effects, the resource groups they are applied to, and the compliance state of your resources. The schema presents a range of attributes of the policy assignment for your analysis, like the assignment name, id, type, and associated parameters.

## Examples

### Basic info
Explore the specific policies applied within your Azure environment. This query can help you gain insights into policy assignments, which is beneficial for maintaining compliance and managing resources effectively.

```sql
select
  id,
  policy_definition_id,
  name,
  type
from
  azure_policy_assignment;
```

### Get SQL auditing and threat detection monitoring status for the subscription
Assess the status of SQL auditing and threat detection monitoring for a specific subscription. This can help improve your security measures by identifying areas that need attention or improvement.

```sql
select
  id,
  policy_definition_id,
  display_name,
  parameters -> 'sqlAuditingMonitoringEffect' -> 'value' as sqlAuditingMonitoringEffect
from
  azure_policy_assignment;
```

### Get SQL encryption monitoring status for the subscription
Explore the status of SQL encryption monitoring for your subscription. This allows you to assess the security measures in place and ensure that sensitive data is appropriately protected.

```sql
select
  id,
  policy_definition_id,
  display_name,
  parameters -> 'sqlEncryptionMonitoringEffect' -> 'value' as sqlEncryptionMonitoringEffect
from
  azure_policy_assignment;
```