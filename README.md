CodeDepot.tech ~ component package

**for any particular language**
1. Create a "component" library|module|package
  * is essentially a nanomsg wrapper
    * set socket type
    * set transport type
    * set url
    * has a run function that 
      * feeds the underlying function with input
      * sends the underlying function's output to out-socket
2. Create a generator app
  * supply a filepath (or url) and a function name.
  * function is extracted and made into its on library|module|package.
    * the "component" helper library is imported
  * both the input param types and output types are then extracted.
  * messagepack is used to make a message type for both input and output
  * a wrapper function is generated that 
    * imports component package
    * 
    * unmarshals the input off of the in-socket
    * calls the underlying function
    * marshals the output to a msgpack msg 
    * sends it to the next component via out-socket

# Types

  - Id Type Contains:
    - name    string
    - number  int64

  - cbp package:

    - Component:
      - id                Id
      - inSockets         []*Socket
      - outSockets        []*Socket
      - configSockets     []*Socket
      - reportingSockets  []*Socket

    - Connection:
      - id  Id
      - upstreamComponent   *Component
      - downstreamComponent *Component

    - Composite:
      - id          Id
      - components  []*Component
      - connections []*Connection
      - head        *Component
      - tail        *Component

# Component Repository

  - it is important that the author information of the orignal code be maintained
  - the repository item will contain:

  - The following tools are required
    - github client
    - language parser
    - database of (url, funcname)
      - url to the file that contains the function definition
      - the exact name of the function
    - component executable generator
    - repo/artifact for code/exe storage


# User Space

## Function Wrapper

  - creates a main
  - creates a function wrapper for their business logic
    - gets inChannel outChannel from Component
    - calls supplied function with args from msg of inChannel
    - sends output to outChannel

## API
  - http
  - web  
  - cli 