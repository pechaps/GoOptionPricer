package pricer

 // The implementation of Black Scholes model.

import (
	"errors"
	"math"

	"github.com/chobie/go-gaussian"
)

type EuropeanPricer struct {
	PricingInput
}

func (ep *EuropeanPricer) Price() (float64, error) {
	switch ep.optionType {
	case Call:
		return ep.callPrice(), nil
	case Put:
		return ep.putPrice(), nil
	default:
		return 0., errors.New("unvalid option type")
	}
}

func (ep *EuropeanPricer) callPrice() float64 {
	return ep.DiscountFactor() * (ep.NdPlus()*ep.ForwardSpot() - ep.NdMinus()*ep.strike)
}

func (ep *EuropeanPricer) putPrice() float64 {
	return ep.callPrice() - ep.DiscountFactor()*(ep.ForwardSpot()-ep.strike)
}

func (EuropeanPricer *EuropeanPricer) dPlus() float64 {
	return (math.Log(EuropeanPricer.ForwardSpot()/EuropeanPricer.strike) + 0.5*math.Pow(EuropeanPricer.volatility, 2)*EuropeanPricer.dT) / (EuropeanPricer.volatility * math.Sqrt(EuropeanPricer.dT))
}

func (EuropeanPricer *EuropeanPricer) dMinus() float64 {
	return EuropeanPricer.dPlus() - EuropeanPricer.volatility*math.Sqrt(EuropeanPricer.dT)
}

func (EuropeanPricer *EuropeanPricer) NdPlus() float64 {
	return StandardNormalCdf(EuropeanPricer.dPlus())
}

func (EuropeanPricer *EuropeanPricer) NdMinus() float64 {
	return StandardNormalCdf(EuropeanPricer.dMinus())
}

func StandardNormalCdf(x float64) float64 {
	return gaussian.NewGaussian(0, 1).Cdf(x)
}
