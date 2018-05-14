package services

import (
	"github.com/dekichan/msisdninfo/types"
	"github.com/nyaruka/phonenumbers"
)

type MsisdnService struct {
	msisdn      string
	phoneNumber *phonenumbers.PhoneNumber
}

func CreateMsisdnService() IMsisdnService {
	return &MsisdnService{}
}

func (msisdnService *MsisdnService) Parse(msisdn string) types.TransformResponseMsg {
	phoneNumber, err := phonenumbers.Parse(msisdn, "SI")

	if err != nil {
		// throw some exception?
	}

	msisdnService.saveAsE164(msisdn)
	msisdnService.phoneNumber = phoneNumber

	return msisdnService.toResponseMsg()
}

func (msisdnService *MsisdnService) saveAsE164(msisdn string) {
	msisdnService.msisdn = msisdn
}

func (msisdnService *MsisdnService) toResponseMsg() types.TransformResponseMsg {
	return types.TransformResponseMsg{
		CountryCode:       msisdnService.phoneNumber.GetCountryCode(),
		CountryIdentifier: "tmp",
		MnoIdentifier:     "mno",
		SubscriberNumber:  "subs",
	}
}
