---
title: "Steampipe Table: azure_policy_definition - Query Azure Policy Definitions using SQL"
description: "Allows users to query Azure Policy Definitions to gain insights into the policy definitions within Azure Policy service. The table provides details such as policy definition ID, name, type, mode, and metadata."
---

# Table: azure_policy_definition - Query Azure Policy Definitions using SQL

Azure Policy is a service in Azure that you use to create, assign and, manage policies. These policies enforce different rules and effects over your resources, so those resources stay compliant with your corporate standards and service level agreements. Azure Policy meets this need by evaluating your resources for non-compliance with assigned policies.

## Table Usage Guide

The 'azure_policy_definition' table provides insights into policy definitions within Azure Policy service. As a security engineer, explore policy-specific details through this table, including policy definition ID, name, type, mode, and metadata. Utilize it to uncover information about policies, such as their compliance status, the specific rules they enforce, and their effects on your resources. The schema presents a range of attributes of the policy definition for your analysis, like the policy definition ID, name, type, mode, and associated metadata.

## Examples

### Basic info
Explore policy definitions within your Azure environment to gain insights into their specific details such as ID, name, and type. This can be particularly useful for understanding and managing the rules and regulations applied to your resources.

```sql
select
  id,
  name,
  display_name,
  type,
  jsonb_pretty(policy_rule) as policy_rule
from
  azure_policy_definition;
```

### Get the policy definition by display name
Explore the policy definitions by their display names to understand their rules and types. This is particularly useful for managing and enforcing specific policies, such as enabling private endpoint connections on Batch accounts.

```sql
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