curl --location 'localhost:9881' \
--header 'Content-Type: application/json' \
--data '{
    "hello": "Fluent"
}' \
--output results.txt \
--include
