package gval

import (
	"fmt"
	"testing"
)

var NaN float64 = -1

type Indicator struct {
	DateTime int64
}

type Bar struct {
	DateTime                         int64
	Open, High, Low, Close, PreClose float64
	Volume, Amount                   float64
}

type MaObj struct {
	Indicator
	Value float64
}

// MA
// Y = (C1 + C2 +...+ Cn)/n
func MA(bars []float64, n int) []float64 {
	mas := make([]float64, 0)
	for i := 0; i < len(bars); i++ {
		if i < n-1 {
			mas = append(mas, NaN)
			continue
		}
		var sumOfClose float64
		for j := n - 1; j >= 0; j-- {
			close := bars[i-j]
			sumOfClose += close
		}

		mas = append(mas, sumOfClose/float64(n))
	}
	return mas
}

var MaCross = NewLanguage(Function("maCross", func(arguments ...interface{}) (interface{}, error) {
	if len(arguments) != 4 {
		return nil, fmt.Errorf("maCross() expects exactly four int argument")
	}
	arr := arguments[0].([]float64)
	crossType := arguments[1].(string)
	short := arguments[2].(int)
	long := arguments[3].(int)

	maShort := MA(arr, short)
	maLong := MA(arr, long)

	fmt.Println(maShort)
	fmt.Println(maLong)

	for i := 1; i < len(maShort); i++ {
		if maShort[i] == NaN || maLong[i] == NaN {
			continue
		}
		if crossType == "golden" && maShort[i-1] < maLong[i-1] && maShort[i] >= maLong[i] {
			return true, nil
		}
		if crossType == "death" && maShort[i-1] > maLong[i-1] && maShort[i] <= maLong[i] {
			return true, nil
		}
	}
	return false, nil
}))

func TestGval(t *testing.T) {
	//marketCap > 1000000 && marketCap <= 2000000 &&
	hasMaCross, err := Evaluate("maCross(arr, crossType, short, long)", map[string]interface{}{
		"marketCap": 1500000,
		"arr":       []float64{10.1, 10.2, 10.3, 10.4, 10.2, 9.8, 9.7, 9.9, 10.1, 10.5, 10.6, 10.9, 10.7, 10.8, 10.6},
		"crossType": "golden",
		"short":     5,
		"long":      10,
	}, MaCross)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("结果：", hasMaCross)
}
