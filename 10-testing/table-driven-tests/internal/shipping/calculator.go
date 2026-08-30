package shipping

import "errors"

var (
	ErrInvalidWeight = errors.New(
		"weight must be between 1 and 20000 grams",
	)

	ErrInvalidZone = errors.New(
		"invalid shipping zone",
	)

	ErrInvalidService = errors.New(
		"invalid service level",
	)
)

func Calculate(request Request) (Quote, error) {
	if request.WeightGrams <= 0 ||
		request.WeightGrams > 20000 {
			return Quote{}, ErrInvalidWeight
		}
	
	basePrice, surchargePerKilogram, err := 
		ratesForZone(request.Zone)
	if err != nil {
		return Quote{}, err
	}

	extraKilograms := calculateExtraKilograms(
		request.WeightGrams,
	)

	weightSurcharge := int64(extraKilograms) * surchargePerKilogram

	subtotal := basePrice + weightSurcharge

	serviceSurcharge, err := calculateServiceSurcharge(
		request.Service,
		subtotal,
	)
	if err != nil {
		return Quote{}, err
	}

	return Quote{
		BasePrice: basePrice,
		WeightSurcharge: weightSurcharge,
		ServiceSurcharge: serviceSurcharge,
		Total: subtotal + serviceSurcharge,
	}, nil
	
}

func ratesForZone(
	zone Zone,
) (int64, int64, error) {
	switch zone {
	case ZoneLocal:
		return 10000, 2000, nil
	
	case ZoneRegional:
		return 20000, 5000, nil

	case ZoneInternational:
		return 50000, 15000, nil

	default:
		return 0, 0, ErrInvalidZone
	}
}

func calculateExtraKilograms(
	weightGrams int,
) int {
	if weightGrams <= 1000 {
		return 0
	}

	extraGrams := weightGrams - 1000

	return (extraGrams + 999) / 1000
}

func calculateServiceSurcharge(
	service ServiceLevel,
	subtotal int64,
) (int64, error) {
	switch service {
	case ServiceStandard:
		return 0, nil

	case ServiceExpress:
		return subtotal / 2, nil

	default:
		return 0, ErrInvalidService
	}
}