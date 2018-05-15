package services

import (
	"github.com/dekichan/msisdninfo/types"
)

type IMsisdnService interface {
	Parse(msisdn string) (types.TransformResponseMsg, error)
}
