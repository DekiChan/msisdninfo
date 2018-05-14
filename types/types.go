package types

type HomeResponseMsg struct {
	// json:"name" tags tells how to name the prop in json string
	Status  int    `json:"name"`
	Message string `json:"message"`
}

type TransformResponseMsg struct {
	// network operator name
	MnoIdentifier string `json:"mno_identifier"`
	// country calling code (e.g. 386 in case of Slovenian number)
	CountryCode int `json:"country_code"`
	// ISO 3166-1-alpha-2 formatted
	CountryIdentifier string `json:"country_identifier"`
	// number without country and operator codes
	SubscriberNumber string `json:"subscriber_number"`
}

type ErrorResponseMsg struct {
	Message string `json:"message"`
}
