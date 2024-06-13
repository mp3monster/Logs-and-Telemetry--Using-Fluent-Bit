# Build and run the container and plugins

The plugin code is compiled and deployed to the Fluent Bit folder as part of the container build process.  We can run the build process using the Docker command in the folder */chapter10/go/src* 

`docker build . -t fluent-bit-gdb -f Dockerfile`

Note that the compilation errors will be displayed on the console. A compilation failure will not prevent the docker build process from progressing. We can then run the container with the command:

`docker run -it --rm fluent-bit-gdb`

It may be useful to combine the two statements into a single script, particularly during the development process.

We have provided scripts:
- build.[bat|sh] 
- clean.[bat|sh] 
These simplify the process with build incorporating the docker commands. Clean will remove the created container and release all the resources used. This is handy - if you have tinkered with the current image and want to be sure of a clean build or release resources from failed builds this clean script will tidy up.

### Use Cases

- slow queries and other DB metrics being managed as tables in the schema. Many databases will log details like slow queries within special tables in the database. Knowing this information can be very helpful in tuning the DB performance.

- Application audit - applications sometimes may record in the database auditable information such as user login and logout. When data is added or deleted. This information can be beneficial as part of an aggregated view of what is happening within the systems from a security perspective.

- Retrieve metadata that can be used to help enrich logs, metrics, and traces and provide additional meaningful context.



## Configuration Attributes

| Attribute Name   | Description                                                  | Input | Output | Example value                |
| ---------------- | ------------------------------------------------------------ | ----- | ------ | ---------------------------- |
| plugin_instance_id | Optional to give the configuration - se we can see in the logs which plugin instance is generating log events. Only allowed a-zA-Z0-9 | Y | Y | plugin1 |
| db_host          | Host address for the database server                         | Y     | Y      | 192.168.0.1                  |
| db_port          | The network port to communicate to the database with e.g. 5432 for Postgres or 3361 for MySQL | Y     | Y      | 5432                         |
| db_type          | To identify the database type (and therefore correct DB driver to use) the correct DB type is needed from a predefined list of values. Currently, the only valid values are Postgres and MySQL | Y     | Y      | mysql                     |
| db_user          | The name of the user to authenticate as when communicating with the database | Y     | Y      | postgresUser                 |
| db_password      | The associated DB password for the named user. This needs to be in clear text | Y     | Y      | myPassword                   |
| db_name          | A DB Server may support multiple databases, therefore we need to identify which database by its name. | Y     | Y      | local                        |
| table_name       | The name of the table from which we're going to retrieve records from or add records to. | Y     | Y      | myTable                      |
| query_cols       | Identify the columns that need to be queried or have values inserted. If no value is defined in the input, then the * wildcard is assumed and all columns will be retrieved. On the insert, if columns are named then only these columns will receive values. When provided the columns need to be expressed as a comma-separated list | Y     | Y      | a_column, b_column, c_column |
| ordering_col     | To retrieve the log records in the correct order we need to know which column to Order By in the constructed SQL. If not value is provided, then no order by clause is used and the records will be received based on the order the DB engine provides. We track the ordering_col so that each query cycle we don't reread any earlier records. | Y     | N      | mySeqId                      |
| pk               | The primary key so, if we're asked to delete records once read, we can ensure that the correct records are deleted | Y     | N      | myId                         |
| delete           | A boolean flag to indicate whether the records read should be removed from the database once they're in the buffer. Deleting the records means we can't re-consume those records. | Y     | N      | true                         |
| where_expression | It may be desirable to filter the records pulled from the source table. For example only retrieving records of a particular type or that have a specific attribute. e.g. a history of queries, and we only want those marked as slow, or where the execution time was greater than a predetermined threshold. If No value is provided then no where clause will be incorporated. This needs to be a correct SQL syntax | Y     | N      | execution_time > 500         |
| query_frequency  | The interval at which we will query the database to look for new records. This is an integer defining seconds | Y     | N      | 5                            |


## Notes About the Build dependencies and the Dockerfile implications

### Makefiles
In each of the folders - (in_gdb) and (out_gdb) is a Makefile - so it is possible to run the build locally.  We originally used the Makefiles in the Docvkerfile, but have since simplified the Dockerfile so it directly calls to go build process.

### Docker / Local Build dependency
As the plugin development is dependent upon cgo we have an implicit dependency for GLIBC. This means without a version of GLIBC in the OS image. In addition to this, we also depend upon the Go networking to support the remote calling to the databases.

This means we need to be careful of our environment to ensure we have the correct version of GLIBC (2.32 or better). This can be a challenge for some environments.  In the Docker container, we have addressed this by selecting a base image with the correct GLIBC necessary (the Go Lang bullseye image).  We have then provided the build process with additional configuration controls such as setting the CGO feature to on in the environment variables and various   (Thanks to Patrick Stephens - for helping us get these right https://github.com/patrick-stephens).

If the correct GLIBC isn't available then the execution of Fluent Bit will fail with an error message reporting the expected version.

The following illustrates what a successful run of the clean and build script can look like:

`~>clean`

`~>docker rmi fluent-bit-gdb`
`Untagged: fluent-bit-gdb:latest`
`Deleted: sha256:7dbd2f830faf84a2911cb11dfbb7e38d87b3f1215dc7c148c6d5a3236d86e88c`

`~>docker builder prune -f`
`ID                                              RECLAIMABLE     SIZE            LAST ACCESSED`
`hop8oaygy4cqp9jj3t9qcfqiz*                      true            7.053MB         About a minute ago`
`yb23xa0pu2cdjlqhk5h3w8otb*                      true    46.94kB         About a minute ago`
`hm9ikrpxjrw8nrudg2t7qb644                       true    76B             About a minute ago`
`mqxnvbac1twmcnenhg7t6mzuc*                      true    0B              About a minute ago`
`nvahtqo6vozcjy6tpagiu6f0d*                      true    1.176kB         About a minute ago`
`k2py3vj05apuiqtz26d7gn615                       true    78.81MB         About a minute ago`
`r6aq486nrnwxcwkztffcjdxxy                       true    822B            About a minute ago`
`tegk3xbz1agxuct5tyomgp711                       true    6.204MB         About a minute ago`
`x588x1ho980k0oq7xcqvy49xv                       true    4.403MB         About a minute ago`
`f4exn9o4eicpmi4ugl76c69zh                       true    357B            About a minute ago`
`791qmjdk1dk35e5m26o4ygt8h                       true    6.356MB         About a minute ago`
`t7vpqahyd051xv3vf5n43de94                       true    6.799kB         About a minute ago`
`n32067gxz9eoqigky3jk7zmb7                       true    7.007kB         About a minute ago`
`lvvfouar9dbn65ha2rqxmrnqh                       true    25.17kB         About a minute ago`
`d8iedn04ovby49mcnklrhwk5e                       true    25.17kB         About a minute ago`
`rd24tp7na463erpuv3r0gpk7k                       true    46.94kB         About a minute ago`
`e56tio083podr3qbtxmtn35dg                       true    0B              About a minute ago`
`w24wbva5m8yobwmaipr54s9vy                       true    0B              About a minute ago`
`iph4aj4c4coa4noyqtwl9v0l7                       true    0B              About a minute ago`
`5gs9lzfcpmdwzp5yf8nad7wme                       true    0B              About a minute ago`
`2bcvbsqhtmozc1nbkqm4flhlq                       true    0B              About a minute ago`
`xym0peipw5rfhb1di9s40e7r5                       true    0B              About a minute ago`
`kst3djfj2kl1e87yyk37ql7cj                       true    0B              About a minute ago`
`xzbydluidcl33tn7olldmxh9s                       true    0B              About a minute ago`
`Total:  103MB`

`~>build`

`~>docker build . -t fluent-bit-gdb -f Dockerfile`
`[+] Building 30.7s (24/24) FINISHED                                                                                                                                                          docker:default`
 `=> [internal] load build definition from Dockerfile                                                                                                                                                   0.1s`
 `=> => transferring dockerfile: 1.18kB                                                                                                                                                                 0.0s`
 `=> [internal] load metadata for docker.io/fluent/fluent-bit:3.0.6-debug                                                                                                                               0.8s`
 `=> [internal] load metadata for docker.io/library/golang:1.21-bullseye                                                                                                                                1.3s`
 `=> [auth] library/golang:pull token for registry-1.docker.io                                                                                                                                          0.0s`
 `=> [auth] fluent/fluent-bit:pull token for registry-1.docker.io                                                                                                                                       0.0s`
 `=> [internal] load .dockerignore                                                                                                                                                                      0.0s`
 `=> => transferring context: 2B                                                                                                                                                                        0.0s`
 `=> [gobuilder  1/11] FROM docker.io/library/golang:1.21-bullseye@sha256:9fb61b54c8daa098ed318cb712032c884c85db3993c92d501921caaf309686c6                                                             14.3s`
 `=> => resolve docker.io/library/golang:1.21-bullseye@sha256:9fb61b54c8daa098ed318cb712032c884c85db3993c92d501921caaf309686c6                                                                          0.0s`
 `=> => sha256:9fb61b54c8daa098ed318cb712032c884c85db3993c92d501921caaf309686c6 9.10kB / 9.10kB                                                                                                         0.0s`
 `=> => sha256:218e4348e1199729915ff12a225132b02939db0a72311cca53df57f04ae8610e 2.87kB / 2.87kB                                                                                                         0.0s`
 `=> => sha256:3d53ef4019fc129ba03f90790f8f7f28fd279b9357cf3a71423665323b8807d3 55.10MB / 55.10MB                                                                                                       3.2s`
 `=> => sha256:2f6f16bc4179a0d96fdddf91e9cff7084916b616e521aa4d6a87300aa8e7950b 2.32kB / 2.32kB                                                                                                         0.0s`
 `=> => sha256:08f0bf643eb6745d5c7e9bada33de1786ab2350240206a1956fa506a1b47b129 15.76MB / 15.76MB                                                                                                       0.8s`
 `=> => sha256:6b037c2b46ab4e54a261a0ca65b12b93e00ca052e72765c9cc4caf1262a2b86c 54.59MB / 54.59MB                                                                                                       3.3s`
 `=> => sha256:fe8da5e369b749c1fab2438dff3b2ee387abc3d72a0d4e47ac257c021af0c71b 85.93MB / 85.93MB                                                                                                       4.7s`
 `=> => extracting sha256:3d53ef4019fc129ba03f90790f8f7f28fd279b9357cf3a71423665323b8807d3                                                                                                              2.6s`
 `=> => sha256:360cc13c9ded6b9e4517c34383b52c4abc14ad706ab7516ff8fc4f9c0d67338c 67.01MB / 67.01MB                                                                                                       5.3s`
 `=> => sha256:d60761787513d26b5b3cff7017b509165ed08ffe1a7bf1b7be952adcd04e7ef4 125B / 125B                                                                                                             3.5s`
 `=> => sha256:4f4fb700ef54461cfa02571ae0db9a0dc1e0cdb5577484a6d75e68dc38e8acc1 32B / 32B                                                                                                               3.6s`
 `=> => extracting sha256:08f0bf643eb6745d5c7e9bada33de1786ab2350240206a1956fa506a1b47b129                                                                                                              0.3s`
 `=> => extracting sha256:6b037c2b46ab4e54a261a0ca65b12b93e00ca052e72765c9cc4caf1262a2b86c                                                                                                              1.8s`
 `=> => extracting sha256:fe8da5e369b749c1fab2438dff3b2ee387abc3d72a0d4e47ac257c021af0c71b                                                                                                              1.7s`
 `=> => extracting sha256:360cc13c9ded6b9e4517c34383b52c4abc14ad706ab7516ff8fc4f9c0d67338c                                                                                                              3.7s`
 `=> => extracting sha256:d60761787513d26b5b3cff7017b509165ed08ffe1a7bf1b7be952adcd04e7ef4                                                                                                              0.0s`
 `=> => extracting sha256:4f4fb700ef54461cfa02571ae0db9a0dc1e0cdb5577484a6d75e68dc38e8acc1                                                                                                              0.0s`
 `=> CACHED [stage-1 1/5] FROM docker.io/fluent/fluent-bit:3.0.6-debug@sha256:d4f8efddc55e45a4c45a1eeb58207a848b41f4201bc529b22d4b3f4b796b8182                                                          0.0s`
 `=> [internal] load build context                                                                                                                                                                      0.0s`
 `=> => transferring context: 47.64kB                                                                                                                                                                   0.0s`
 `=> [gobuilder  2/11] WORKDIR /root                                                                                                                                                                    0.5s`
 `=> [gobuilder  3/11] COPY / /root/                                                                                                                                                                    0.1s`
 `=> [gobuilder  4/11] COPY /common/* /root/out/                                                                                                                                                        0.1s`
 `=> [gobuilder  5/11] COPY /common/* /root/in/                                                                                                                                                         0.1s`
 `=> [gobuilder  6/11] COPY /out_gdb/* /root/out                                                                                                                                                        0.1s`
 `=> [gobuilder  7/11] COPY /in_gdb/* /root/in                                                                                                                                                          0.1s`
 `=> [gobuilder  8/11] RUN go mod edit -replace github.com/fluent/fluent-bit-go=github.com/fluent/fluent-bit-go@master                                                                                  0.3s`
 `=> [gobuilder  9/11] RUN go mod tidy                                                                                                                                                                  0.8s`
 `=> [gobuilder 10/11] RUN go build -C out -buildmode=c-shared  -a -gcflags=all="-C -l -B" -ldflags="-w -s" -trimpath  -tags netgo,osusergo -o /root/out/out_gdb.so                                     6.2s`
 `=> [gobuilder 11/11] RUN go build -C in -buildmode=c-shared  -a -gcflags=all="-C -l -B" -ldflags="-w -s" -trimpath  -tags netgo,osusergo -o /root/in/in_gdb.so                                        6.1s`
 `=> [stage-1 2/5] COPY --from=gobuilder /root/out/out_gdb.so /fluent-bit/bin/                                                                                                                          0.1s`
 `=> [stage-1 3/5] COPY --from=gobuilder /root/in/in_gdb.so /fluent-bit/bin/                                                                                                                            0.1s`
 `=> [stage-1 4/5] COPY --from=gobuilder /root/fluent-bit.conf /fluent-bit/etc/                                                                                                                         0.1s`
 `=> [stage-1 5/5] COPY --from=gobuilder /root/plugins.conf /fluent-bit/etc/                                                                                                                            0.1s`
 `=> exporting to image                                                                                                                                                                                 0.2s`
 `=> => exporting layers                                                                                                                                                                                0.2s`
 `=> => writing image sha256:5958bf3a3268f04b58aa78a908131c1a1f54c287815a8704895a703676615edd                                                                                                           0.0s`
 `=> => naming to docker.io/library/fluent-bit-gdb                                                                                                                                                      0.0s`

`View build details: docker-desktop://dashboard/build/default/default/w3zeywsdq5lwfbn68q8ex8dql`

`What's Next?`
  `View a summary of image vulnerabilities and recommendations ? docker scout quickview`

`~>docker run -it --rm fluent-bit-gdb`
`Fluent Bit v3.0.6`
* `Copyright (C) 2015-2024 The Fluent Bit Authors`
* `Fluent Bit is a CNCF sub-project under the umbrella of Fluentd`
* `https://fluentbit.io`

`___________.__                        __    __________.__  __          ________
\_   _____/|  |  __ __   ____   _____/  |_  \______   \__|/  |_  ___  _\_____  \`
 `|    __)  |  | |  |  \_/ __ \ /    \   __\  |    |  _/  \   __\ \  \/ / _(__  <
 |     \   |  |_|  |  /\  ___/|   |  \  |    |    |   \  ||  |    \   / /       \
 \___  /   |____/____/  \___  >___|  /__|    |______  /__||__|     \_/ /______  /`
     `\/                     \/     \/               \/                        \/`

`2024/06/05 10:09:48 [out_gdb] Register called`
`2024/06/05 10:09:48 [out_gdb] Registration result =false`
`2024/06/05 10:09:48 [in_gdb] Register called`
`[2024/06/05 10:09:48] [ info] [fluent bit] version=3.0.6, commit=9af65e2c36, pid=1`
`[2024/06/05 10:09:48] [ info] [storage] ver=1.5.2, type=memory, sync=normal, checksum=off, max_chunks_up=128`
`[2024/06/05 10:09:48] [ info] [cmetrics] version=0.9.0`
`[2024/06/05 10:09:48] [ info] [ctraces ] version=0.5.1`
`[2024/06/05 10:09:48] [ info] [input:in_gdb:in_gdb.0] initializing`
`[2024/06/05 10:09:48] [ info] [input:in_gdb:in_gdb.0] storage_strategy='memory' (memory only)`
`[2024/06/05 10:09:48] [ info] [input:in_gdb:in_gdb.0] thread instance initialized`
`2024/06/05 10:09:48 [out_gdb]Defaulting query columns to *`
`2024/06/05 10:09:48 [out_gdb] instance_1 Init connection test successful true`
`2024/06/05 10:09:48 Adding to context params==>{"pgnname":"out_gdb","instNme":"instance_1","host":"192.168.1.135","port":"3306","usr":"demo","pw":"demo","dbnme":"demo","cols":"*","seqr":"a_key","tbl":"plugindest","pk":"a_key","dbtype":"mysql","freq":1}
[2024/06/05 10:09:48] [ info] [sp] stream processor started
[2024/06/05 10:09:48] [ info] [output:stdout:stdout.1] worker #0 started
2024/06/05 10:09:49 [] Query constructed:SELECT COUNT(*) FROM pluginsrc`
`2024/06/05 10:09:49 [] Query constructed:SELECT a_key, a_string FROM pluginsrc ORDER BY a_key LIMIT 1`
`2024/06/05 10:09:49 execQuery row being sent = map[a_key:10001 a_string:record one]`
`2024/06/05 10:09:49 KeyList=[10001] Last Sequence Id=10001,  delete is false`
`2024/06/05 10:09:49 [in_gdb] InputCallback - retrieved data [2024-06-05 10:09:49.105166176 +0000 UTC m=+0.429097637 map[a_key:10001 a_string:record one]]`
`2024/06/05 10:09:50 [] Query constructed:SELECT COUNT(*) FROM pluginsrc WHERE a_key > 10001
2024/06/05 10:09:50 [] Query constructed:SELECT a_key, a_string FROM pluginsrc WHERE a_key > 10001 ORDER BY a_key LIMIT 1
2024/06/05 10:09:50 execQuery row being sent = map[a_key:10002 a_string:record one]
2024/06/05 10:09:50 KeyList=[10002] Last Sequence Id=10002,  delete is false
2024/06/05 10:09:50 [in_gdb] InputCallback - retrieved data [2024-06-05 10:09:50.105176353 +0000 UTC m=+1.429107804 map[a_key:10002 a_string:record one]]
2024/06/05 10:09:51 [] Query constructed:SELECT COUNT(*) FROM pluginsrc WHERE a_key > 10002`
`2024/06/05 10:09:51 [] Query constructed:SELECT a_key, a_string FROM pluginsrc WHERE a_key > 10002 ORDER BY a_key LIMIT 1`
`2024/06/05 10:09:51 execQuery row being sent = map[a_key:10003 a_string:record 3]`
`2024/06/05 10:09:51 KeyList=[10003] Last Sequence Id=10003,  delete is false`
`2024/06/05 10:09:51 [in_gdb] InputCallback - retrieved data [2024-06-05 10:09:51.105131549 +0000 UTC m=+2.429063000 map[a_key:10003 a_string:record 3]]`
`2024/06/05 10:09:52 [] Query constructed:SELECT COUNT(*) FROM pluginsrc WHERE a_key > 10003`
`2024/06/05 10:09:52 [] Query constructed:SELECT a_key, a_string FROM pluginsrc WHERE a_key > 10003 ORDER BY a_key LIMIT 1`
`2024/06/05 10:09:52 execQuery row being sent = map[a_key:10004 a_string:record 3]`
`2024/06/05 10:09:52 KeyList=[10004] Last Sequence Id=10004,  delete is false`
`2024/06/05 10:09:52 [in_gdb] InputCallback - retrieved data [2024-06-05 10:09:52.105193364 +0000 UTC m=+3.429124825 map[a_key:10004 a_string:record 3]]`
`^C[2024/06/05 10:09:53] [engine] caught signal (SIGINT)`
`[2024/06/05 10:09:53] [ warn] [engine] service will shutdown in max 5 seconds`
`[0] db1: [[1717582189.105166176, {}], {"a_key"=>"10001", "a_string"=>"record one"}]`
`[1] db1: [[1717582190.105176353, {}], {"a_string"=>"record one", "a_key"=>"10002"}]`
`[2] db1: [[1717582191.105131549, {}], {"a_key"=>"10003", "a_string"=>"record 3"}]`
`[3] db1: [[1717582192.105193364, {}], {"a_key"=>"10004", "a_string"=>"record 3"}]`
`2024/06/05 10:09:53 [out_gdb]instance_1 Flush called with context`
`2024/06/05 10:09:53 [out_gdb]instance_1 FLBPluginFlushCtx about to process:0 with timestamp 2024-06-05 10:09:49.105166176 +0000 UTC`
`[out_gdb]instance_1 insert expression: INSERT INTO plugindest (a_string,a_key) VALUES ('record one','10001')2024/06/05 10:09:53 [out_gdb]instance_1 FLBPluginFlushCtx about to process:0 with timestamp 2024-06-05 10:09:50.105176353 +0000 UTC`
`[out_gdb]instance_1 insert expression: INSERT INTO plugindest (a_string,a_key) VALUES ('record one','10002')2024/06/05 10:09:53 [out_gdb]instance_1 FLBPluginFlushCtx about to process:0 with timestamp 2024-06-05 10:09:51.105131549 +0000 UTC`
`[out_gdb]instance_1 insert expression: INSERT INTO plugindest (a_string,a_key) VALUES ('record 3','10003')2024/06/05 10:09:53 [out_gdb]instance_1 FLBPluginFlushCtx about to process:0 with timestamp 2024-06-05 10:09:52.105193364 +0000 UTC`
`[out_gdb]instance_1 insert expression: INSERT INTO plugindest (a_key,a_string) VALUES ('10004','record 3')[2024/06/05 10:09:53] [ info] [engine] service has stopped (0 pending tasks)`
`2024/06/05 10:09:53 [out_gdb]instance_1 Flush called with context`
`2024/06/05 10:09:53 [out_gdb] Flushing environment params`
`[2024/06/05 10:09:53] [ info] [output:stdout:stdout.1] thread worker #0 stopping...`
`[2024/06/05 10:09:53] [ info] [output:stdout:stdout.1] thread worker #0 stopped`
`2024/06/05 10:09:53 [in_gdb] Flushing environment params`
`2024/06/05 10:09:53 [out_gdb] Unregister called`
`Terminate batch job (Y/N)?`
