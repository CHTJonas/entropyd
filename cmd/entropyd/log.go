package main

import "fmt"

type logTuple struct {
	key   string
	value string
}

func logString(key, value string) *logTuple {
	return &logTuple{
		key,
		value,
	}
}

func logInt(key string, value int) *logTuple {
	return logString(key, fmt.Sprint(value))
}

func logError(key string, err error) *logTuple {
	return logString(key, err.Error())
}

func log(msg string, tuples ...*logTuple) {
	fmt.Printf("msg=%s", msg)
	for _, v := range tuples {
		fmt.Printf(", %s=%s", (*v).key, (*v).value)
	}
	fmt.Println()
}
