set fn=log%1%.json
echo %fn
curl -X POST --location 127.0.0.1:9881 --header Content-Type:application/json --data @%fn%
