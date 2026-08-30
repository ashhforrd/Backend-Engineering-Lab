package shipping

type Zone string

const (
	ZoneLocal Zone = "LOCAL"
	ZoneRegional Zone = "REGIONAL"
	ZoneInternational Zone = "INTERNATIONAL"
)

type ServiceLevel string

const (
	ServiceStandard ServiceLevel = "STANDARD"
	ServiceExpress  ServiceLevel = "EXPRESS"
)

type Request struct {
	WeightGrams int
	Zone Zone
	Service ServiceLevel
}

type Quote struct {
	BasePrice int64
	WeightSurcharge int64
	ServiceSurcharge int64
	Total int64
}