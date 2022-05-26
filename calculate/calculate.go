package calculate

import (
	"errors"
	"math"
	"strconv"
)

func toFloat(a, b string) (float64, float64, error) {
	if a == "" {
		a = "0"
	}
	if b == "" {
		b = "0"
	}
	fa, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return 0, 0, err
	}
	fb, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return 0, 0, err
	}
	return fa, fb, nil
}

// Add 加
func Add(a, b string) (string, error) {
	fa, fb, err := toFloat(a, b)
	if err != nil {
		return "", err
	}
	return strconv.FormatFloat(fa+fb, 'f', -1, 64), nil
}

// Sub 减
func Sub(a, b string) (string, error) {
	fa, fb, err := toFloat(a, b)
	if err != nil {
		return "", err
	}
	return strconv.FormatFloat(fa-fb, 'f', -1, 64), nil
}

// Mul 乘
func Mul(a, b string) (string, error) {
	fa, fb, err := toFloat(a, b)
	if err != nil {
		return "", err
	}
	return strconv.FormatFloat(fa*fb, 'f', -1, 64), nil
}

// Div 除 reverse: true a除b(b除以a), false a除以b
func Div(a, b string, reverse bool) (string, error) {
	fa, fb, err := toFloat(a, b)
	if err != nil {
		return "", err
	}
	if reverse {
		fa, fb = fb, fa
	}
	if fb == 0 {
		return "", errors.New("被除数不能为0")
	}
	return strconv.FormatFloat(fa/fb, 'f', -1, 64), nil
}

// Mod 取余
func Mod(a, b string) (string, error) {
	fa, fb, err := toFloat(a, b)
	if err != nil {
		return "", err
	}
	return strconv.FormatFloat(math.Mod(fa, fb), 'f', -1, 64), nil
}

// Abs 绝对值
func Abs(a string) (string, error) {
	fa, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return "", err
	}
	return strconv.FormatFloat(math.Abs(fa), 'f', -1, 64), nil
}

// Floor 向下取整
func Floor(a string) (string, error) {
	fa, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return "", err
	}
	return strconv.FormatFloat(math.Floor(fa), 'f', -1, 64), nil
}

// Ceil 向上取整
func Ceil(a string) (string, error) {
	fa, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return "", err
	}
	return strconv.FormatFloat(math.Ceil(fa), 'f', -1, 64), nil
}

// Round 四舍五入
func Round(a string) (string, error) {
	fa, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return "", err
	}
	return strconv.FormatFloat(math.Round(fa), 'f', -1, 64), nil
}

// Pow 幂运算
func Pow(a, b string) (string, error) {
	fa, fb, err := toFloat(a, b)
	if err != nil {
		return "", err
	}
	return strconv.FormatFloat(math.Pow(fa, fb), 'f', -1, 64), nil
}
