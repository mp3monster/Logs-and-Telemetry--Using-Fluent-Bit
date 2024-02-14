The following details suggested enhancements that could be applied to the plugin to make it more enterprise-capable.

### Functional Behavior

In the functional space, then we could improve in the following ways:

##### Datable column Value as the event timestamp

Rather than apply the timestamp for the record based on when we ingest the record, it would be good if there was the option to set the timestamp to match one of the columns of the table. So that the timestamp reflected when the event actually happened.

### Additional DB Drivers and performance

Currently, the implementation only makes use of the Postgres database driver. The ideal goal is that the plugin can be configured to work with multiple DB server types as it sticks to using ANSI SQL.

Therefore, the dependencies would need updating, and the logic of creating the connection would updated to use the relevant driver implementation.

To minimize complexities with the drivers, we would recommend adopting pure Go drivers (some drivers map the invocations onto a C implementation).

We would recommend looking at:
- Oracle
- MS SQL Server

There is a good list of available drivers [here](https://zchee.github.io/golang-wiki/SQLDrivers/).

##### MySQL Driver Issue

The MySQL drive that has been used is a pure Go implementation. However we have found that if the query doesn't yield any rows, it actually throws a nil pointer/memory error which triggers a panic. To mitigate this we have implemented a *select count(* * ) step. This issue needs to raised and addressed with the driver.

##### Query Optimization

The Query implementation allows for a multiple row result. This is against the possibility we can cache the result and play lines back  once the context issue is addressed  for the input plugin. This will make the performance more efficient rather than querying once per row. The alternative to this is to modify the code to use Go's single record query implementation.

Update Optimization

Rather than inserting each record as it is read from the output plugin callback - if there are multiple records, we then insert them in a single transaction, being more efficient with the database.

### Caching of DB connections

Between invocations, the DB driver structure is not cached, and we reconstruct a new connection. As creating connections can be inefficient, caching the connection between invocations. There are some additional complexities that would need to be considered:

- Impact of the cached object holding the actual connection that crosses the C/Go layer.
- Addressing the possibility that the connection may be timed out by the server between calls from the Plugin. Therefore resilience and re-establishing the connection would need to be incorporated into each callback.
- One mitigation here is to pre-emptively refresh a connection based on some sort of schedule.

### Handling 1st time starts

When the plugin starts for the 1st time, it is possible that there is a lot of data. We need to manage the process of pulling all these records in without saturating Fluent Bit and ensure we don't lose track as a result of a Fluent Bit restart.  This could look a bit like the Tail plugin feature.

### Sourcing Password rather than in configuration

Currently, the database credentials are passed through from the configuration of the pipeline.  It would be particularly good if we could retrieve the credentials via other mechanisms, such as retrieving them directly from a credentials repository such as Keycloak.

### Unit Testing

The testing for the book has been manual, rather than proper unit tests with automation. Along with the generation of godoc.
