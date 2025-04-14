---
title: "Steampipe Table: azure_provider - Query Azure Providers using SQL"
description: "Allows users to query Azure Providers, specifically the registration status, namespace, and other related properties, providing insights into the Azure resource providers' registration state."
folder: "Resource Manager"
---

# Table: azure_provider - Query Azure Providers using SQL

An Azure Provider is a service that supplies the resources you can deploy and manage through Resource Manager. Each resource provider offers operations for working with the resources that are deployed. Some common resource providers are Microsoft.Compute, which supplies the virtual machine resource, Microsoft.Storage, which supplies the storage account resource, and Microsoft.Web, which supplies resources related to web apps.

## Table Usage Guide

The `azure_provider` table provides insights into Azure providers within Microsoft Azure Resource Manager. As a DevOps engineer, explore provider-specific details through this table, including registration status, namespace, and other related properties. Utilize it to uncover information about providers, such as their registration state, the resources they supply, and their corresponding operations.

## Examples

### Basic info
Determine the areas in which your Azure provider is registered. This is useful for understanding your Azure resources and their distribution.

```sql+postgres
select
  id,
  namespace,
  registration_state
from
  azure_provider;
```

```sql+sqlite
select
  id,
  namespace,
  registration_state
from
  azure_provider;
```

### List of azure providers which are not registered for use
Identify the Azure providers that are not yet registered for use. This is useful to ensure all necessary providers are properly set up and ready for use in your Azure environment.

```sql+postgres
select
  namespace,
  registration_state
from
  azure_provider
where
  registration_state = 'NotRegistered';
```

```sql+sqlite
select
  namespace,
  registration_state
from
  azure_provider
where
  registration_state = 'NotRegistered';
```