package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const intRange = 4

func (ind Individual) String() (s string) {
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

func genIndividual(configEntries map[string]*ConfigEntry) (ind Individual) {

	for k, v := range configEntries {
		var value ConfigValue

		if k == "Language" {
			continue
		}

		value.entry = k
		if strings.LastIndex(v.Type, "Flags") != -1 {
			value.value = "{"

			for idx, o := range v.Options {
				if idx != 0 {
					value.value += ", "
				}
				if rand.Intn(2) == 0 {
					value.value += o + ": false"
				} else {
					value.value += o + ": true"
				}
			}
			value.value += "}"
		} else {
			switch v.Type {
			case "bool":
				if rand.Intn(2) == 0 {
					value.value = "false"
				} else {
					value.value = "true"
				}
				break
			case "int":
				value.value = fmt.Sprintf("%d", rand.Intn(intRange*2)-intRange)
				break
			case "unsigned":
				value.value = fmt.Sprintf("%d", rand.Intn(intRange))
				break
			default:

				if len(v.Options) != 0 {
					value.value = v.Options[rand.Intn(len(v.Options)-1)]
				} else {
					fmt.Println("Err: no type/:parameter for " + k)
					continue
				}
				break
			}
		}

		ind.configValues = append(ind.configValues, value)
	}

	return ind
}

func (ind *Individual) UpdateScore(clangPath string, perfectSource string) error {
	ind.score = math.MaxUint32
	conf := ind.String()

	out, err := exec.Command(clangPath, "-style=\""+conf+"\"", perfectSource).Output()
	if err != nil {
		return err
	}

	ioutil.WriteFile("/tmp/reformated", out, os.ModePerm)

	out, err = exec.Command("/usr/bin/diff", "/tmp/reformated", perfectSource).Output()
	/*
		if err != nil {
			return err
		}*/
	ind.score = uint32(len(out))
	return nil
}
