---
title: "Steampipe Table: azure_provider - Query Azure Resource Providers using SQL"
description: "Allows users to query Azure Resource Providers."
---

# Table: azure_provider - Query Azure Resource Providers using SQL

Azure Resource Providers are services that supply the resources you can deploy and manage through Resource Manager. Each resource provider offers operations for working with the resources that are deployed. Some common resource providers are Microsoft.Compute, which supplies the virtual machine resource, Microsoft.Storage, which supplies the storage account resource, and Microsoft.Web, which supplies resources related to web apps.

## Table Usage Guide

The 'azure_provider' table provides insights into Resource Providers within Microsoft Azure. As a DevOps engineer, explore provider-specific details through this table, including the provider's namespace, registration state, and resource types. Utilize it to uncover information about providers, such as those that are registered or unregistered, the resources they supply, and their capabilities. The schema presents a range of attributes of the Resource Provider for your analysis, like the provider ID, registration state, and resource types.

## Examples

### Basic info
Explore the registration status of your Azure provider to understand its operational state and ensure it's properly configured. This can be useful in maintaining the efficiency of your cloud infrastructure.

```sql
select
  id,
  namespace,
  registration_state
from
  azure_provider;
```

### List of azure providers which are not registered for use
Explore which Azure providers are not registered for use. This can be particularly useful in identifying potential gaps in your Azure services setup, helping to ensure all necessary providers are correctly registered and operational.

```sql
select
  namespace,
  registration_state
from
  azure_provider
where
  registration_state = 'NotRegistered';
```