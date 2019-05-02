package bigfloalt

import (
	"math"
	"math/big"
)

func Format(num string, base int, decimals int) (string, error) {
	f, _, err := new(big.Float).Parse(num, base)
	if err != nil {
		return "", nil
	}
	value := new(big.Float).Quo(f, big.NewFloat(math.Pow10(decimals)))
	return value.Text('f', -1), nil
}

func Add(x, y string) (string, error) {
	fx, _, err := new(big.Float).Parse(x, 10)
	if err != nil {
		return "", err
	}
	fy, _, err := new(big.Float).Parse(y, 10)
	if err != nil {
		return "", err
	}

	r := new(big.Float).Add(fx, fy)
	return r.Text('f', -1), nil
}

func Sub(x, y string) (string, error) {
	fx, _, err := new(big.Float).Parse(x, 10)
	if err != nil {
		return "", err
	}
	fy, _, err := new(big.Float).Parse(y, 10)
	if err != nil {
		return "", err
	}

	r := new(big.Float).Sub(fx, fy)
	return r.Text('f', -1), nil
}

func Gt(x, y string) (bool, error) {
	fx, _, err := new(big.Float).Parse(x, 10)
	if err != nil {
		return false, err
	}
	fy, _, err := new(big.Float).Parse(y, 10)
	if err != nil {
		return false, err
	}

	return fx.Cmp(fy) == 1, nil
}
