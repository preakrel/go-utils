package calculate

import (
	"testing"
)

func TestAdd(t *testing.T) {
	t.Log(Add("1", "2"))
	t.Log(Add("1", "-2"))
	t.Log(Add("1.5", "2.8"))
	t.Log(Sub("1.1", "2.2"))
	t.Log(Sub("5.2", "2.2"))
	t.Log(Sub("10", "7"))
}

func TestDiv(t *testing.T) {
	t.Log(Mul("10", "2"))
	t.Log(Mul("2.2", "-2"))
	t.Log(Mul("-2.8", "-2"))
	t.Log(Div("10", "2", false))
	t.Log(Div("10", "3", false))
	t.Log(Div("10", "2.5", false))
	t.Log(Div("10", "1.8", true))
	t.Log(Div("10", "0", true))
	t.Log(Div("10", "0", false))
}

func TestMod(t *testing.T) {
	t.Log(Mod("10", "3"))
	t.Log(Mod("10", "-3"))
	t.Log(Mod("-10", "3"))
	t.Log(Mod("-10", "-3"))
}

func TestAbs(t *testing.T) {
	t.Log(Abs("103"))
	t.Log(Abs("10.3"))
	t.Log(Abs("-10.3"))
	t.Log(Abs("-103"))
}

func TestFloor(t *testing.T) {
	t.Log(Ceil("1.2"))
	t.Log(Ceil("-12.8"))
	t.Log(Floor("15"))
	t.Log(Floor("15.8"))
	t.Log(Floor("-15.8"))
}

func TestRound(t *testing.T) {
	t.Log(Round("11.2"))
	t.Log(Round("11.8"))
	t.Log(Round("-11.4"))
	t.Log(Round("-11.5"))
}

func TestPow(t *testing.T) {
	t.Log(Pow("1.2", "2"))
	t.Log(Pow("3", "5"))
	t.Log(Pow("6", "6"))
}
