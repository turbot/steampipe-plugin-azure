# Table:  azure_ad_group

Azure Active Directory groups is used to manage access to your cloud-based apps, on-premises apps, and your resources.

## Examples

### AD groups' users basic info

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