package services

import (
	"fmt"
	"regexp"

	"github.com/dekichan/msisdninfo/types"
	"github.com/nyaruka/phonenumbers"
)

// first crude validation rule for a phone number
// will match 7 to 15 digits, optionally prefixed with '+' or '00'
const MSISDN_REGEX = `^(\+|00)?[0-9]{7,15}$`

// This service parses msisdn as provided by the user
// It implements IMsisdnService interface
type MsisdnService struct {
	msisdn        string
	msisdnUnpref  string
	phoneNumber   *phonenumbers.PhoneNumber
	carrierMapper IPhonenumberToCarrierMapper
	carrierInfo   types.CarrierInfo
}

type MsisdnError struct {
	Message string
	err     error
}

// Implement error interface
func (err MsisdnError) Error() string {
	var msg string

	if err.err != nil {
		msg = fmt.Sprintf("%s: %s", err.Message, err.err.Error())
	} else {
		msg = err.Message
	}

	return msg
}

// Returns an instance of MsisdnService that
// implements IMsisdnService interface
func CreateMsisdnService() IMsisdnService {
	mapper := CreatePhonenumberToCarrierMapper()

	return &MsisdnService{
		carrierMapper: mapper,
	}
}

// Parse msisdn as provided by the user
// If msisdn is valid, parsed object will contain filled types.TransofrmResponseMsg and err will be null
// If msisdn is invalid, parsed will contain empty types.TransfromResponseMsg and err will be an instance of MsisdnError
func (msisdnService *MsisdnService) Parse(msisdn string) (parsed types.TransformResponseMsg, err error) {
	if !isMsisdnValid(msisdn) {
		// throw some error
		return types.TransformResponseMsg{}, MsisdnError{Message: "Oops, invalid msisdn."}
	}

	msisdnPrefixed := toZeroPrefixed(msisdn)
	msisdnUnprefixed := toUnprefixed(msisdn)
	phoneNumber, err := phonenumbers.Parse(msisdnPrefixed, "SI")

	if err != nil {
		return types.TransformResponseMsg{}, MsisdnError{Message: "Oops, unable to parse msisdn", err: err}
	} else if !phonenumbers.IsValidNumber(phoneNumber) {
		return types.TransformResponseMsg{}, MsisdnError{Message: "Oops, unable to parse msisdn: invalid number"}
	}

	msisdnService.phoneNumber = phoneNumber
	msisdnService.msisdn = msisdnPrefixed
	msisdnService.msisdnUnpref = msisdnUnprefixed

	cc := msisdnService.phoneNumber.GetCountryCode()
	carrierOk, carrierInfo := msisdnService.carrierMapper.GetCarrier(int(cc), msisdnUnprefixed)

	if !carrierOk {
		return types.TransformResponseMsg{}, MsisdnError{Message: "Oops, invalid msisdn: no matching network operator found"}
	}

	msisdnService.carrierInfo = carrierInfo
	return msisdnService.toResponseMsg(), nil
}

// Returns msisdn with '00' prefix
// phonenumbers.Parse() doesn't parse country code if there is no
// prefix - in that case it uses default locale provided
// Ie, for slovenian number: 0038640123456
func toZeroPrefixed(msisdn string) string {
	var prefixed string

	if msisdn[0] == '+' {
		prefixed = "00" + trimLeftChars(msisdn, 1)
	} else if msisdn[:2] != "00" {
		prefixed = "00" + msisdn
	} else {
		prefixed = msisdn
	}

	return prefixed
}

// Returns msisdn without '00' or '+' prefixes
// Ie, for slovenian number: 38640123456
func toUnprefixed(msisdn string) string {
	var unprefixed string

	if msisdn[0] == '+' {
		unprefixed = trimLeftChars(msisdn, 1)
	} else if msisdn[:2] == "00" {
		unprefixed = trimLeftChars(msisdn, 2)
	} else {
		unprefixed = msisdn
	}

	return unprefixed
}

// Trims n chars from the start of string s
// ie trimLeftChars("1234", 2) == "34" is true
func trimLeftChars(s string, n int) string {
	return s[n:]
}

// Returns fields types.TransformResponseMsg object
// types.TransformResponseMsg has the data as needed for final response
func (msisdnService *MsisdnService) toResponseMsg() types.TransformResponseMsg {
	countryCode := msisdnService.phoneNumber.GetCountryCode()
	// map value is an array since country code could be the same for multiple
	// regions. We take the first one since it's usually a (larger) country
	countryIdentifier := phonenumbers.CountryCodeToRegion[int(countryCode)][0]
	subscriberNumber := getSubscriberNumber(msisdnService.carrierInfo.Prefix, msisdnService.msisdnUnpref)

	return types.TransformResponseMsg{
		CountryCode:       countryCode,
		CountryIdentifier: countryIdentifier,
		MnoIdentifier:     msisdnService.carrierInfo.Name,
		SubscriberNumber:  subscriberNumber,
	}
}

// Makes basic msisdn regex validation as defined in MSISDN_REGEX
func isMsisdnValid(msisdn string) bool {
	matched, _ := regexp.MatchString(MSISDN_REGEX, msisdn)

	return matched
}

// Get number without country and carrier codes
// Slovenian example:
// msisdn = 38640123456, subsciberNumber = 123456
func getSubscriberNumber(prefix string, msisdn string) string {
	prefixIdx := len(prefix)
	return msisdn[prefixIdx:]
}
