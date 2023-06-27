package main

import (
	"fmt"
	"option_pricer/pricer"
)

func main() {

	exerciseStyle := pricer.American // American, European
	optionType := pricer.Call        // Call, Put
	spotPrice := 100.0               // Current value of the spot
	strike := 100.0                  // Strike price of the option
	volatility := 0.2                // Anualized volatility
	riskFreeRate := 0.03             // Risk free rate up to the option maturity
	dividendYield := 0.01            // Continuous dividend yield up to the option option maturity
	dT := 0.5                        // Year to maturity

	pricer := pricer.NewPricer(
		exerciseStyle,
		optionType,
		spotPrice,
		strike,
		volatility,
		riskFreeRate,
		dividendYield,
		dT,
	)

	if pricer != nil {
		fmt.Println("Exercise Style:", exerciseStyle)
		fmt.Println("Option Type:", optionType)
		price, pricing_error := pricer.Price()
		if pricing_error != nil {
			fmt.Println(pricing_error)
		} else {
			fmt.Println("Price:", price)
		}
	} else {
		fmt.Println("Invalid Exercise Style")
	}
}
