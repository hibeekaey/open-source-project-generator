package constants

const (
	FrontendAppMain  = "main"
	FrontendAppAdmin = "admin"
	FrontendAppHome  = "home"
)

type AppsStruct struct {
	Frontend []string
}

var Apps = AppsStruct{
	Frontend: []string{
		FrontendAppMain,
		FrontendAppAdmin,
		FrontendAppHome,
	},
}
