The following details suggested enhancements that could be applied to the plugin to make it more enterprise-capable.



### Additional DB Drivers

Currently, the implementation only makes use of the Postgres database driver. The ideal goal is that the plugin can be configured to work with multiple DB server types as it sticks to using ANSI SQL.

Therefore, the dependencies would need updating, and the logic of creating the connection would updated to use the relevant driver implementation.

To minimize complexities with the drivers, we would recommend adopting pure Go drivers (some drivers map the invocations onto a C implementation).

We would recommend looking at:
- Oracle
- MS SQL Server

There is a good list of available drivers [here](https://zchee.github.io/golang-wiki/SQLDrivers/).

### Caching of DB connections

Between invocations, the DB driver structure is not cached, and we reconstruct a new connection. As creating connections can be inefficient, caching the connection between invocations. There are some additional complexities that would need to be considered:

- Impact of the cached object holding the actual connection that crosses the C/Go layer.
- Addressing the possibility that the connection may be timed out by the server between calls from the Plugin. Therefore resilience and re-establishing the connection would need to be incorporated into each callback.
- One mitigation here is to pre-emptively refresh a connection based on some sort of schedule.

### Handling 1st time starts

When the plugin starts for the 1st time, it is possible that there is a lot of data. We need to manage the process of pulling all these records in without saturating Fluent Bit and ensure we don't lose track as a result of a Fluent Bit restart.  This could look a bit like the Tail plugin feature.

### Sourcing Password rather than in configuration

Currently, the database credentials are passed through from the configuration of the pipeline.  It would be particularly good if we could retrieve the credentials via other mechanisms, such as retrieving them directly from a credentials repository such as Keycloak.
