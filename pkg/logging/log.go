package logging

import "fmt"

type LogTuple struct {
	key   string
	value string
}

func LogString(key, value string) *LogTuple {
	return &LogTuple{
		key,
		value,
	}
}

func LogInt(key string, value int) *LogTuple {
	return LogString(key, fmt.Sprint(value))
}

func LogError(key string, err error) *LogTuple {
	return LogString(key, err.Error())
}

func Log(msg string, tuples ...*LogTuple) {
	fmt.Printf("msg=%s", msg)
	for _, v := range tuples {
		fmt.Printf(", %s=%s", (*v).key, (*v).value)
	}
	fmt.Println()
}
