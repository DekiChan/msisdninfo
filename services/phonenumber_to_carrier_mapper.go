package services

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/dekichan/msisdninfo/types"
)

// captures string in format <any number of digits>|<any number of any characters>
// example: 386651|SŽ - Infrastruktura
// const DATA_LINE_REGEX = `^\d+\|(\p{L}|\s|\w|\-)+$`
const DATA_LINE_REGEX = `^\d+\|.+$`

// directory containing carrier data files
const CARRIER_DATA_DIR = "./carriers/"

type PhonenumberToCarrierMapper struct {
	CountryCode    int
	Msisdn         string
	carrierDataDir string
}

// Returns an instance of PhoneNumberToCarrierMapper that
// implements IPhonenumberToCarrierMapper interface
func CreatePhonenumberToCarrierMapper() IPhonenumberToCarrierMapper {
	return &PhonenumberToCarrierMapper{
		carrierDataDir: CARRIER_DATA_DIR,
	}
}

// Returns CarrierInfo object with properties Name and Prefix
// Based on country/region code it opens appropriate file with carrier data
// and tries to find carrier data based on matching prefix in file and given msisdn
func (mapper *PhonenumberToCarrierMapper) GetCarrier(countryCode int, msisdn string) (hasInfo bool, carrierInfo types.CarrierInfo) {
	// open carrier data file & split by lines
	filePath := fmt.Sprintf("%s/%d.txt", mapper.carrierDataDir, countryCode)
	inFile, err := os.Open(filePath)

	if err != nil {
		return false, types.CarrierInfo{}
	}

	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if lineHasData(line) {
			match, info := processLine(line, msisdn)
			if match {
				return true, info
			}
		}
	}

	return false, types.CarrierInfo{}
}

// Reads the line and returns CarrierInfo object if msisdn has prefix
// return variable hasInfo determines whether msisdn has prefix contained in line
// return variable info has carrier info (empty if hasInfo is false)
func processLine(line string, msisdn string) (hasInfo bool, info types.CarrierInfo) {
	hasInfo = false
	info = types.CarrierInfo{}

	prefix, name := splitToPrefixAndName(line)

	if prefixMatches(prefix, msisdn) {
		hasInfo = true
		info = types.CarrierInfo{
			Name:   name,
			Prefix: prefix,
		}
	}

	return
}

// Checks whether a given line has carrier prefix/name data
// ie if it is of form "<any number of digits>|<any number of any characters>"
func lineHasData(line string) bool {
	matched, _ := regexp.MatchString(DATA_LINE_REGEX, line)
	return matched
}

// Splits the line in two parts over the pipe "|" character
// If there's more than one "|" it splits only over the first one
func splitToPrefixAndName(line string) (prefix string, name string) {
	splitLine := strings.SplitN(line, "|", 2)

	prefix = splitLine[0]
	name = splitLine[1]

	return
}

// Checks whether string has a given prefix
// ie msisdn = "12345", prefix = "12" return true
//    msisdn = "12345", prefix = "45" return false
func prefixMatches(prefix string, msisdn string) bool {
	prefixIdx := len(prefix)
	return prefix == msisdn[:prefixIdx]
}
