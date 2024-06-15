REM Note - in the book we' pointed to the use of Docker compose. But if you want to short cvircuit that approach, then this should allow Prometheus to work. This command has more parameters than the start-prometheus standard script
docker run --rm -p 9091:9091 -p 2022:2022 --network=host --name=myPrometheus --mount type=bind,source=.\prometheus\prometheus.yml,destination=/etc/prometheus/prometheus.yml prom/prometheus --log.level=debug --config.file=/etc/prometheus/prometheus.yml --web.listen-address="0.0.0.0:9091"
