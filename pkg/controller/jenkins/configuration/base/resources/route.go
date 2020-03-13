package resources

import (
	routev1 "github.com/openshift/api/route/v1"
)

// PopulateRouteFromService takes the ServiceName and Creates the Route based on it
func PopulateRouteFromService(route *routev1.Route, name string) {
	route.Spec = routev1.RouteSpec{
		TLS: &routev1.TLSConfig{
			InsecureEdgeTerminationPolicy: routev1.InsecureEdgeTerminationPolicyRedirect,
			Termination:                   routev1.TLSTerminationEdge,
		},
		To: routev1.RouteTargetReference{
			Kind: "Service",
			Name: name,
		},
	}
}
