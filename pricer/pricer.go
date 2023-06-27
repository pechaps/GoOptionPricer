package pricer

import (
	"math"
)

type OptionPricer interface {
	ForwardSpot() float64
	DiscountFactor() float64
	Price() (float64, error)
}

type ExerciseStyle string

const (
	American ExerciseStyle = "american"
	European ExerciseStyle = "european"
)

type OptionType string

const (
	Call OptionType = "call"
	Put  OptionType = "put"
)

type PricingInput struct {
	optionType    OptionType
	spotPrice     float64
	strike        float64
	volatility    float64
	riskFreeRate  float64
	dividendYield float64
	dT            float64
}

func (pi *PricingInput) ForwardSpot() float64 {
	return pi.spotPrice * math.Exp((pi.riskFreeRate-pi.dividendYield)*pi.dT)
}

func (pi *PricingInput) DiscountFactor() float64 {
	return math.Exp(-pi.riskFreeRate * pi.dT)
}

func NewPricer(
	exerciseStyle ExerciseStyle,
	optionType OptionType,
	spotPrice float64,
	strike float64,
	volatility float64,
	riskFreeRate float64,
	dividendYield float64,
	dT float64,
) OptionPricer {
	switch exerciseStyle {
	case American:
		return &AmericanPricer{
			PricingInput: PricingInput{
				optionType:    optionType,
				spotPrice:     spotPrice,
				strike:        strike,
				volatility:    volatility,
				riskFreeRate:  riskFreeRate,
				dividendYield: dividendYield,
				dT:            dT,
			},
		}
	case European:
		return &EuropeanPricer{
			PricingInput: PricingInput{
				optionType:    optionType,
				spotPrice:     spotPrice,
				strike:        strike,
				volatility:    volatility,
				riskFreeRate:  riskFreeRate,
				dividendYield: dividendYield,
				dT:            dT,
			},
		}
	default:
		return nil
	}
}
