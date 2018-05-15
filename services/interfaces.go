package services

import (
	"github.com/dekichan/msisdninfo/types"
)

type IMsisdnService interface {
	Parse(msisdn string) (types.TransformResponseMsg, error)
}

type IPhonenumberToCarrierMapper interface {
	GetCarrier(countryCode int, msisdn string) types.CarrierInfo
}
