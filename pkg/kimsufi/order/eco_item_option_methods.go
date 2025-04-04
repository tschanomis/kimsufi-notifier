package order

import (
	"fmt"
	"slices"
	"strings"
)

// Get returns the option that matches the provided family.
func (i EcoItemOptions) Get(family string) *EcoItemOption {
	for index := range i {
		o := &i[index]
		if o.Family == family {
			return o
		}
	}

	return nil
}

func (i EcoItemOptions) GetMandatoryOptions(filter func(EcoItemOptions, EcoItemOption) bool) EcoItemOptions {
	var options = EcoItemOptions{}

	for _, option := range i {
		if !option.Mandatory {
			continue
		}

		if filter != nil && !filter(options, option) {
			continue
		}

		options = append(options, option)
	}

	return options
}

// GetCheapestMandatoryOptions returns the cheapest mandatory options.
// It compares prices matching PriceDuration and PricingMode.
func (i EcoItemOptions) GetCheapestMandatoryOptions() EcoItemOptions {
	defautPriceConfig := EcoItemPriceConfig{
		Duration:    PriceDuration,
		PricingMode: PricingMode,
	}

	var options = EcoItemOptions{}

	for _, option := range i {
		if !option.Mandatory {
			continue
		}

		current := options.Get(option.Family)
		if current == nil {
			options = append(options, option)
			continue
		}

		newPrice := option.GetPriceByConfig(defautPriceConfig)
		if newPrice == nil {
			continue
		}

		currentPrice := current.GetPriceByConfig(defautPriceConfig)
		if currentPrice == nil {
			continue
		}

		if newPrice.PriceInUcents < currentPrice.PriceInUcents {
			options = append(options, option)
		}
	}

	return options
}

// GetPriceByConfig returns the first price that matches the provided EcoItemPriceConfig.
func (i EcoItemOption) GetPriceByConfig(priceConfig EcoItemPriceConfig) *EcoItemOptionPrice {
	for _, price := range i.Prices {
		if price.Duration == priceConfig.Duration && price.PricingMode == priceConfig.PricingMode {
			return &price
		}
	}

	return nil
}

// ToOptions converts EcoItemOptions to Options.
func (i EcoItemOptions) ToOptions() Options {
	var options []Option

	for _, o := range i {
		options = append(options, o.Option)
	}

	return options
}

// NewOptionsFromMap creates a new Options from a map.
// Using the map keys as family and the values as planCode.
func NewOptionsFromMap(optionsMap map[string]string) Options {
	var options Options

	for family, planCode := range optionsMap {
		o := Option{
			Family:   family,
			PlanCode: planCode,
		}
		options = append(options, o)
	}

	return options
}

// NewOptionsFromSlice creates a new Options from a map.
// Using the map keys as family and the values as planCode.
func NewOptionsFromSlice(optionsSlice []string) (Options, error) {
	var options Options

	for _, o := range optionsSlice {
		oSplit := strings.Split(o, "=")
		if len(oSplit) != 2 {
			return nil, fmt.Errorf("invalid option: %s", o)
		}

		opt := Option{
			Family:   oSplit[0],
			PlanCode: oSplit[1],
		}

		options = append(options, opt)
	}

	return options, nil
}

func (opts Options) SplitByPlanCode(planCode string) (Options, Options) {
	var matching Options
	var other Options

	for _, o := range opts {
		if o.PlanCode == planCode {
			matching = append(matching, o)
		} else {
			other = append(other, o)
		}
	}

	return matching, other
}

func NewOptionsCombinationsFromSlice(optionsSlice Options) []Options {
	var groups = optionsSlice.Groups()
	var total = 1
	for _, options := range groups {
		total *= len(options)
	}

	var combinations = make([]Options, total)
	for _, options := range groups {
		index := 0
		for i := 0; i < total; i++ {
			if i%len(options) == 0 {
				index = 0
			}

			opt := options[index]
			index++

			combinations[i] = append(combinations[i], opt)
		}
	}

	return combinations
}

func (opts Options) Set(opt Option) []Option {
	for i, o := range opts {
		if o.Family == opt.Family {
			opts[i] = opt
			return opts
		}
	}

	return append(opts, opt)
}

// Merge merges the provided Options with the current Options.
// It does not overwrite existing options.
func (opts Options) Merge(other []Option) Options {
	families := opts.Families()

	for _, o := range other {
		if !slices.Contains(families, o.Family) {
			opts = append(opts, o)
		}
	}

	return opts
}

func (opts Options) Families() []string {
	var families []string

	for _, o := range opts {
		families = append(families, o.Family)
	}

	return families
}

func (opts Options) PlanCodes() []string {
	var planCodes []string

	for _, o := range opts {
		planCodes = append(planCodes, o.PlanCode)
	}

	return planCodes
}

func (opts Options) Groups() map[string]Options {
	group := map[string]Options{}

	for _, o := range opts {
		if _, ok := group[o.Family]; !ok {
			group[o.Family] = Options{}
		}

		group[o.Family] = append(group[o.Family], o)
	}

	return group
}
