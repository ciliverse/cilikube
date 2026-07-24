package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func gwSection(name string) *gatewayv1.SectionName {
	s := gatewayv1.SectionName(name)
	return &s
}

// ShowcaseGatewaySeedObjects returns Gateway API sample objects for the fake clientset.
func ShowcaseGatewaySeedObjects() []runtime.Object {
	className := gatewayv1.ObjectName("showcase-nginx")
	ns := gatewayv1.Namespace("production")
	hostname := gatewayv1.Hostname("orders.showcase.local")
	path := gatewayv1.HTTPPathMatch{
		Type:  ptrPathMatchPathPrefix(),
		Value: strPtr("/"),
	}
	svcPort := gatewayv1.PortNumber(80)

	return []runtime.Object{
		&gatewayv1.GatewayClass{
			ObjectMeta: metav1.ObjectMeta{
				Name:   "showcase-nginx",
				Labels: map[string]string{"showcase": "true"},
			},
			Spec: gatewayv1.GatewayClassSpec{
				ControllerName: "showcase.cilikube.local/nginx-gateway",
			},
			Status: gatewayv1.GatewayClassStatus{
				Conditions: []metav1.Condition{{
					Type:   string(gatewayv1.GatewayClassConditionStatusAccepted),
					Status: metav1.ConditionTrue,
					Reason: "Accepted",
				}},
			},
		},
		&gatewayv1.Gateway{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "public-gw",
				Namespace: "production",
				Labels:    appLabels("orders-api"),
			},
			Spec: gatewayv1.GatewaySpec{
				GatewayClassName: className,
				Listeners: []gatewayv1.Listener{{
					Name:     "http",
					Protocol: gatewayv1.HTTPProtocolType,
					Port:     80,
					AllowedRoutes: &gatewayv1.AllowedRoutes{
						Namespaces: &gatewayv1.RouteNamespaces{
							From: ptrFromNamespaces(gatewayv1.NamespacesFromSame),
						},
					},
				}},
			},
			Status: gatewayv1.GatewayStatus{
				Addresses: []gatewayv1.GatewayStatusAddress{{
					Type:  ptrAddressType(gatewayv1.IPAddressType),
					Value: "203.0.113.10",
				}},
			},
		},
		&gatewayv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "orders-route",
				Namespace: "production",
				Labels:    appLabels("orders-api"),
			},
			Spec: gatewayv1.HTTPRouteSpec{
				CommonRouteSpec: gatewayv1.CommonRouteSpec{
					ParentRefs: []gatewayv1.ParentReference{{
						Name:        gatewayv1.ObjectName("public-gw"),
						Namespace:   &ns,
						SectionName: gwSection("http"),
					}},
				},
				Hostnames: []gatewayv1.Hostname{hostname},
				Rules: []gatewayv1.HTTPRouteRule{{
					Matches: []gatewayv1.HTTPRouteMatch{{
						Path: &path,
					}},
					BackendRefs: []gatewayv1.HTTPBackendRef{{
						BackendRef: gatewayv1.BackendRef{
							BackendObjectReference: gatewayv1.BackendObjectReference{
								Name: gatewayv1.ObjectName("orders-api"),
								Port: &svcPort,
							},
						},
					}},
				}},
			},
		},
	}
}

func ptrPathMatchPathPrefix() *gatewayv1.PathMatchType {
	t := gatewayv1.PathMatchPathPrefix
	return &t
}

func ptrFromNamespaces(v gatewayv1.FromNamespaces) *gatewayv1.FromNamespaces {
	return &v
}

func ptrAddressType(t gatewayv1.AddressType) *gatewayv1.AddressType {
	return &t
}
