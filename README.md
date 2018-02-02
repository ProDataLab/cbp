CodeDepot.tech ~ component package

The following are my notes while the initial code is translated from concept to implementation

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


