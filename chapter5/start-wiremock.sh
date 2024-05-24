docker run -it --rm -p 8080:8080 -p 8443:8443 --name wiremock --volume wiremock:/home/wiremock/mappings wiremock/wiremock:2.35.0  --verbose
# https://wiremock.org/docs/standalone/docker/
# http://localhost:8080/__admin/mappings
# http://localhost:8080/__admin/requests
