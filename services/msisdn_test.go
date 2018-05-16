package services

import (
	"reflect"
	"testing"

	"github.com/dekichan/msisdninfo/types"
)

// !!!!! SEE ALSO FILE msisd_test.go where some data is defined

// Needed because mapper otherwise gets wrong carriers data dir
// See comment at createTestMapper()
func createTestMsisdnService() IMsisdnService {
	mapper := createTestMapper()

	return &MsisdnService{
		carrierMapper: mapper,
	}
}

func TestCreateRegularMsisdnService(t *testing.T) {
	createdService := CreateMsisdnService()
	service := &MsisdnService{}

	if reflect.TypeOf(service) != reflect.TypeOf(createdService) {
		t.Error("Can't create regular MsisdnService")
	}
}

func TestParseInvalidMsisdns(t *testing.T) {
	msisdnService := createTestMsisdnService()

	for _, msisdn := range invalidMsisdns {
		resp, err := msisdnService.Parse(msisdn)
		if err == nil {
			t.Error("MsisdnService shouldn't return nil error for invalid msisdn")
		} else if !isResponseEmpty(resp) {
			t.Error("MsisdnService should return empty response object for invalid msisdn")
		}
	}
}

func TestParseValidMsisdns(t *testing.T) {
	msisdnService := createTestMsisdnService()

	for _, msisdn := range validSloA1Msisdns {
		resp, err := msisdnService.Parse(msisdn)
		if err != nil {
			t.Error("MsisdnService should return nil error for valid msisdn")
		} else if !isResponseValid(resp, validSloA1MsisdnResult) {
			t.Error("MsisdnService.Parse result is wrong.")
		}
	}
}

// Checks whether all types.TransformResponseMsg fields are empty
// ie set to default values
func isResponseEmpty(r types.TransformResponseMsg) bool {
	return r.CountryCode == 0 &&
		r.CountryIdentifier == "" &&
		r.MnoIdentifier == "" &&
		r.SubscriberNumber == ""
}

// Checks whether all fields in r match those in expected
func isResponseValid(r types.TransformResponseMsg, expected types.TransformResponseMsg) bool {
	return r.CountryCode == expected.CountryCode &&
		r.CountryIdentifier == expected.CountryIdentifier &&
		r.MnoIdentifier == expected.MnoIdentifier &&
		r.SubscriberNumber == expected.SubscriberNumber
}
