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
const DATA_DIR = "./carriers/"

type PhonenumberToCarrierMapper struct {
	CountryCode    int
	Msisdn         string
	carrierDataDir string
}

func CreatePhonenumberToCarrierMapper() IPhonenumberToCarrierMapper {
	return &PhonenumberToCarrierMapper{
		carrierDataDir: DATA_DIR,
	}
}

func (mapper *PhonenumberToCarrierMapper) GetCarrier(countryCode int, msisdn string) types.CarrierInfo {
	mapper.CountryCode = countryCode
	mapper.Msisdn = msisdn

	return mapper.getCarrierInfo(countryCode, msisdn)
}

func (mapper *PhonenumberToCarrierMapper) getCarrierInfo(countryCode int, msisdn string) types.CarrierInfo {
	// open carrier data file & split by lines
	filePath := fmt.Sprintf("%s/%d.txt", mapper.carrierDataDir, countryCode)
	inFile, _ := os.Open(filePath)

	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	fmt.Println("in mapper")
	fmt.Println(msisdn)

	for scanner.Scan() {
		line := scanner.Text()
		if lineHasData(line) {
			match, info := processLine(line, msisdn)
			// fmt.Println(match, info)
			if match {
				return info
			}
		}
	}

	return types.CarrierInfo{}
}

func processLine(line string, msisdn string) (bool, types.CarrierInfo) {
	prefix, name := splitToPrefixAndName(line)

	if prefixMatches(prefix, msisdn) {
		return true, types.CarrierInfo{
			Name:   name,
			Prefix: prefix,
		}
	}

	return false, types.CarrierInfo{}
}

func lineHasData(line string) bool {
	matched, _ := regexp.MatchString(DATA_LINE_REGEX, line)
	return matched
}

func splitToPrefixAndName(line string) (prefix string, name string) {
	splitLine := strings.SplitN(line, "|", 2)

	prefix = splitLine[0]
	name = splitLine[1]

	return
}

func prefixMatches(prefix string, msisdn string) bool {
	prefixIdx := len(prefix)
	return prefix == msisdn[:prefixIdx]
}
