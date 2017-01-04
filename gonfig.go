package gonfig

import (
	"encoding/json"
	"os"
	"reflect"
	"strconv"
)

var configurationData interface{} = nil

func GetConf(filename string, configuration interface{}) error {
	var err error = nil

	if configurationData != nil {
		configuration = configurationData
		return nil
	}

	err = getFromJson(filename, configuration)
	if err == nil {
		getFromEnvVariables(configuration)
	}

	configurationData = configuration

	return err
}

func getFromJson(filename string, configuration interface{}) error {

	if len(filename) == 0 {
		return nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		return err
	}

	return nil
}

func getFromEnvVariables(configuration interface{}) {
	typ := reflect.TypeOf(configuration)
	// if a pointer to a struct is passed, get the type of the dereferenced object
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		p := typ.Field(i)
		value := os.Getenv(p.Name)
		if !p.Anonymous && len(value) > 0 {
			// struct
			s := reflect.ValueOf(configuration).Elem()

			if s.Kind() == reflect.Struct {
				// exported field
				f := s.FieldByName(p.Name)
				if f.IsValid() {
					// A Value can be changed only if it is
					// addressable and was not obtained by
					// the use of unexported struct fields.
					if f.CanSet() {
						// change value
						kind := f.Kind()
						if kind == reflect.Int || kind == reflect.Int64 {
							setStringToInt(f, value, 64)
						} else if kind == reflect.Int32 {
							setStringToInt(f, value, 32)
						} else if kind == reflect.Int16 {
							setStringToInt(f, value, 16)
						} else if kind == reflect.Uint || kind == reflect.Uint64 {
							setStringToUInt(f, value, 64)
						} else if kind == reflect.Uint32 {
							setStringToUInt(f, value, 32)
						} else if kind == reflect.Uint16 {
							setStringToUInt(f, value, 16)
						} else if kind == reflect.Bool {
							setStringToBool(f, value)
						} else if kind == reflect.Float64 {
							setStringToFloat(f, value, 64)
						} else if kind == reflect.Float32 {
							setStringToFloat(f, value, 32)
						} else if kind == reflect.String {
							f.SetString(value)
						}
					}
				}
			}
		}
	}
}

func setStringToInt(f reflect.Value, value string, bitSize int) {
	convertedValue, err := strconv.ParseInt(value, 10, bitSize)

	if err == nil {
		if !f.OverflowInt(convertedValue) {
			f.SetInt(convertedValue)
		}
	}
}

func setStringToUInt(f reflect.Value, value string, bitSize int) {
	convertedValue, err := strconv.ParseUint(value, 10, bitSize)

	if err == nil {
		if !f.OverflowUint(convertedValue) {
			f.SetUint(convertedValue)
		}
	}
}

func setStringToBool(f reflect.Value, value string) {
	convertedValue, err := strconv.ParseBool(value)

	if err == nil {
		f.SetBool(convertedValue)
	}
}

func setStringToFloat(f reflect.Value, value string, bitSize int) {
	convertedValue, err := strconv.ParseFloat(value, bitSize)

	if err == nil {
		if !f.OverflowFloat(convertedValue) {
			f.SetFloat(convertedValue)
		}
	}
}
