package pricer

// The implementation of Coxx Ross Rubinstein model

import (
	"errors"
	"math"
)

type AmericanPricer struct {
	PricingInput
	priceTree  [][]float64
	payoffTree [][]float64
}

func (ap *AmericanPricer) numberOfStep() int {
	nb_of_step_min := 2
	nb_of_step_max := 100
	return min(max(int(365*ap.dT), nb_of_step_min), nb_of_step_max)
}

func (ap *AmericanPricer) dt() float64 {
	return ap.dT / float64(ap.numberOfStep())
}

func (ap *AmericanPricer) up() float64 {
	return math.Exp(ap.volatility * math.Sqrt(ap.dt()))
}

func (ap *AmericanPricer) down() float64 {
	return 1. / ap.up()
}

func (ap *AmericanPricer) qUp() float64 {
	return (math.Exp((ap.riskFreeRate-ap.dividendYield)*ap.dt()) - ap.down()) / (ap.up() - ap.down())
}

func (ap *AmericanPricer) qDown() float64 {
	return 1. - ap.qUp()
}

func (ap *AmericanPricer) incrementalDiscountFactor() float64 {
	return math.Exp(-ap.riskFreeRate * ap.dt())
}

func (ap *AmericanPricer) Price() (float64, error) {
	ap.intitializeTrees()
	ap.setPriceTree()
	ap.initPayoffTree()
	for i := ap.numberOfStep() - 1; i >= 0; i-- {
		previousPayoff := ap.payoffTree[i+1]
		payoff := ap.payoffTree[i]
		previousPayoff_a := previousPayoff[0]
		for j := 0; j < len(payoff); j++ {
			previousPayoff_b := previousPayoff[j+1]
			payoff[j] = (previousPayoff_a*ap.qUp() + previousPayoff_b*ap.qDown()) * ap.incrementalDiscountFactor()
			previousPayoff_a = previousPayoff_b
		}
		ap.checkEarlyExercise(i)
	}
	return ap.payoffTree[0][0], nil
}

func (ap *AmericanPricer) intitializeTrees() {
	var tempPriceTree [][]float64
	var tempPayoffTree [][]float64
	for i := 0; i <= ap.numberOfStep(); i++ {
		var tempPriceNode []float64
		var tempPayoffNode []float64
		for j := 0; j <= i; j++ {
			tempPriceNode = append(tempPriceNode, 0.)
			tempPayoffNode = append(tempPayoffNode, 0.)
		}
		tempPriceTree = append(tempPriceTree, tempPriceNode)
		tempPayoffTree = append(tempPayoffTree, tempPayoffNode)

	}
	ap.priceTree = tempPriceTree
	ap.payoffTree = tempPayoffTree
}

func (ap *AmericanPricer) setPriceTree() {
	ap.priceTree[0][0] = ap.spotPrice
	for i := 0; i < ap.numberOfStep(); i++ {
		previousNodes := ap.priceTree[i]
		currentNodes := ap.priceTree[i+1]
		for j := 0; j <= i; j++ {
			currentNodes[j] = previousNodes[j] * ap.up()
		}
		currentNodes[i+1] = previousNodes[i] * ap.down()
	}
}

func (ap *AmericanPricer) initPayoffTree() error {
	finalPricesNode := ap.priceTree[ap.numberOfStep()]
	finalPayoffNode := ap.payoffTree[ap.numberOfStep()]
	switch ap.optionType {
	case Call:
		for i := 0; i < len(finalPayoffNode); i++ {
			finalPayoffNode[i] = math.Max(0., finalPricesNode[i]-ap.strike)
		}
		return nil
	case Put:
		for i := 0; i < len(finalPayoffNode); i++ {
			finalPayoffNode[i] = math.Max(0., ap.strike-finalPricesNode[i])
		}
		return nil
	default:
		return errors.New("unvalid option type")
	}

}

func (ap *AmericanPricer) checkEarlyExercise(node int) error {
	payoffNode := ap.payoffTree[node]
	pricesNode := ap.priceTree[node]
	switch ap.optionType {
	case Call:
		for i := 0; i < len(payoffNode); i++ {
			payoffNode[i] = math.Max(payoffNode[i], pricesNode[i]-ap.strike)
		}
		return nil
	case Put:
		for i := 0; i < len(payoffNode); i++ {
			payoffNode[i] = math.Max(payoffNode[i], ap.strike-pricesNode[i])
		}
		return nil
	default:
		return errors.New("unvalid option type")
	}
}
