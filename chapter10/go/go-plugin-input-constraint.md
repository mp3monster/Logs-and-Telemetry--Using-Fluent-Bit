# Go Plugin Input Constraint 



### Challenge

To be able to have multiple plugins and track details for different instances of the same plugin, the FLBPluginInputCallback method needs to have a version that can be given the context.  We see this in the output plugin's methods *FLBPluginFlush* and *FLBPluginFlushCtx*.  

Looking within the code of Fluent Bit ([Go Proxy](https://github.com/fluent/fluent-bit/blob/689afa14ae8c06e8ae32c930eaeb7daa305705db/src/proxy/go/go.c#L241)), it appears to be possible to implement.

Without this, the input plugin can only generate hardwired values.

The issues have been reported in the issue [here](https://github.com/fluent/fluent-bit/issues/8464).  A second feature request has been identified where Go can be used for custom filters - found [here](https://github.com/fluent/fluent-bit/issues/8465).

## Workarounds

In the meantime, there are some possible workarounds.

#### Environment Vars

Push the configuration values into environment variables and then pull them back when needed by the FLBPluginInputCallback method. This method doesn't have access to any contextual information, so it can only be applied once.

A variation of this would be to store the configurations in a configuration file and load the file as part of the callback.

#### Hardwiring

As we are in complete control of the plugin, we can hardwire some or all of the attributes of the plugin into the plugin. By exploiting the approach shown with the Docker image - if the utilities and code are shipped, this becomes more like how Lua is deployed.

By each source having its own Plugin, the differing configuration needs are overcome. 

If you consider hardwiring some values, consider using the default tag mechanism described [here](https://pkg.go.dev/gopkg.in/mcuadros/go-defaults.v1#section-readme).

#### Compile Multiple instances

If your initial thought is to cut and paste the code multiple times and then cringe, that is entirely understandable. But we don't have to go as far as that.  The current implementation uses the environment variables mechanism; within the logic, we create the environment variable name using a constant declared that provides the plugin name.

The Go Lang build allows us to select which files to compile together (a great article in [Dave Cheney's blog](https://dave.cheney.net/2013/10/12/how-to-use-conditional-compilation-with-the-go-build-tool) on how this works.) By incorporating the constant into its file, we can create multiple versions of the file. Then, each build selects an instance of the file and a different target static object file.   The Fluent Bit plugins file needs to be updated to pick up the additional plugin objects.

This then just requires each use case to use a different instance of the plugin built. This obviously does mean that the runtime executable will get larger for each occurrence of a plugin's use.

#### Containerize a staging Layer.

While containers carry an overhead, separating the handling of the source to a containerized layer means you can inject different configurations into the container. Then, if you want to use the source plugin, each source has its own container instance and directs the captured events to a target using the Fluent Bit capabilities using HTTP, forward, etc.