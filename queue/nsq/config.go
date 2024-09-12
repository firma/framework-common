package nsq

type NsqConfig struct {
	NsqdAddr            string
	NsqLookupdAddr      string
	EnableLookup        bool
	LookupdPollInterval int64
	MaxInFlight         int
	AuthSecret          string
}
