package gosign

type CheckResponse struct {
	ServiceCookie bool
	IP            string
	Principal     string // user name
	Factors       []string
	Realm         string // first factor
}
