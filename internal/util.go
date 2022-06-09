package internal

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	KibiByte uint64 = 1024
	MebiByte        = 1024 * KibiByte
	GibiByte        = 1024 * MebiByte
	TebiByte        = 1024 * GibiByte
)

var (
	reThreshold         = regexp.MustCompile(`(?i)^(\d+)\s*([TGMK]i?B)?$`)
	ErrThresholdInvalid = errors.New("threshold invalid")
)

func ParseThreshold(threshold string) (uint64, error) {
	match := reThreshold.FindStringSubmatch(threshold)
	if match == nil {
		return 0, ErrThresholdInvalid
	}

	value, err := strconv.ParseUint(match[1], 10, 64)
	if err != nil {
		return 0, err
	}

	var level uint64

	switch u := strings.ToLower(match[2]); u {
	case "kb", "kib":
		level = value * KibiByte
	case "", "mb", "mib":
		level = value * MebiByte
	case "gb", "gib":
		level = value * GibiByte
	case "tb", "tib":
		level = value * TebiByte
	default:
		return 0, fmt.Errorf("invalid unit")
	}

	return level, nil
}
