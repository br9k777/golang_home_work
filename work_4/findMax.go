package findMax

import (
	"fmt"
	"os"
	"reflect"

	"go.uber.org/zap"
)

var (
	log *zap.Logger
)

func init() {
	var err error
	if log, err = zap.NewProduction(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	zap.ReplaceGlobals(log)
}

var reflectValueOf = reflect.ValueOf

// var reflectSwapper = reflect.Swapper

//FindMax поиск максимального элемента используя функцию
func FindMax(slice interface{}, max func(i, j int) bool) interface{} {

	rv := reflectValueOf(slice)

	if slice == nil || rv.IsNil() || rv.Kind() != reflect.Slice {
		return slice
	}

	// swap := reflectSwapper(slice)
	// var getFunc func
	length := rv.Len()
	// log.Info("Массив ", zap.Int(`размер`, length))
	switch length {
	case 0:
		return nil
	case 1:
		return rv.Index(0)
	}

	// quickSort_func(lessSwap{max, swap}, 0, length, maxDepth(length))
	var maxElement int
	for i := 1; i < length; i++ {
		if max(i, maxElement) {
			maxElement = i
		}
	}
	return rv.Index(maxElement)
}
