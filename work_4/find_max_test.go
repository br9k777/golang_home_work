package findMax

import (
	"fmt"
	"os"
	"testing"

	"go.uber.org/zap"
)

//не разобрался еще как лучше включать dev логирование для тестов
func TestSetLogger(t *testing.T) {
	var err error
	if log, err = zap.NewDevelopment(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func TestGetMax(t *testing.T) {
	// t.Skip()

	s := []int{1, 2, 3}

	fmt.Printf("slice %v Наибольший элемент %#v\n", s, FindMax(s, func(i, j int) bool { return s[i] > s[j] }))
	s2 := []string{"qw", "sdq", "ersc"}

	fmt.Printf("slice %v Наибольший элемент %#v\n", s2, FindMax(s2, func(i, j int) bool { return len(s2[i]) > len(s2[j]) }))

}
