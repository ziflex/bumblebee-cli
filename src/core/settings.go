package core

type Settings struct {
	Prefix    string
	Directory string
}

func NewDefaultSettings() *Settings {
	return &Settings{
		Prefix:    PRIMUSRUN,
		Directory: "/usr/share/applications",
	}
}
