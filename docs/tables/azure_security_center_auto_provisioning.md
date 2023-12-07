---
title: "Steampipe Table: azure_security_center_auto_provisioning - Query Azure Security Center Auto Provisioning Settings using SQL"
description: "Allows users to query Azure Security Center Auto Provisioning Settings, providing insights into the automatic deployment of security services and controls."
---

# Table: azure_security_center_auto_provisioning - Query Azure Security Center Auto Provisioning Settings using SQL

Azure Security Center Auto Provisioning is a feature within Microsoft Azure that allows for the automatic deployment of security services and controls. It plays a crucial role in ensuring that the necessary security services are in place across all Azure resources, making it easier to maintain and monitor the security posture of your Azure environment. It promotes consistency and reduces the chance of misconfiguration or oversight in security controls deployment.

## Table Usage Guide

The `azure_security_center_auto_provisioning` table provides insights into the automatic deployment of security services and controls within Azure Security Center. As a Security or DevOps engineer, explore the details of auto provisioning settings through this table, including the target resource type and auto provisioning status. Utilize it to maintain optimal and consistent security posture across your Azure resources, and to ensure that all necessary security services are automatically deployed as needed.

## Examples

### Basic info
Determine the areas in which automatic provisioning is enabled in your Azure Security Center to enhance your security posture and reduce manual configuration efforts.

```sql+postgres
select
  id,
  name,
  type,
  auto_provision
from
  azure_security_center_auto_provisioning;
```

```sql+sqlite
select
  id,
  name,
  type,
  auto_provision
from
  azure_security_center_auto_provisioning;
```

### List subscriptions that have automatic provisioning of VM monitoring agent enabled
Discover the segments that have automatic virtual machine monitoring agent provisioning enabled. This can be beneficial to assess the elements within your system that are being automatically monitored, ensuring system performance and security.

```sql+postgres
select
  id,
  name,
  type,
  auto_provision
from
  azure_security_center_auto_provisioning
where
  auto_provision = 'On';
```

```sql+sqlite
select
  id,
  name,
  type,
  auto_provision
from
  azure_security_center_auto_provisioning
where
  auto_provision = 'On';
```