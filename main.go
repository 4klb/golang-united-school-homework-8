package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	id        = flag.String("id", "", "usage")
	operation = flag.String("operation", "", "usage")
	item      = flag.String("item", "", "usage")
	filename  = flag.String("fileName", "", "usage")
	errorFile = ": no such file or directory"
)

type Arguments map[string]string

// type Args struct {
// 	FileName  string `json:"fileName"`
// 	Item      string `json:"item"`
// 	Operation string `json:"operation"`
// }

func Add(args Arguments, name string) error {
	file, err := Openfile(args, name)
	if err != nil {
		return err
	}

	result, err := MarshalData(args)
	ioutil.WriteFile(file.Name(), result, 0644)
	if err != nil {
		return err
	}
	return nil
}

func List(name string, writer io.Writer) error {
	// var args Args
	bytes, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	if len(bytes) == 0 {
		return nil
	}

	writer.Write(bytes)
	// fmt.Println(args.Item)

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

func Openfile(args Arguments, filename string) (*os.File, error) {
	file, err := os.Open(filename)
	if err != nil && err != errors.New("open "+filename+errorFile) {
		file, err = CreateFile(filename)
	} else if err != nil {
		return file, err
	}
	return file, nil
}

func CreateFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_CREATE, 0644)
	if err != nil {
		return file, err
	}
	log.Println("file:", filename, "was created")
	return file, nil
}

func OperationMissingError(hasOperation, hasWrongOperation bool, wrongOperation string) error {
	if !hasOperation {
		return errors.New("-operation flag has to be specified")
	}
	if hasWrongOperation {
		return errors.New("Operation " + wrongOperation + " not allowed!")
	}
	return nil
}

func Perform(args Arguments, writer io.Writer) (err error) {
	// var fileN string
	// var flags []string = []string{"id", "operation", "item", "fileName"}
	// var operations []string = []string{"add", "list", "findById", "remove"}

	// var hasOperation bool
	// var hasWrongOperation bool
	// var wrongOperation string

	// for key, value := range args {
	// 	if key == flags[3] {
	// 		fileN = value
	// 		args["fileName"] = *filename
	// 	} else if key == flags[2] {
	// 		args["item"] = *item
	// 	} else if key == flags[1] {
	// 		if value == "" {
	// 			hasOperation = false
	// 		}
	// 		args["operation"] = *operation
	// 	} else if key == flags[0] {
	// 		args["id"] = *id
	// 	} else if value != operations[0] || value != operations[1] || value != operations[2] || value != operations[3] {
	// 		hasWrongOperation = true
	// 		hasOperation = true
	// 		wrongOperation = value
	// 	}
	// }
	// err = OperationMissingError(hasOperation, hasWrongOperation, wrongOperation)
	// if err != nil {
	// 	return err
	// }

	// log.Println(hasOperation)

	// for _, value := range args[] {
	// 	switch value {
	// 	case "add":
	// 		err = Add(args, fileN)
	// 	case "list":
	// 		err = List(fileN, writer)
	// 	case "findById":
	// 		err = FindByID()
	// 	case "remove":
	// 		err = DeleteList()

	// 	default:
	// return errors.New("-operation flag has to be specified")
	// 		// return errors.New("Operation " + key + " not allowed!")
	// 	}
	// 	if err != nil {
	// 		log.Println(err)
	// 		return err
	// 	}
	// }
	fmt.Println(args)
	switch args["operation"] {
	case "add":
		err = Add(args, args["fileName"])
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
		return errors.New("-operation flag has to be specified")
	default:
		return errors.New("Operation " + args["operarion"] + " not allowed!")
	}
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }
	// // for key, _ := range args {
	// 	if key != flags[0] || key != flags[1] || key != flags[2] {

	// 	}
	// }

	return nil
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
