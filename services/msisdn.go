package services

import (
	"fmt"
	"regexp"

	"github.com/dekichan/msisdninfo/types"
	"github.com/nyaruka/phonenumbers"
)

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

// will match 7 to 15 digits, optionally prefixed with '+' or '00'
const MSISDN_REGEX = `^(\+|00)?[0-9]{7,15}$`

func (err MsisdnError) Error() string {
	var msg string

	if err.err != nil {
		msg = fmt.Sprintf("%s: %s", err.Message, err.err.Error())
	} else {
		msg = err.Message
	}

	return msg
}

func (err MsisdnError) ToResponseError() types.ErrorResponseMsg {
	return types.ErrorResponseMsg{err.err.Error()}
}

func CreateMsisdnService() IMsisdnService {
	mapper := CreatePhonenumberToCarrierMapper()

	return &MsisdnService{
		carrierMapper: mapper,
	}
}

func (msisdnService *MsisdnService) Parse(msisdn string) (types.TransformResponseMsg, error) {
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
	fmt.Println(fmt.Sprintf("Local msisdn unpref: %s", msisdnService.msisdn))

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

func trimLeftChars(s string, n int) string {
	c := 0
	for i := range s {
		if c >= n {
			return s[i:]
		}
		c++
	}
	return s[:0]
}

func (msisdnService *MsisdnService) toResponseMsg() types.TransformResponseMsg {
	fmt.Println(msisdnService.phoneNumber)
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
