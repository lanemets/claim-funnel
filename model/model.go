package model

type ClaimId struct {
	Value string
}

type ProcessDefinitionId struct {
	Value string
}

type Address struct {
	AddressType string
	Line1       string
	Line2       string
	City        string
	State       string
	PostalCode  string
	CountryCode string
}

type Company struct {
	BusinessType string
	Name         string
	Tin          string
}

type Phone struct {
	PhoneType   string
	CountryCode string
	Number      string
	Ext         string
}

type Person struct {
	FirstName  string
	MiddleName string
	LastName   string
	Dob        string
	Ein        string
}

type Profile struct {
	ExternalID  string
	ProfileType string
	Person      Person
	Company     Company
	Address     Address
	Phone       Phone
}

type Claim struct {
	Email        string
	Amount       string
	CurrencyCode string
	//TODO: deprecated, should be clientId in new version of claim service
	ClientReferenceNumber string
	Description           string
	Status                string
}
