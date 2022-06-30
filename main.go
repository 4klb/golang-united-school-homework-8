package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	id                    = flag.String("id", "", "usage")
	operation             = flag.String("operation", "", "usage")
	item                  = flag.String("item", "", "usage")
	filename              = flag.String("fileName", "", "usage")
	errorFlagFile         = "-fileName flag has to be specified"
	OperationMissingError = "-operation flag has to be specified"
	errNoSuchFile         = "open : no such file or directory"
)

type Arguments map[string]string

func Add(args Arguments) error {
	if len(args["item"]) == 0 {
		return errors.New("")
	}
	err := Openfile(args)
	if err != nil {
		return err
	}

	result, err := MarshalData(args)
	ioutil.WriteFile(args["fileName"], result, 0644)
	if err != nil {
		return err
	}
	return nil
}

func List(name string, writer io.Writer) error {
	bytes, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	if len(bytes) == 0 {
		return nil
	}

	writer.Write(bytes)

	return nil
}

func FindByID() error {
	return nil
}

func DeleteList() error {
	return nil
}

func MarshalData(args Arguments) ([]byte, error) {
	js, err := json.Marshal(args)
	if err != nil {
		return js, err
	}
	return js, nil
}

// func UnmarshalData(js []byte) (Args, error) {
// 	args := Args{}
// 	err := json.Unmarshal(js, &args)
// 	if err != nil {
// 		return args, err
// 	}
// 	return args, err
// }

func Openfile(args Arguments) (err error) {
	switch args["fileName"] {
	case "":
		return errors.New(errorFlagFile)
	default:
		err = CreateFile(args["fileName"])
	}
	return nil
}

func CreateFile(filename string) error {
	_, err := os.OpenFile(filename, os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	log.Println("file:", filename, "was created")
	return nil
}

func Perform(args Arguments, writer io.Writer) (err error) {
	if len(args["fileName"]) == 0 {
		return errors.New(errorFlagFile)
	}
	switch args["operation"] {
	case "add":
		err = Add(args)
		return err
	case "list":
		err = List(args["fileName"], writer)
		return err
	case "findById":
		err = FindByID()
		return err
	case "remove":
		err = DeleteList()
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
	// str := `-operation "add" -item "{"id": "1", "email": "email@test.com", "age": 23}" -fileName "users.json"`
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
