package yttdatavalues

import (
	"testing"
)

func TestGetFlattenedDataValues(t *testing.T) {
	antreaDataValue := `
#@data/values
#@overlay/match-child-defaults missing_ok=True
---
infraProvider: vsphere
antrea:
  config:
    serviceCIDR: 100.64.0.0/13
    trafficEncapMode: encap
    noSNAT: false
    tlsCipherSuites: TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384
    featureGates:
      AntreaProxy: true
      EndpointSlice: false
      AntreaPolicy: true
      NodePortLocal: false
      AntreaTraceflow: true
`
	expected := map[string]interface{}{
		"antrea.config.trafficEncapMode":             "encap",
		"antrea.config.noSNAT":                       false,
		"antrea.config.featureGates.NodePortLocal":   false,
		"antrea.config.serviceCIDR":                  "100.64.0.0/13",
		"infraProvider":                              "vsphere",
		"antrea.config.tlsCipherSuites":              "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384",
		"antrea.config.featureGates.AntreaProxy":     true,
		"antrea.config.featureGates.EndpointSlice":   false,
		"antrea.config.featureGates.AntreaPolicy":    true,
		"antrea.config.featureGates.AntreaTraceflow": true,
	}

	flattenedDataValues, err := GetFlattenedDataValues([]byte(antreaDataValue))
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range expected {
		if val, ok := flattenedDataValues[k]; ok {
			if val != v {
				t.Fatalf("Expected %v got %v", v, val)
			}
		} else {
			t.Fatalf("key %s not found", k)
		}
	}

}
