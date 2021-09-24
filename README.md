# affinity-test-apps

A client and a server app that can be used to test connection affinity when 
Kubernetes service endpoints change. The server responds client's requests
with its unique id (e.g., hostname). The client then validates that the received
returned id in the server's reply continues to be the same for all the
subsequent requests.

A container image containing both can be fetched from
[here](https://hub.docker.com/r/cilium/affinity-test-apps).

