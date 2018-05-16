package services

import (
	"testing"
)

var validMapperMsisdns = []struct {
	msisdn  string
	carrier string
}{
	{"38640123456", "A1"},
	{"38641123456", "Telekom Slovenije"},
	{"386651123456", "SŽ - Infrastruktura"},
}

var invalidMapperMsisdns = []string{
	"0038640123456", // zero prefixed
	"+38640123456",  // plus prefixed
	"38611123456",   // inexistent carrier
	"38540123456",   // different country
}

// We need this because working directory when running tests is
// msisdninfo/services, while data dir is msisdninfo/carriers
// Default create uses data dir "./carriers"
func createTestMapper() *PhonenumberToCarrierMapper {
	return &PhonenumberToCarrierMapper{
		carrierDataDir: "../carriers/",
	}
}

// see comment at createTestMapper()
// also cover default service factory function
func TestGetCarrierWithInvalidDataDir(t *testing.T) {
	mapper := CreatePhonenumberToCarrierMapper()
	hasInfo, info := mapper.GetCarrier(386, "38640123456")

	if hasInfo || info.Name != "" {
		t.Error("CreatePhonenumberToCarrierMapper should have wrong data dir and return empty result")
	}
}

func TestGetCarrierWithInvalidMsisdns(t *testing.T) {
	mapper := createTestMapper()

	for _, msisdn := range invalidMapperMsisdns {
		hasInfo, info := mapper.GetCarrier(386, msisdn)
		if hasInfo {
			t.Error("Invalid msisdn shouldn't return hasInfo = true")
		} else if info.Name != "" {
			t.Error("Invalid msisdn shouldn't return non-empty info")
		} else if !hasInfo && info.Name != "" {
			t.Error("Invalid msisdn caused hasInfo = false and non-empty info")
		}
	}
}

func TestGetCarrierWithValidMsisdns(t *testing.T) {
	mapper := createTestMapper()

	for _, tt := range validMapperMsisdns {
		hasInfo, info := mapper.GetCarrier(386, tt.msisdn)
		if !hasInfo || tt.carrier != info.Name {
			t.Error("Valid msisdn but no return information")
		}
	}
}
