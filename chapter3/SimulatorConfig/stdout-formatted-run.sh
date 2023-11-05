# $flbBookRootDir - see Appendix A for configuration
docker run -v .:/vol/log -v $flbBookRootDir/chapter3/SimulatorConfig/:/vol/conf \
           -v $flbBookRootDir/TestData/:/vol/test-data \
           --env run_props=stdout-formatted.properties \
           --env data=medium-source.txt \
           logsimcontainer-logger 