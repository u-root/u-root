# Example Service for Service of Services

the Service of Services (SoS) is a web-based UI for configuring system settings, on localhost:8000. There are a number of services in u-root by default, but you can easily build your own.

An SoS service consists of a number of files:
* `your_sos.go` contains main() and starts your service, in `cmds/your_sos`
* `service.go` contains all the functionality of your service, in `cmds/your_sos`
* `server.go` routes data between the webpage and the service, in `cmds/your_sos`
* `your_sos.html`, the HTML and Javascript frontend, placed in `pkg/sos/html`

The logic in service.go is heavily dependent on what you are trying to implement, but everything else follows this general pattern:
1. main constructs a new service, passes it to a new server, and starts the server.
2. the server gets a port from the SoS, builds its router, and adds itself to the SoS table.
3. the server's router maps URLs to user-defined Go functions.
4. the server's displayStateHandle function updates the service and renders the webpage with service data.
5. the webpage contains input fields, which are sent back as JSON through URLs defined in the router.
6. The server's router passes the JSON to a Go function, which it decodes and sends to a service method.
7. The service method operates on this input, updating its fields once it's done.
8. The webpage reloads on completion, rerendering with the updated service data.

the example Service in this dir is heavily commented with more details.

To start a service on startup, call it in your uinit.

To test services on your build machine:
1. navigate to `cmds/sos`, run `go build .`, then `sudo ./sos`
2. navigate to your service, `cmds/your_sos`, run `go build .`, then `./your_sos`

You can then open the SoS table in a browser at localhost:8000 and access your service.
This will use the builtin HTML string, not your HTML file, due to path issues. Until we
implement a flag, you can manually change `htmlRoot` in `pkg/sos/server.go` and rebuild sos.
