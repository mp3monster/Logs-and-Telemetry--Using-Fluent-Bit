from groovy
LABEL maintainer "Phil Wilkins docker@mp3monster.org"
LABEL Description="Dynamic execution of the LogSimulator docker image" Vendor="mp3monster" Version="1.0"
VOLUME ["/vol/conf", "/vol/log", "/vol/test-data"]
EXPOSE 80-9000
# use /vol/log - for targeting log outputs
COPY ./run.sh .
CMD ./run.sh