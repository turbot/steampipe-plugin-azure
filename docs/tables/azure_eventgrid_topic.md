# Table: azure_eventgrid_topic

The event grid topic provides an endpoint where the source sends events. The publisher creates the event grid topic, and decides whether an event source needs one topic or more than one topic. A topic is used for a collection of related events. To respond to certain types of events, subscribers decide which topics to subscribe to.

## Examples

### Basic info

```sql
select
  name,
  id,
  type,
  provisioning_state
from
  azure_eventgrid_topic;
```

### List domains not configured with private endpoint connections

```sql
select
  name,
  id,
  type,
  private_endpoint_connections
from
  azure_eventgrid_topic
where
  private_endpoint_connections is null;
```

### List domains with local authentication disabled

```sql
select
  name,
  id,
  type,
  disable_local_auth
from
  azure_eventgrid_topic
where
  disable_local_auth;
```
