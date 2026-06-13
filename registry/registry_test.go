package registry

import (
	"testing"

	"github.com/stellhub/stellar/config"
)

func TestInstanceFromConfig(t *testing.T) {
	cfg := &config.RegistryConfig{
		Namespace:  "default",
		Service:    "order-service",
		InstanceID: "order-service-1",
		Zone:       "zone-a",
		TTL:        "30s",
		Labels: map[string]string{
			"version": "v1",
		},
		Metadata: map[string]string{
			"owner": "platform",
		},
		ServiceEndpoints: []config.RegistryServiceEndpointConfig{{
			Name:     "http",
			Protocol: "http",
			Host:     "127.0.0.1",
			Port:     8080,
		}},
	}

	instance, ok := InstanceFromConfig(cfg)
	if !ok {
		t.Fatalf("expected instance from config")
	}
	if instance.Namespace != "default" || instance.Service != "order-service" || instance.InstanceID != "order-service-1" {
		t.Fatalf("unexpected instance identity %#v", instance)
	}
	if instance.TTLSeconds != 30 {
		t.Fatalf("unexpected ttl seconds %d", instance.TTLSeconds)
	}
	if len(instance.Endpoints) != 1 || instance.Endpoints[0].Port != 8080 {
		t.Fatalf("unexpected endpoints %#v", instance.Endpoints)
	}
}

func TestFilterInstancesByZoneAndLabels(t *testing.T) {
	values := []Instance{
		{
			Namespace:  "default",
			Service:    "order-service",
			InstanceID: "order-service-1",
			Zone:       "zone-a",
			Labels:     map[string]string{"version": "v1"},
		},
		{
			Namespace:  "default",
			Service:    "order-service",
			InstanceID: "order-service-2",
			Zone:       "zone-b",
			Labels:     map[string]string{"version": "v2"},
		},
	}

	filtered := filterInstances(values, Query{
		Namespace: "default",
		Service:   "order-service",
		Zone:      "zone-a",
		Labels:    []string{"version=v1"},
	})
	if len(filtered) != 1 || filtered[0].InstanceID != "order-service-1" {
		t.Fatalf("unexpected filtered instances %#v", filtered)
	}
}
