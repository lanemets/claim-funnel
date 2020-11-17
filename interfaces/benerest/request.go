package benerest

type CreateClaimRequest struct {
	Claim struct {
		Email        string `json:"email"`
		Amount       string `json:"amount"`
		CurrencyCode string `json:"currencyCode"`
		Description  string `json:"description"`
		Status       string `json:"status"`
	} `json:"claim"`
	Profile struct {
		ExternalID  string `json:"externalId"`
		ProfileType string `json:"profileType"`
		Person      struct {
			FirstName  string `json:"firstName"`
			MiddleName string `json:"middleName"`
			LastName   string `json:"lastName"`
			Dob        string `json:"dob"`
			Ein        string `json:"ein"`
		} `json:"person"`
		Company struct {
			BusinessType string `json:"businessType"`
			Name         string `json:"name"`
			Tin          string `json:"tin"`
		} `json:"company"`
		Address struct {
			AddressType string `json:"addressType"`
			Line1       string `json:"line1"`
			Line2       string `json:"line2"`
			City        string `json:"city"`
			State       string `json:"state"`
			PostalCode  string `json:"postalCode"`
			CountryCode string `json:"countryCode"`
		} `json:"address"`
		Phone struct {
			PhoneType   string `json:"phoneType"`
			CountryCode string `json:"countryCode"`
			Number      string `json:"number"`
			Ext         string `json:"ext"`
		} `json:"phone"`
	} `json:"profile"`
}
