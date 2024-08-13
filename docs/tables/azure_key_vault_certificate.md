---
title: "Steampipe Table: azure_key_vault_certificate - Query Azure Key Vault Certificates using SQL"
description: "Allows users to query Azure Key Vault Certificates, providing access to certificate details, including key ID, key properties, and secret properties."
---

# Table: azure_key_vault_certificate - Query Azure Key Vault Certificates using SQL

Azure Key Vault is a cloud service that provides a secure store for secrets, keys, and certificates. Key Vault certificates are managed digital certificates that aid in secure communication and identity verification, which are essential in various IT and cloud scenarios.

## Table Usage Guide

The `azure_key_vault_certificate` table allows users to explore and manage certificates within Azure Key Vault. This table is especially useful for security engineers and cloud administrators who need to oversee the state, configuration, and properties of certificates stored in Azure Key Vault.

## Examples

### Basic info
Review the general status and details of your Azure Key Vault certificates. This query is fundamental for routine checks and ensuring that certificates are up-to-date and correctly enabled.

```sql+postgres
select
  name,
  vault_name,
  enabled,
  created,
  updated
from
  azure_key_vault_certificate;
```

```sql+sqlite
select
  name,
  vault_name,
  enabled,
  created,
  updated
from
  azure_key_vault_certificate;
```

### List disabled certificates
Identify certificates that are currently disabled in Azure Key Vault. This query helps maintain proper security measures and effective access control.

```sql+postgres
select
  name,
  vault_name,
  enabled
from
  azure_key_vault_certificate
where
  not enabled;
```

```sql+sqlite
select
  name,
  vault_name,
  enabled
from
  azure_key_vault_certificate
where
  not enabled;
```

### List certificates expiring in 10 days
Discover certificates within Azure Key Vault that are nearing their expiration date. This is crucial for proactive certificate renewal and avoiding potential security risks.

```sql+postgres
select
  name,
  enabled,
  not_before,
  created,
  expires
from
  azure_key_vault_certificate
where
  expires <= now() + interval '10 days';
```

```sql+sqlite
select
  name,
  enabled,
  not_before,
  created,
  expires
from
  azure_key_vault_certificate
where
  datetime(expires) <= datetime('now', '+10 days');
```

### Get key properties of certificates
Analyze the key properties of certificates, including their exportability, type, size, and reuse policies. This information is vital for understanding the security and operational characteristics of each certificate.

```sql+postgres
select
  name,
  id,
  key_properties ->> 'Exportable' as exportable,
  key_properties ->> 'KeyType' as key_type,
  key_properties ->> 'KeySize' as key_size,
  key_properties ->> 'ReuseKey' as reuse_key
from
  azure_key_vault_certificate;
```

```sql+sqlite
select
  name,
  id,
  json_extract(key_properties, '$.Exportable') as exportable,
  json_extract(key_properties, '$.KeyType') as key_type,
  json_extract(key_properties, '$.KeySize') as key_size,
  json_extract(key_properties, '$.ReuseKey') as reuse_key
from
  azure_key_vault_certificate;
```

### Get X509 properties of certificates
Examine the X509 properties of certificates, such as the subject, extended key usage (EKUs), alternative names, key usage, and validity. This query is crucial for detailed certificate analysis and compliance checks.

```sql+postgres
select
  name,
  id,
  x509_certificate_properties ->> 'Subject' as subject,
  x509_certificate_properties -> 'Ekus' as ekus,
  x509_certificate_properties -> 'SubjectAlternativeNames' as subject_alternative_names,
  x509_certificate_properties ->> 'KeyUsage' as key_usage,
  x509_certificate_properties ->> 'ValidityInMonths' as validity_in_months
from
  azure_key_vault_certificate;
```

```sql+sqlite
select
  name,
  id,
  json_extract(x509_certificate_properties, '$.Subject') as subject,
  json_extract(x509_certificate_properties, '$.Ekus') as ekus,
  json_extract(x509_certificate_properties, '$.SubjectAlternativeNames') as subject_alternative_names,
  json_extract(x509_certificate_properties, '$.KeyUsage') as key_usage,
  json_extract(x509_certificate_properties, '$.ValidityInMonths') as validity_in_months
from
  azure_key_vault_certificate;
```
