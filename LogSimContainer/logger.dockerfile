from groovy
LABEL maintainer "Phil Wilkins docker@mp3monster.org"
LABEL Description="LogSimulator docker image" Vendor="mp3monster" Version="1.1 Beta"
VOLUME ["/vol/conf", "/vol/log", "/vol/test-data"]
# use /vol/log - for targeting log outputs
COPY run.sh .
RUN  wget -nv https://raw.githubusercontent.com/mp3monster/LogGenerator/FBBook/docker/run.sh; 
CMD ./run.sh