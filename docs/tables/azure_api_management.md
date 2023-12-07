---
title: "Steampipe Table: azure_api_management - Query Azure API Management Services using SQL"
description: "Allows users to query Azure API Management Services, specifically providing insights into the management of APIs for both on-premises and cloud solutions."
---

# Table: azure_api_management - Query Azure API Management Services using SQL

Azure API Management is a solution that allows organizations to publish, manage, secure, and analyze their APIs in a unified way. It provides the core competencies to ensure a successful API program through developer engagement, business insights, analytics, security, and protection. Azure API Management enables you to create an API gateway and developer portal in minutes.

## Table Usage Guide

The `azure_api_management` table provides insights into API management services within Azure. As a DevOps engineer, leverage this table to explore details about your API management services, including their configurations, locations, and associated resources. Utilize it to manage and secure your APIs, monitor their performance, and understand their usage patterns.

## Examples

### Public and private IP address info of each API management
Determine the areas in which each API management system operates by understanding their public and private IP addresses. This aids in assessing network accessibility and identifying potential security concerns.

```sql+postgres
select
  name,
  public_ip_addresses,
  private_ip_addresses
from
  azure_api_management;
```

```sql+sqlite
select
  name,
  public_ip_addresses,
  private_ip_addresses
from
  azure_api_management;
```


### API management publisher info
Explore the publisher details associated with your Azure API management to maintain effective communication and ensure smooth operations. This allows you to identify who is in charge of specific APIs, facilitating efficient management and collaboration.

```sql+postgres
select
  name,
  publisher_name,
  publisher_email
from
  azure_api_management;
```

```sql+sqlite
select
  name,
  publisher_name,
  publisher_email
from
  azure_api_management;
```


### List of premium API managements and their computing capacity
Identify instances where premium API management services are being used and assess their computing capacity. This can be useful in evaluating your resource allocation and optimizing your API management strategy.

```sql+postgres
select
  name,
  sku_name,
  sku_capacity
from
  azure_api_management
where
  sku_name = 'Premium';
```

```sql+sqlite
select
  name,
  sku_name,
  sku_capacity
from
  azure_api_management
where
  sku_name = 'Premium';
```


### List of API management without application tag key
Determine the areas in which API management in Azure lacks an 'application' tag. This could be useful for managing and organizing your resources, as well as ensuring compliance with tagging policies.

```sql+postgres
select
  name,
  tags
from
  azure_api_management
where
  not tags :: JSONB ? 'application';
```

```sql+sqlite
select
  name,
  tags
from
  azure_api_management
where
  json_extract(tags, '$.application') is null;
```