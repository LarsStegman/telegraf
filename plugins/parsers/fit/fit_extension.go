package fit

import (
	"github.com/tormoder/fit"
	"math"
)

func RightContribution(lr fit.LeftRightBalance) float64 {
	if lr == fit.LeftRightBalanceInvalid {
		return math.NaN()
	}

	amount := float64(lr & fit.LeftRightBalanceMask)
	amountIsRight := lr&fit.LeftRightBalanceRight != 0
	if amountIsRight {
		return amount
	} else {
		return 100 - amount
	}
}

func LeftContribution(lr fit.LeftRightBalance) float64 {
	return 100 - RightContribution(lr)
}

func RightContribution100(lr fit.LeftRightBalance100) float64 {
	if lr == fit.LeftRightBalance100Invalid {
		return math.NaN()
	}

	amount := float64(lr&fit.LeftRightBalance100Mask) / 100
	amountIsRight := lr&fit.LeftRightBalance100Right != 0
	if amountIsRight {
		return amount
	} else {
		return 100 - amount
	}
}

func LeftContribution100(lr fit.LeftRightBalance100) float64 {
	return 100 - RightContribution100(lr)
}
