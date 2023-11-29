# Table: azure_ad_group

**Deprecated. Use [azuread_group](https://hub.steampipe.io/plugins/turbot/azuread/tables/azuread_group) instead.**

Azure Active Directory groups is used to manage access to your cloud-based apps, on-premises apps, and your resources.

## Examples

### Basic info
Determine the areas in which your Azure Active Directory groups are mail-enabled. This could be beneficial for managing group email communications and understanding which groups have specific email settings.

```sql
select
  display_name,
  object_id,
  mail,
  mail_enabled,
  mail_nickname
from
  azure_ad_group;
```


### List of AD groups where security is not enabled
Determine the areas in which Azure Active Directory groups are not security-enabled. This is crucial for identifying potential vulnerabilities and enhancing the security posture of your organization.

```sql
select
  display_name,
  object_id,
  security_enabled
from
  azure_ad_group
where
  not security_enabled;
```


### List of AD groups where mail is not enabled
Determine the areas in which Azure Active Directory groups have not enabled mail. This can be useful in identifying groups that may not be receiving important communications or updates.

```sql
select
  display_name,
  mail,
  mail_enabled
from
  azure_ad_group
where
  not mail_enabled;
```