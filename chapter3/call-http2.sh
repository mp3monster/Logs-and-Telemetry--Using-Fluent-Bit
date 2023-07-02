curl --location 'localhost:9891' \
--header 'Content-Type: application/json' \
--data '{
    "hello": "to another Fluent endpoint"
}' \
--include \
--output response.txt
