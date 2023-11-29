---
title: "Steampipe Table: azure_api_management - Query Azure API Management Services using SQL"
description: "Allows users to query Azure API Management Services."
---

# Table: azure_api_management - Query Azure API Management Services using SQL

Azure API Management is a fully managed service that helps customers publish, secure, transform, maintain, and monitor APIs. With Azure API Management, organizations can ensure that their APIs are always available and performing as expected, and that their valuable data is secure. The service also includes a developer portal to help onboard developers and foster a developer community.

## Table Usage Guide

The 'azure_api_management' table provides insights into API Management Services within Azure. As a DevOps engineer, explore service-specific details through this table, including API names, locations, and associated metadata. Utilize it to uncover information about services, such as those with specific SKUs, the regions they are deployed in, and the verification of their identities. The schema presents a range of attributes of the API Management Service for your analysis, like the service name, resource group, subscription ID, and associated tags.

## Examples

### Public and private IP address info of each API management
Gain insights into the public and private IP addresses associated with each API management system in your Azure environment. This allows for better network management and security monitoring.

```sql
select
  name,
  public_ip_addresses,
  private_ip_addresses
from
  azure_api_management;
```


### API management publisher info
Gain insights into the publishers of your Azure API management service, including their names and contact emails, to facilitate effective communication and management.

```sql
select
  name,
  publisher_name,
  publisher_email
from
  azure_api_management;
```


### List of premium API managements and their computing capacity
Identify premium Azure API management services and their computing capacities. This is useful for assessing your organization's API management capabilities and planning for potential upgrades or expansions.

```sql
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
Identify instances where API management in Azure is missing the 'application' tag. This can aid in pinpointing areas where tagging conventions may not have been followed, helping to improve resource management and compliance.

```sql
select
  name,
  tags
from
  azure_api_management
where
  not tags :: JSONB ? 'application';
```