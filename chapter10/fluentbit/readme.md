## Background

This folder contains test and example configurations for our Go based plugin.



## Plugin Configuration Attributes



The Plugin supports both input and output operations for databases. Initially, we've only incorporated drivers for two types of RelationalDB (PostgreSQL and MySQL), but the logic will easily support the addition of more drivers.



The input plugin name is **in_gdb**, and the corresponding output plugin is called **out_gdb**.



| Attribute        | Description                                                  | <u>*I*</u>n/<u>*O*</u>ut/<u>*B*</u>oth | Example Configuration   |
| ---------------- | ------------------------------------------------------------ | -------------------------------------- | ----------------------- |
| db_host          | The network address of the database server                   | B                                      | 192.168.1.1             |
| db_port          | The port exposed for connecting the database server.         | B                                      | 3306                    |
| db_user          | Username needed as part of the access credentials            | B                                      | joe                     |
| db_password      | Password for the username provided                           | B                                      | bloggs                  |
| ordering_col     | If we want the events to be retrieved in a specific order we need to provide the name of the column by which the records should be ordered. | I                                      | orderId                 |
| table_name       | The name of the table that is to be read or inserted into    | B                                      | myTable                 |
| db_name          | The database name contains the relevant table.               | B                                      | myDB                    |
| db_type          | Defines the database type to be used. This allows us to select the correct driver and make any appropriate adjustments to the SQL syntax necessary. Currently, only values of **mysql** and **psql** are supported. | B                                      | mysql                   |
| pk               | The primary key. We need to know this to target record deletion if we want the delete option to work. | I                                      | myPK                    |
| delete           | It takes a **true** or **false** value to determine whether it should delete the record once a value has been read from a database table. A **true** value will cause the record to be deleted once the record has been successfully consumed. | I                                      | true                    |
| query_cols       | Identifies the names of the columns that should be read from the query. This should be a comma-separated list | I                                      | colA, colB, another_col |
| where_expression | If we want to be selective about records retrieved, we need to supply a where statement. | I                                      | a=b                     |
| log_level        | This overrides the default log level settings for the plugin's use | B                                      | debug                   |

**Note**: IF Fluent Bit and the plugin are running within a container, then ensure that the db_host is visible inside the container. This can be solved several ways - such as explicitly defining the host address. Configuring the container orchestration so it uses the correct network so the address will resolve correctly.