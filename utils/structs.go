package utils

type API struct {
	Artists   string
	Locations string
	Dates     string
	Relations string
}

type Artists struct {
	Id           int
	Image        string
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
	Locations    string
	ConcertDates string
	Relations    string
}

type Locations struct {
	Id        int
	Locations []string
	Dates     string
}

type Dates struct {
	Id    int
	Dates []string
}

type Relations struct {
	Id             int
	DatesLocations map[string][]string
}

type InformationPage struct {
	Relations Relations
	Locations Locations
	Dates     Dates
	Artist    Artists
}
