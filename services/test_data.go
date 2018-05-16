package services

import "github.com/dekichan/msisdninfo/types"

type TestMsisdns struct {
	Invalid          []string
	ValidSloA1       []string
	ValidSloA1Result types.TransformResponseMsg
}

func GetTestMsisdns() TestMsisdns {
	return TestMsisdns{
		Invalid:          invalidMsisdns,
		ValidSloA1:       validSloA1Msisdns,
		ValidSloA1Result: validSloA1MsisdnResult,
	}
}

var invalidMsisdns = []string{
	"-xxx3443553465",   // invalid characters
	"123456",           // too short
	"1234567890123456", // too long
	"99940123456",      // inexistent country code
	"38629123456",      // inexistent MNO
	"+38642123456",     // inexistent MNO
}

// If you're changing this you might need to change validSloA1MsisdnResult too
var validSloA1Msisdns = []string{
	"38640123456",
	"+38640123456",
	"0038640123456",
}

// If you're changing this you might need to change validSloA1Msisdns too
var validSloA1MsisdnResult = types.TransformResponseMsg{
	CountryCode:       386,
	CountryIdentifier: "SI",
	MnoIdentifier:     "A1",
	SubscriberNumber:  "123456",
}
