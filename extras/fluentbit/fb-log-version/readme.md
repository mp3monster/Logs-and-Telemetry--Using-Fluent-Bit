# Logging Fluent Bit deployment information

It can be useful to log the application and related information when a service is started (or restarted / hot deployed). Logging that information can show when and where a (re)start occurred. Fluent Bit does write to the console on startup its version information. But capturing this and excluding all other information is a little messy until now.

As of release 2.2, a new filter has been introduced, [sysinfo](https://docs.fluentbit.io/manual/pipeline/filters/sysinfo). The filter can add to a log event with the details of the Fluent Bit version and host environment.  As it is a filter, we do need to trigger the filter from an input.  We don't want to tag every log event with this information.

The solution to only running the Filter to append the necessary information is to have a dummy input that only executes once (controlled by the samples) attribute. The Dummy input creates a simple JSON payload, which can be extended by setting an environment variable called __ENV_INFO__, at which point that value will get incorporated into the dummy output.

To use this configuration, the parent configuration file simply needs to include *(@include) report-fb-version.conf*, and then have a match declaration that will trap *FLBVersion* 
