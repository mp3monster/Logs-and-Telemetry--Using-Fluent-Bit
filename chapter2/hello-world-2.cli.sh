fluent-bit -q -i dummy -t dummy1 -p dummy="{\"hello\":\"my world\"}"  -i dummy -t dummy2 -p dummy="{\"more\":\"stuff\"}" -o stdout -m '*'
