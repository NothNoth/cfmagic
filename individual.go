package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strings"
)

const intRange = 120

func (ind Individual) String() (s string) {
	return ind.toOneLineConfig()
}

func (ind Individual) toOneLineConfig() (s string) {
	s = "{"

	for idx, cv := range ind.configValues {
		if idx != 0 {
			s += ", "
		}
		s += fmt.Sprint(cv.entry + ": " + cv.value)
	}
	s += "}"
	return s
}

func (ind Individual) toClangFormatConfigFile() (s string) {
	s = "---\n"

	for _, cv := range ind.configValues {
		s += cv.entry + ": " + cv.value + "\n"
	}
	s += "...\n"
	return s
}
func generateConfigValue(entry string, configEntries map[string]*ConfigEntry) string {

	if entry == "Language" {
		return ""
	}
	if entry == "DisableFormat" {
		return ""
	}
	if entry == "BreakAfterJavaFieldAnnotations" {
		return ""
	}

	if entry == "IndentWidth" {
		return "4"
	}

	if entry == "BreakBeforeBraces" {
		return "Custom"
	}

	v := configEntries[entry]
	if strings.LastIndex(v.Type, "Flags") != -1 {
		value := "{"
		for idx, o := range v.Options {
			if idx != 0 {
				value += ", "
			}
			if rand.Intn(2) == 0 {
				value += o + ": false"
			} else {
				value += o + ": true"
			}
		}
		value += "}"
		return value
	}

	switch v.Type {
	case "bool":
		if rand.Intn(2) == 0 {
			return "false"
		} else {
			return "true"
		}
		break
	case "int":
		return fmt.Sprintf("%d", rand.Intn(intRange*2)-intRange)
		break
	case "unsigned":
		return fmt.Sprintf("%d", rand.Intn(intRange))
		break
	default:
		if len(v.Options) != 0 {
			return v.Options[rand.Intn(len(v.Options)-1)]
		} else {
			fmt.Println("Err: no type/:parameter for " + entry)
			return ""
		}
	}

	return ""
}

func genIndividual(configEntries map[string]*ConfigEntry) (ind Individual) {

	ind.score = math.MaxUint32
	ind.identifier = uint32(rand.Intn(math.MaxUint32))
	for entry := range configEntries {
		value := generateConfigValue(entry, configEntries)
		if len(value) != 0 {
			ind.configValues = append(ind.configValues, ConfigValue{entry: entry, value: value})
		}
	}

	return ind
}

func (ind *Individual) UpdateScore(clangPath string, perfectSource string) error {
	ind.score = math.MaxUint32
	conf := ind.toClangFormatConfigFile()

	//fileName := fmt.Sprintf("/tmp/reformated_%d", ind.identifier)
	fileName := "/tmp/reformated"

	ioutil.WriteFile(path.Join(path.Dir(perfectSource), ".clang-format"), []byte(conf), os.ModePerm)

	out, err := exec.Command(clangPath, "-style=file", perfectSource).Output()
	if err != nil {
		return err
	}

	ioutil.WriteFile(fileName, out, os.ModePerm)

	out, err = exec.Command("/usr/bin/diff", fileName, perfectSource).Output()
	/*
		if err != nil {
			return err
		}*/

	ind.score = uint32(bytes.Count(out, []byte("\n")))

	return nil
}

func (mother *Individual) mix(father *Individual, configEntries map[string]*ConfigEntry) (baby Individual) {
	for i := 0; i < len(mother.configValues); i++ {
		var entry string
		var value string
		entry = mother.configValues[i].entry

		if rand.Intn(2) == 0 {
			value = mother.configValues[i].value
		} else {
			for j := 0; j < len(father.configValues); j++ {
				if father.configValues[i].entry == entry {
					value = father.configValues[i].value
					break
				}
			}
			if len(value) == 0 {
				value = mother.configValues[i].value
			}
		}

		if uint32(rand.Intn(100)) <= mutationRate {
			mutatedValue := generateConfigValue(entry, configEntries)
			if len(mutatedValue) != 0 {
				value = mutatedValue
			}
		}
		baby.configValues = append(baby.configValues, ConfigValue{entry: entry, value: value})
	}

	baby.identifier = uint32(rand.Intn(math.MaxUint32))
	baby.score = uint32(rand.Intn(math.MaxUint32))

	return
}
