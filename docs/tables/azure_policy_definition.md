---
title: "Steampipe Table: azure_policy_definition - Query Azure Policy Definitions using SQL"
description: "Allows users to query Azure Policy Definitions, specifically the details of policy definitions within Azure Policy, providing insights into policy details and compliance information."
---

# Table: azure_policy_definition - Query Azure Policy Definitions using SQL

Azure Policy is a service in Azure that you use to create, assign, and manage policies. These policies enforce different rules and effects over your resources, so those resources stay compliant with your corporate standards and service level agreements. Azure Policy does this by running evaluations of your resources and scanning for those not compliant with the policies you have created.

## Table Usage Guide

The `azure_policy_definition` table provides insights into policy definitions within Azure Policy. As a security analyst, explore policy-specific details through this table, including policy rules, effects, and associated metadata. Utilize it to uncover information about policies, such as those with specific effects, the relationships between policies, and the verification of policy rules.

## Examples

### Basic info
Explore the policies defined within your Azure environment to better understand their purpose and type. This can be beneficial to gain insights into your current security configurations and to identify areas for potential improvement.

```sql+postgres
select
  id,
  name,
  display_name,
  type,
  jsonb_pretty(policy_rule) as policy_rule
from
  azure_policy_definition;
```

```sql+sqlite
select
  id,
  name,
  display_name,
  type,
  policy_rule
from
  azure_policy_definition;
```

### Get the policy definition by display name
Determine the specifics of a policy definition based on its display name. This is particularly useful in scenarios where you need to understand the details of a policy without having to navigate through multiple layers of your Azure policy definitions.

```sql+postgres
select
  id,
  name,
  display_name,
  type,
  jsonb_pretty(policy_rule) as policy_rule
from
  azure_policy_definition
where
  display_name = 'Private endpoint connections on Batch accounts should be enabled';
```

```sql+sqlite
select
  id,
  name,
  display_name,
  type,
  policy_rule
from
  azure_policy_definition
where
  display_name = 'Private endpoint connections on Batch accounts should be enabled';
```