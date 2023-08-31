docker run --rm -p 9090:9090 -p 2021:2021    --name=myPrometheus --mount type=bind,source=.\prometheus\prometheus.yml,destination=/etc/prometheus/prometheus.yml prom/prometheus 

REM Feature Flags documentation https://prometheus.io/docs/prometheus/latest/feature_flags/
REM Remote write feature flag -   --enable-feature=remote-write-receiver