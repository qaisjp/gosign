package gosign

type CheckResponse struct {
	IP        string
	Principal string // user name
	Factors   []string
	Realm     string // first factor
}
