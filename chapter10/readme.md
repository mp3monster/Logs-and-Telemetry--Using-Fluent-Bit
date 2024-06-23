The repository for this chapter is divided into several parts. It covers the creation of a Fluent Bit plugin using Go.

- [code-baseline]() contains an empty skeleton for building Go plugins for Fluent Bit.
- [go](https://github.com/mp3monster/Fluentbit-with-Kubernetes/tree/main/chapter10/go) - this contains the chapter demo plugin source code, build scripts, and Docker containers.
- [fluentbit](https://github.com/mp3monster/Fluentbit-with-Kubernetes/tree/main/chapter10/fluentbit) - contains test configurations to exercise the Go plugin at an integrated level.  We also have documentation on all the parameters the plugin configuration needs.
- [sql](https://github.com/mp3monster/Fluentbit-with-Kubernetes/tree/main/chapter10/sql) - this contains simple SQL scripts to help configure [PostgreSQL](https://www.postgresql.org/) and [MySQL](https://www.mysql.com/) databases.