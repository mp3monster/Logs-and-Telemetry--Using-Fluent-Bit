# Logs and Telemetry using Fluent Bit, Kubernetes, streaming and more
This repository contains the code, configurations, test data, and utilities for Manning's book [Logs and Telemetry using Fluent Bit]([https://www.manning.com/books/fluent-bit-with-kubernetes?utm_source=Phil&utm_medium=affiliate](https://www.manning.com/books/logs-and-telemetry?utm_source=Phil&utm_medium=affiliate)), as well as the book's examples and solutions. We have incorporated some additional examples within the chapter folders.

A quick summary of the book's chapters:

1. Introduction to Fluent Bit
2. From Zero to Hello World
3. Capturing Inputs
4. Getting inputs from Containers and Kubernetes
5. Outputting events
6. Parsing to extract more meaning
7. Filtering and transforming events
8. Stream Processors for time series calculations and filtering
9. Building processors and Fluent Bit extension options
10. Building Plugins
11. Putting Fluent Bit into action - an enterprise use case

Additional read-me documents are incorporated into the different folders to provide domain- or chapter-specific information.

![](https://blog.mp3monster.org/wp-content/uploads/2024/07/wilkins-hi.jpg?w=529)

### Demo Configurations

The chapters container demo configurations to illustrate different aspects of Fluent Bit.  This includes for some chapters scripts which wrap utilities / tools in a container. 

##### Cross-Platform Constraints

It is important to note that **a couple of demos CAN'T run on  (Mac and Windows)** as not all plugins support all platforms. We have focussed on all the demos working for Linux as this is the typical platform for containerized/Kubernetes solutions. The root of the issue is that there are a few plugins that have not been built for macOS (particularly for Apple silicon), and for Windows some of the OS level services do not have a direct equivalent.

We've noted the constraints in the text of the book, and the book also has a table showing which plugins are supported by which platforms.

##### Classic & YAML Configurations

Fluent Bit is slowly pushing towards YAML as the configuration standard, although with a couple of exceptions, the classic and YAML configuration formats work for all plugins. As a result, we have provided YAML versions of the configuration files (although during development, we've largely worked with the classic format first).

#### Extras

In addition to the content in the book, we have also incorporated some extra resources, such as additional configurations and details, such as Rail Road diagrams.  They are in a folder called [extras](https://github.com/mp3monster/Logs-and-Telemetry--Using-Fluent-Bit/tree/main/extras) or, where more appropriate, in the relevant folder.

##### Railroad diagrams

Details of the railroad diagrams can be found here, and an example diagram:

![](https://github.com/mp3monster/Logs-and-Telemetry--Using-Fluent-Bit/blob/main/extras/Syntax%20RailRoad%20Diagrams/Classic%20Configuration%20Format%20Railroad%20Diagram.png?raw=true)

## About the Book
This book covers Fluent Bit v2 onwards, focusing on its application and configuration. It's designed to cater to your specific interests in logging and monitoring, even if you haven't read Logging In Action.

It is being written as a partner title to [Logging In Action](https://www.manning.com/books/logging-in-action?a_aid=Phil) (although you don't need to have read [Logging in Action](https://www.manning.com/books/logging-in-action?a_aid=Phil)) to benefit from this book.



### Log Simulator

The book uses several third-party tools, all of which are referenced in Appendix A. The Log Simulator can be a very helpful tool for developing and testing Fluent Bit configurations (although it is currently limited to logs). All the resources can be found in the [Log Generator GitHub](https://github.com/mp3monster/LogGenerator) repository.

The docker scripts we've provided to run the Log Simulator (are in the folder `SimulatorConfig` in the relevant chapter folders). These scripts exist to just simplify the Docker command - primarily mapping the volumes and environment variable configuration settings for a specific scenario.

The command does include the parameters `-ti --init` .  This means that the Docker container will terminate with the use of `ctrl-c`.  These parameters do the following:

- `-ti` tells Docker to use a pseudo tty (terminal input) and to run interactively. As a result, you will see everything the Log Simulator sends to the console.
- `--init` command modifies the way the process is started within the container so that the 

The use of this `-ti --init` is nicely explained [here](https://www.baeldung.com/ops/docker-init-parameter) and you can read more about how `--init` works [here](https://github.com/krallin/tini).



## More information

For more information, you can either go to my [blog](https://blog.mp3monster.org/publication-contributions/fluent-bit-with-kubernetes/)
