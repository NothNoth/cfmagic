package main

import (
	"fmt"
	"math/rand"
	"strings"
)

const intRange = 120

var ignoredSettings = map[string]bool{
	"Language":                       true,
	"DisableFormat":                  true,
	"BreakAfterJavaFieldAnnotations": true,
}

type ConfigEntryType int

const (
	String ConfigEntryType = iota
	Int
	Unsigned
	Bool
	Flags
)

type ConfigValue struct {
	entry     string
	valueType ConfigEntryType

	valueString   string
	valueBool     bool
	valueInt      int
	valueUnsigned uint32
	valueFlags    map[string]bool
}

func (cv ConfigValue) String() (s string) {
	switch cv.valueType {
	case String:
		return cv.entry + ": " + cv.valueString
	case Int:
		return cv.entry + ": " + fmt.Sprintf("%d", cv.valueInt)
	case Unsigned:
		return cv.entry + ": " + fmt.Sprintf("%d", cv.valueUnsigned)
	case Bool:
		return cv.entry + ": " + fmt.Sprintf("%t", cv.valueBool)
	case Flags:
		s = cv.entry + ": {"
		idx := 0
		for f, v := range cv.valueFlags {
			if idx != 0 {
				s += ", "
			}
			idx++

			s += f + ": " + fmt.Sprintf("%t", v)
		}
		s += "}"
		return s
	}

	return ""
}

func generateConfigValue(entry string, configEntries map[string]*ConfigEntry) *ConfigValue {
	var cv ConfigValue
	cv.entry = entry

	found, _ := ignoredSettings[entry]
	if found == true {
		return nil
	}

	v := configEntries[entry]
	if strings.LastIndex(v.Type, "Flags") != -1 {
		cv.valueType = Flags
		cv.valueFlags = make(map[string]bool)
		for _, f := range v.Options {

			if rand.Intn(2) == 0 {
				cv.valueFlags[f] = false
			} else {
				cv.valueFlags[f] = true
			}
		}
		return &cv
	}

	switch v.Type {
	case "bool":
		cv.valueType = Bool

		if rand.Intn(2) == 0 {
			cv.valueBool = false
			return &cv
		} else {
			cv.valueBool = true
			return &cv
		}
	case "int":
		cv.valueType = Int
		cv.valueInt = rand.Intn(intRange*2) - intRange
		return &cv
	case "unsigned":
		cv.valueType = Unsigned
		cv.valueUnsigned = uint32(rand.Intn(intRange))
		return &cv
	default:
		if len(v.Options) != 0 {
			cv.valueType = String
			cv.valueString = v.Options[rand.Intn(len(v.Options))]
			return &cv
		} else {
			fmt.Println("Err: no type/:parameter for " + entry)
			return nil
		}
	}

	return nil
}
