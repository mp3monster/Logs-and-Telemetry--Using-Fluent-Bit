REM # docker run --rm --name my-prometheus --network-alias prometheus --publish 9090:9090 --volume prometheus.yml:/etc/prometheus/prometheus.yml --detach prom/prometheus
REM docker run --rm --name my-prometheus  --publish 9090:9090 --volume prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
REM docker run --rm -p 9090:9090 -p 2021:2021  -p 9885:9885   --name=myPrometheus --mount type=bind,source=.\prometheus\prometheus.yml,destination=/etc/prometheus/prometheus.yml prom/prometheus 
docker run --rm -p 9090:9090 -p 2021:2021    --name=myPrometheus --mount type=bind,source=.\prometheus\prometheus.yml,destination=/etc/prometheus/prometheus.yml prom/prometheus 

REM Feature Flags documentation https://prometheus.io/docs/prometheus/latest/feature_flags/
REM Remote write feature flag -   --enable-feature=remote-write-receiver