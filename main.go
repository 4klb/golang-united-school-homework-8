package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

var (
	id                    = flag.String("id", "", "usage")
	operation             = flag.String("operation", "", "usage")
	item                  = flag.String("item", "", "usage")
	filename              = flag.String("fileName", "", "usage")
	errorFlagFile         = "-fileName flag has to be specified"
	OperationMissingError = "-operation flag has to be specified"
	errFlagItem           = "-item flag has to be specified"
	errNoSuchFile         = "open : no such file or directory"
	errSameId             = "Item with id 1 already exists"
	errIdFlagMissing      = "-id flag has to be specified"
)

type Arguments map[string]string

type Item struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func Add(args Arguments, writer io.Writer) error {
	if len(args["item"]) == 0 {
		return errors.New(errFlagItem)
	}

	err := CheckIsFileOpen(args)
	if err != nil {
		return err
	}

	err = CheckItemID(args, writer)
	if err != nil {
		err = WriteOut([]byte(errSameId), writer)
		if err != nil {
			return err
		}
	}

	str := "[" + args["item"] + "]"

	err = ioutil.WriteFile(args["fileName"], []byte(str), 0644)
	if err != nil {
		return err
	}

	return nil
}

func List(args Arguments, writer io.Writer) ([]byte, error) {
	bytes, err := Readfile(args["fileName"], writer)
	if err != nil {
		return nil, err
	}
	err = WriteOut(bytes, writer)
	if err != nil {
		return nil, err
	}

	if len(bytes) == 0 {
		return nil, nil
	}
	return bytes, nil
}

func FindByID(args Arguments, writer io.Writer) error {
	if args["id"] == "" {
		return errors.New(errIdFlagMissing)
	} else {
		jsonData, err := Readfile(args["fileName"], writer)
		if err != nil {
			return err
		}

		unmarshaledItem, err := UnmarshalData(jsonData)
		if err != nil {
			log.Println("err", err)
			return err

		}

		for _, value := range unmarshaledItem {
			if value.Id == args["id"] {
				jsonData, err := MarshalData(value)
				if err != nil {
					return err
				}
				err = WriteOut(jsonData, writer)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func Remove(args Arguments, writer io.Writer) error {
	var item []Item

	if args["id"] == "" {
		return errors.New(errIdFlagMissing)
	} else {
		jsonData, err := Readfile(args["fileName"], writer)
		if err != nil {
			return err
		}

		unmarshaledItem, err := UnmarshalData(jsonData)
		if err != nil {
			return err
		}

		exist := IsIDExist(unmarshaledItem, args)

		if !exist {
			err = WriteOut([]byte("Item with id "+args["id"]+" not found"), writer)
			if err != nil {
				return err
			}
			return nil
		}

		for _, value := range unmarshaledItem {
			if value.Id != args["id"] {
				item = append(item, value)
			}
		}

		jsonData, err = MarshalData(item)
		if err != nil {
			return err
		}

		EditFile(args, jsonData)

		bytes, err := Readfile(args["fileName"], writer)
		if err != nil {
			return err
		}

		err = WriteOut(bytes, writer)
		if err != nil {
			return err
		}
	}
	return nil
}

func IsIDExist(item []Item, args Arguments) bool {
	var exist bool
	for _, value := range item {
		if value.Id == args["id"] {
			exist = true
		}
	}
	return exist
}

func EditFile(args Arguments, data []byte) error {
	os.Remove(args["fileName"])

	err := CreateFile(args["fileName"])
	if err != nil {
		return err
	}

	_, err = Openfile(args["fileName"])
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(args["fileName"], data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func CheckItemID(args Arguments, writer io.Writer) error {
	var idToAdd string
	var id2 int

	// jsonData, err := UnmarshalData([]byte(args["item"]))
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }

	// log.Println(jsonData)

	idToAdd = args["item"][7:8]

	file, err := Openfile(args["fileName"])
	if err != nil {
		log.Println(err)
		return err
	}

	jsonData, err := Readfile(file.Name(), writer)
	if err != nil {
		log.Println(err)
		return err
	}

	id2, err = strconv.Atoi(idToAdd)
	if err != nil {
		log.Println(err)
		return err
	}

	id3, err := strconv.Atoi(string(string(jsonData[8:9]))) // TODO
	if err != nil {
		log.Println(err)
		return err
	}

	if id2 == id3 {
		return errors.New(errSameId)
	}

	return nil
}

func Readfile(filename string, writer io.Writer) ([]byte, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return bytes, err
}

func WriteOut(bytes []byte, writer io.Writer) error {
	_, err := writer.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func UnmarshalData(jsonData []byte) ([]Item, error) {
	var item []Item
	err := json.Unmarshal(jsonData, &item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func MarshalData(data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return jsonData, err
	}
	return jsonData, nil
}

func CheckIsFileOpen(args Arguments) (err error) {
	switch args["fileName"] {
	case "":
		return errors.New(errorFlagFile)
	default:
		err = CreateFile(args["fileName"])
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateFile(filename string) error {
	_, err := Openfile(filename)
	if err != nil {
		return err
	}
	log.Println("file:", filename, "was created")
	return nil
}

func Openfile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	// defer file.Close()
	return file, nil
}

func Perform(args Arguments, writer io.Writer) (err error) {

	if len(args["fileName"]) == 0 {
		return errors.New(errorFlagFile)
	}
	switch args["operation"] {
	case "add":
		err = Add(args, writer)
		return err
	case "list":
		_, err = List(args, writer)
		return err
	case "findById":
		err = FindByID(args, writer)
		return err
	case "remove":
		err = Remove(args, writer)
		return err
	case "":
		return errors.New(OperationMissingError)
	default:
		return errors.New("Operation " + args["operation"] + " not allowed!")
	}
}

func parseArgs() Arguments {
	flag.Parse()
	args := make(Arguments)

	args["fileName"] = *filename
	args["item"] = *item
	args["operation"] = *operation
	args["id"] = *id

	return args
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

// args
// str := `-operation "add" -item "{"id": "1", "email": "email@test.com", "age": 23}" -fileName "users.json"`
// go run main.go -operation "add" -item "" -fileName "users.json"
