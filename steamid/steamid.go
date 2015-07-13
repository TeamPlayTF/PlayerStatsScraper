package steamid

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ((Universe << 56) | (Account Type << 52) | (Instance << 32) | Account ID)

func CommIdToSteamId(commid string) (string, error) {
	number, err := strconv.ParseInt(commid, 10, 64)
	if err != nil {
		return "", err
	}

	accountId := (number << 32) >> 32
	instance := ((number >> 32) << 44) >> 44
	accountType := ((number >> 52) << 60) >> 60
	universe := (number >> 56)

	if accountType != 1 {
		return "", errors.New("Non-user ids are not supported")
	}

	return fmt.Sprintf("[U:%v:%v:%v]", universe, accountId, instance), nil

}

func CommIdToLegacySteamId(commid string) (string, error) {
	number, err := strconv.ParseInt(commid, 10, 64)
	if err != nil {
		return "", err
	}

	accountId := (number << 32) >> 32
	universe := (number >> 56)

	if universe != 1 {
		return "", errors.New("Non-public legacy ids are not supported")
	}

	return fmt.Sprintf("STEAM_0:%v:%v", accountId&1, accountId>>1), nil
}

func SteamIdToCommId(steamid string) (string, error) {
	matchedLegacy, _ := regexp.MatchString("STEAM_[0-9]+:[01]:[0-9]+", steamid)
	matchedNew, _ := regexp.MatchString("\\[?[A-Z]+:[0-9]+:[0-9]+[:[0-9]+]?\\]?", steamid)

	if matchedLegacy {
		return legacySteamIdToCommId(steamid)
	} else if matchedNew {
		return newSteamIdToCommId(steamid)
	}

	return "", errors.New("Steam ID didn't match legacy or modern format")
}

func legacySteamIdToCommId(steamid string) (string, error) {
	steamid = strings.Replace(steamid, "STEAM_", "", 1)
	params := strings.Split(steamid, ":")

	if params[0] != "0" {
		return "", errors.New("Non-public legacy ids are not supported")
	}

	Y, _ := strconv.ParseInt(params[1], 10, 64)
	Z, _ := strconv.ParseInt(params[2], 10, 64)

	accountId := Z<<1 + Y

	return strconv.FormatInt((1<<56)|(1<<52)|(1<<32)|accountId, 10), nil
}

func newSteamIdToCommId(steamid string) (string, error) {
	steamid = strings.Replace(steamid, "[", "", 1)
	steamid = strings.Replace(steamid, "]", "", 1)
	params := strings.Split(steamid, ":")

	if params[0] != "U" {
		return "", errors.New("Non-user ids are not supported")
	}

	universe, _ := strconv.ParseInt(params[1], 10, 64)
	accountId, _ := strconv.ParseInt(params[2], 10, 64)

	var instanceId int64 = 1
	if len(params) == 4 {
		instanceId, _ = strconv.ParseInt(params[3], 10, 64)
	}

	return strconv.FormatInt((universe<<56)|(1<<52)|(instanceId<<32)|accountId, 10), nil
}
