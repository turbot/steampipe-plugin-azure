---
title: "Steampipe Table: azure_cdn_frontdoor_custom_domain - Query Azure CDN Front Door Custom Domains using SQL"
description: "Allows users to query CDN Front Door Custom Domains in Azure, providing detailed information about each custom domain, including its host name, provisioning state, validation state, and TLS settings."
folder: "CDN"
---

# Table: azure_cdn_frontdoor_custom_domain - Query Azure CDN Front Door Custom Domains using SQL

An Azure CDN Front Door Custom Domain represents a user-owned hostname that is associated with an Azure Front Door profile. It enables custom branding for content delivered through Azure Front Door by mapping a custom DNS name to the Front Door endpoint.

## Table Usage Guide

The `azure_cdn_frontdoor_custom_domain` table provides insights into custom domains configured within Azure CDN Front Door profiles. As an Infrastructure or Security Engineer, use this table to audit domain validation states, TLS settings, and provisioning status across all Front Door custom domains in your subscription.

## Examples

### Basic custom domain information
Retrieve the name, host name, and provisioning state of all Front Door custom domains.

```sql+postgres
select
  name,
  host_name,
  provisioning_state,
  resource_group
from
  azure_cdn_frontdoor_custom_domain;
```

```sql+sqlite
select
  name,
  host_name,
  provisioning_state,
  resource_group
from
  azure_cdn_frontdoor_custom_domain;
```

### Custom domains by profile
List all custom domains grouped by their Front Door profile.

```sql+postgres
select
  profile_name,
  name,
  host_name,
  domain_validation_state
from
  azure_cdn_frontdoor_custom_domain
order by
  profile_name;
```

```sql+sqlite
select
  profile_name,
  name,
  host_name,
  domain_validation_state
from
  azure_cdn_frontdoor_custom_domain
order by
  profile_name;
```

### Custom domains with pending or failed validation
Identify custom domains that are not yet validated or have a validation failure.

```sql+postgres
select
  name,
  host_name,
  profile_name,
  domain_validation_state,
  resource_group
from
  azure_cdn_frontdoor_custom_domain
where
  domain_validation_state <> 'Approved';
```

```sql+sqlite
select
  name,
  host_name,
  profile_name,
  domain_validation_state,
  resource_group
from
  azure_cdn_frontdoor_custom_domain
where
  domain_validation_state <> 'Approved';
```

### Custom domains with TLS settings
Retrieve TLS configuration details for each custom domain.

```sql+postgres
select
  name,
  host_name,
  profile_name,
  tls_settings
from
  azure_cdn_frontdoor_custom_domain
where
  tls_settings is not null;
```

```sql+sqlite
select
  name,
  host_name,
  profile_name,
  tls_settings
from
  azure_cdn_frontdoor_custom_domain
where
  tls_settings is not null;
```

### Custom domains not yet fully deployed
Find custom domains whose deployment is still in progress or not succeeded.

```sql+postgres
select
  name,
  host_name,
  profile_name,
  deployment_status,
  provisioning_state
from
  azure_cdn_frontdoor_custom_domain
where
  deployment_status <> 'Running';
```

```sql+sqlite
select
  name,
  host_name,
  profile_name,
  deployment_status,
  provisioning_state
from
  azure_cdn_frontdoor_custom_domain
where
  deployment_status <> 'Running';
```
