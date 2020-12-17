package types

import (
	"fmt"
	"testing"
	"time"
)

func TestGetPoolerManagerFromFiles(t *testing.T) {
	m := GetPoolerManagerFromFiles("../../../../template/mnemonic/farm_test", "pooler")
	fmt.Println(len(m))
}

func TestGetPoolerManagerFromFiles1(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
	groupNum := 3
	times := len(arr)/groupNum + 1
	for i := 0; i < times; i++ {
		if i != times-1 {
			fmt.Println(arr[i*groupNum : (i+1)*groupNum])
		} else {
			fmt.Println(arr[i*groupNum:])
		}
	}
}

func TestTimeFormat(t *testing.T) {
	str := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(str)
}
