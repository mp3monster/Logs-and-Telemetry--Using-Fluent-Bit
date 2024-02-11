# Build and run the container and plugins

The plugin code is compiled and deployed to the Fluent Bit folder as part of the container build process.  We can run the build process using the Docker command in the folder */chapter10/go/src* 

`docker build . -t fluent-bit-gdb -f Dockerfile`

Note that the compilation errors will be displayed on the console. A compilation failure will not prevent the docker build process from progressing. We can then run the container with the command:

`docker run -it --rm fluent-bit-gdb`

It may be useful to combine the two statements into a single script, particularly during the development process.

## Use Cases



- slow queries and other DB metrics being managed as tables in the schema. Many databases will log details like slow queries within special tables in the database. Knowing this information can be very helpful to help tune the DB performance.

- Application audit - applications sometimes may record in the database auditable information such as user login and logout. When data is added or deleted. This information can be beneficial as part of an aggregated view of what is happening within the systems from a security perspective.

- Retrieve metadata that can be used to help enrich logs, metrics, and traces and provide additional meaningful context.



# Configuration Attributes

| Attribute Name   | Description                                                  | Input | Output | Example value                |
| ---------------- | ------------------------------------------------------------ | ----- | ------ | ---------------------------- |
| db_host          | Host address for the database server                         | Y     | Y      | 192.168.0.1                  |
| db_port          | The network port to communicate to the database with e.g. 5432 for Postgres or 3361 for MySQL | Y     | Y      | 5432                         |
| db_user          | The name of the user to authenticate as when communicating with the database | Y     | Y      | postgresUser                 |
| db_password      | The associated DB password for the named user. This needs to be in clear text | Y     | Y      | myPassword                   |
| ordering_col     | To retrieve the log records in the correct order we need to know which column to Order By in the constructed SQL. If not value is provided, then no order by clause is used and the records will be received based on the order the DB engine provides. We track the ordering_col so that each query cycle we don't reread any earlier records. | Y     | N      | mySeqId                      |
| table_name       | The name of the table from which we're going to retrieve records from or add records to. | Y     | Y      | myTable                      |
| db_name          | A DB Server may support multiple databases, therefore we need to identify which database by its name. | Y     | Y      | local                        |
| db_type          | To identify the database type (and therefore correct DB driver to use) the correct DB type is needed from a predefined list of values. Currently the only valid values are postgres and mysql | Y     | Y      | mysql                     |
| pk               | The primary key so, if we're asked to delete records once read, we can ensure that the correct records are deleted | Y     | N      | myId                         |
| delete           | A boolean flag to indicate whether the records read should be removed from the database once they're in the buffer. Deleting the records means we can't reconsume those records. | Y     | N      | true                         |
| where_expression | It maybe desirable to filter the records pulled from the source table. For example only retrieving records of a particular type or that have a specific attribute. e.g. a history of queries, and we only want those marked as slow, or where the execution time was greater than a predetermined threshold. If No value is provided then no where clause will be incorporated. This needs to be a correct SQL syntax | Y     | N      | execution_time > 500         |
| query_cols       | Identify the columns that need to be queried or have values inserted to. If no value is defined in the input, then the * wildcard is assumed and all columns will be retrieved. On the insert, if columns are named then only these columns will receive values. When provided the columns need to be expressed as a comma separated list | Y     | Y      | a_column, b_column, c_column |
| query_frequency  | The interval at which we will query the database to look for new records. This is an integer defining seconds | Y     | N      | 5                            |



