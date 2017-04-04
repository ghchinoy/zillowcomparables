# Zillow Comparables

Given a zillow property ID, find comparable properties.

Requires a Zillow API Key (ZWSID) provided as an environment variable and an optional Google Maps API key 
(MAPSAPI) also as an envirionment variable, if distance from original location is desired.


example usage

```
ZWSID=X1-ZWzA28143veodv_a9fwu go run main.go 101473495
```

or

```
MAPSAPI=AIzaStfqUKBiJX12bbinv0iIAwaUzUQpuvrL--6U ZWSID=X1-ZWzA28143veodv_a9fwu go run main.go 101473495
```


## Zillow API

https://www.zillow.com/howto/api/GetDeepComps.htm

## Google Maps API

https://developers.google.com/maps/documentation/distance-matrix/intro


## To Do

* reconcile maps distance matrix with properties
