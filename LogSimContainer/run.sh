# we look to see if we've already got the log simulator code in the environment, if we haven't then retrieve it
# if we sdont have any configuration then pull the default and copy it into place
# if we have environment variable useExtn set to 1, then pull in the externsions
# output the configuration we've got to use
# run the log simulator with the identified run properties and test data

if test -f ./LogSimulator.groovy; then
  echo "got simulator"
else
  echo "getting log simulator"
  wget https://raw.githubusercontent.com/mp3monster/LogGenerator/FBBook/LogSimulator.groovy
fi

if test -f /vol/conf; then 
  echo "copying /vol/conf for use"
  cp /vol/conf/*.conf .
  # need to find the conf file and rename
  cp -r ./*.conf ./default.conf
else
  rm -f default.properties source.txt
  wget -q https://raw.githubusercontent.com/mp3monster/LogGenerator/FBBook/docker/default.properties
  wget -q https://raw.githubusercontent.com/mp3monster/LogGenerator/FBBook/docker/source.txt
fi

#get the extensions if the useExtn env var is set
if [${#useExtn} -eq 1]; then
  echo "Getting custom extensions"
  wget -q https://raw.githubusercontent.com/mp3monster/LogGenerator/FBBook/CustomConsoleOutputter.groovy
  wget -q https://raw.githubusercontent.com/mp3monster/LogGenerator/FBBook/CustomOCINotificationsOutputter.groovy
  wget -q https://raw.githubusercontent.com/mp3monster/LogGenerator/FBBook/CustomOCIOutputter.groovy
  wget -q https://raw.githubusercontent.com/mp3monster/LogGenerator/FBBook/CustomOCIQueueOutputter.groovy
fi

#if something isn't behaving properly - let's see what test data and configs we have
#ls /vol/conf
#echo
#ls /vol/test-data

groovy --version
 
 set data_set = ""
echo ${#data}
#if test -n $data; then 
if [ ${#data} -gt 1 ] ; then 
  data_set = "/vol/test-data/" + $data
fi

cat /vol/conf/$run_props

groovy ./LogSimulator.groovy /vol/conf/$run_props $data_set
