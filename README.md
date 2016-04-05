# Http Async Router

## Purpose

This tiny server is written in technical purpose to route requests between single AVK Billing instance and multiple
ESB instances, thus allowing to run multiple test environments using only one instance of Billing

## Known limitations

Requests initialized from Billing can only be directed to default ESB instance.
Therefore across any amount of instances there must be a default one.

## Concepts

### Requests

Request are identified using `COREQUEST` id provided by Billing

Router is attaching `COREQUEST` value of each request it manages to extract it for
in `X-Request-Id` header of response.

if `X-Request-Id` value is `0` it means what this request is not manageable (not async)
and is just proxied to target without caching of routing information

### Routing

Routing information is based on routes, not IPs, thus allowing to run many instances on single IP.

Therefore EACH instance connected to the router SHOULD use its own configured route.

Routing of a request is accomplished using combination of `COREQUEST` and route it originally came from.
i.e. then `COREQUEST` id is become known by the router, its ID and Source route is remembered in cache.

`COREQUEST` and routing information is cached and stored in-memory for time defined using `bucket-ttl`
parameter of config file.

### Routes (Endpoints)

Route/Endpoint - is a unique http path with configured settings

#### type
each route has a `type` - AVK or ESB accordingly to its client 

#### target-endpoint
Is an URL to which requests passed to this route is being proxied.
response from this URL is returned as router response

#### default-callback-path
Is an URL to which by default all async responses will be proxied.
It SHOULD BE explicitly defined for EACH route.

for ESB type of routes it can be dynamically overridden in runtime using `X-Callback-Path` header in request

#### content-type
Is an expected type of payload passed to request target
for now - it should always be `application/json` for ESB
and `application/xml` (important - NOT `text/xml`) for AVK

### Routing Cluster
Routing cluster is a valid 1-N combination of managed routes.
It SHOULD always consist of SINGLE instance of AVK Billing and 1-N instances of ESB

Clusters are formed implicitly - using combination of routes

Technically - the router allows to manage requests between any amount of instances of ESB and Billing
The ONLY mandatory condition that should be maintained is: 1 Billing vs 1-N Esb instances

Should you decide to configure more than one routing cluster - make sure they are not intersected with each other.

## Configuration

Configuration lookup is powered using [GoLang Viper](https://github.com/spf13/viper) library
and is configured to watch for config variables in:
- `/etc/har/config.yml` (takes precedence if exists)
- `./config.yml` (current workdir)

Following is a sample config file contents, with commentary

```yaml
bind-to: "192.168.253.60:80" # interface to bind in standard iface notation - [ip]:port
bucket-ttl: 6h # cache time - request route is remembered for each request for this time. format is \d+(h|s)
routes:  # route configuration - each instance connected should use separate route
  /request/dev/esb: # route path for ESB node 1
    type: ESB # route type = (ESB|AVK)
    target-endpoint: "http://192.168.253.41/GATE_JSON/index.php" # requests are proxied here (AVK node 1)
    default-callback-path: "http://192.168.253.53/esb/api/v10/esb_gateway" # async requests from AVK will be directed here (ESB node 1 real endpoint)
    content-type: "application/json"  # translated content-type
  /request/stable/esb: # route path for ESB node 2
    type: ESB
    target-endpoint: "http://192.168.253.41/GATE_JSON/index.php" # requests are proxied here (AVK node 1)
    default-callback-path: "http://192.168.253.54/esb/api/v10/esb_gateway" # async requests from AVK will be directed here (ESB node 2 real endpoint)
    content-type: "application/json"
  /esb/api/v10/esb_gateway: # route path for AVK node 1
    type: AVK
    target-endpoint: "http://192.168.253.53/esb/api/v10/esb_gateway" # so called "default" ESB node address
    default-callback-path: "http://192.168.253.41/GATE_JSON/index.php" # async requests from ESB will be directed here (AVK node 1 real endpoint)
    content-type: "application/xml"
```

## Logging

All logs are going to STDOUT

## Launch

Currently no so called "Classic" UNIX daemon features is implemented.

Launch, terminal unbinding and log targeting is done using following command:

	./HAR-v0.3.0 >> har.log 2>&1 &
	
No log rotation is implemented either.