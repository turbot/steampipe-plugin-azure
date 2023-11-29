---
title: "Steampipe Table: azure_security_center_auto_provisioning - Query Azure Security Center Auto Provisioning Settings using SQL"
description: "Allows users to query Azure Security Center Auto Provisioning Settings."
---

# Table: azure_security_center_auto_provisioning - Query Azure Security Center Auto Provisioning Settings using SQL

Azure Security Center is a unified infrastructure security management system that strengthens the security posture of your data centers, and provides advanced threat protection across your hybrid workloads in the cloud - whether they're in Azure or not. Auto Provisioning settings in Azure Security Center enable you to control if security solutions are automatically deployed and provisioned for new resources. This feature is designed to ensure that as new resources are deployed, they are automatically onboarded to the security solutions and policies you have defined.

## Table Usage Guide

The 'azure_security_center_auto_provisioning' table provides insights into the auto provisioning settings within Azure Security Center. As a security administrator, explore setting-specific details through this table, including the current auto provisioning status and target resource type. Utilize it to uncover information about the auto provisioning settings, such as those that are currently active and the resource types they are applied to. The schema presents a range of attributes of the auto provisioning settings for your analysis, like the auto provisioning setting id, provisioning status, and target resource type.

## Examples

### Basic info
Discover the segments that have automatic provisioning enabled in your Azure Security Center to better manage your security policies and configurations. This is useful for maintaining security standards and ensuring consistent configurations across your environment.

```sql
select
  id,
  name,
  type,
  auto_provision
from
  azure_security_center_auto_provisioning;
```

### List subscriptions that have automatic provisioning of VM monitoring agent enabled
Discover the subscriptions that have enabled automatic provisioning for their VM monitoring agent. This allows you to identify potential areas for increased security and efficiency.

```sql
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