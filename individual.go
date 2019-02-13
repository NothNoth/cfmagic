package main

import (
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path"
)

type Individual struct {
	configValues []ConfigValue
	score        uint32
	identifier   uint32
}

func (ind Individual) String() (s string) {
	return ind.toOneLineConfig()
}

func (ind Individual) toOneLineConfig() (s string) {
	s = "{"

	for idx, cv := range ind.configValues {
		if idx != 0 {
			s += ", "
		}
		s += cv.String()
	}
	s += "}"
	return s
}

func (ind Individual) toClangFormatConfigFile() (s string) {
	s = "---\n"

	for _, cv := range ind.configValues {
		s += cv.String() + "\n"
	}
	s += "...\n"
	return s
}

func genIndividual(configEntries map[string]*ConfigEntry) (ind Individual) {

	ind.score = math.MaxUint32
	ind.identifier = uint32(rand.Intn(math.MaxUint32))
	for entry := range configEntries {
		cv := generateConfigValue(entry, configEntries)
		if cv != nil {
			ind.configValues = append(ind.configValues, *cv)
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

	out, _ = exec.Command("/usr/bin/diff", fileName, perfectSource).Output()

	//ind.score = uint32(bytes.Count(out, []byte("\n")))
	ind.score = uint32(len(out))
	/*
		if ind.score == 0 {
			fmt.Println("Job done")
			fmt.Println("Formated code at : " + fileName)
			fmt.Println("Clang config at  : " + path.Join(path.Dir(perfectSource), ".clang-format"))
			panic("woot")
		}*/
	return nil
}

func (mother *Individual) mix(father *Individual, mutationRate uint32, configEntries map[string]*ConfigEntry) (baby Individual) {

	for i := 0; i < len(mother.configValues); i++ {
		var motherValue *ConfigValue
		var fatherValue *ConfigValue
		var babyValue ConfigValue
		motherValue = &mother.configValues[i]
		var j int
		for j = 0; j < len(father.configValues); j++ {
			if father.configValues[i].entry == motherValue.entry {
				fatherValue = &father.configValues[i]
				break
			}
		}
		if fatherValue == nil {
			baby.configValues = append(baby.configValues, *motherValue)
			continue
		}

		//Mix mother and father values
		switch motherValue.valueType {
		case String:
			fallthrough
		case Unsigned:
			fallthrough
		case Int:
			fallthrough
		case Bool:
			if rand.Intn(2) == 0 {
				babyValue = *motherValue
			} else {
				babyValue = *fatherValue
			}
		case Flags:
			babyValue.valueType = Flags
			babyValue.entry = motherValue.entry
			babyValue.valueFlags = make(map[string]bool)
			for f := range motherValue.valueFlags {
				if rand.Intn(2) == 0 {
					babyValue.valueFlags[f] = motherValue.valueFlags[f]
				} else {
					babyValue.valueFlags[f] = fatherValue.valueFlags[f]
				}
			}
		}

		//If mutation, override with new value
		if uint32(rand.Intn(100)) <= mutationRate {
			mutatedValue := generateConfigValue(motherValue.entry, configEntries)
			if mutatedValue != nil {
				babyValue = *mutatedValue
			}
		}
		baby.configValues = append(baby.configValues, babyValue)
	}

	baby.identifier = uint32(rand.Intn(math.MaxUint32))
	baby.score = uint32(rand.Intn(math.MaxUint32))

	return
}
