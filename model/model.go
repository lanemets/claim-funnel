package model

type ClaimId struct {
	Value string
}

type CreateClaimRequest struct {
	Claim   *Claim
	Profile *Profile
}

type Claim struct {
	Email                 string
	Amount                string
	CurrencyCode          string
	ClientReferenceNumber string
	Description           string
	CreatedOn             string
	ExpiresOn             string
}

type Profile struct {
	ExternalId    string
	ProfileType   ProfileType
	ProfileEntity ProfileEntity
	Address       *Address
	Phone         *Phone
}

type Person struct {
	FirstName  string
	MiddleName string
	LastName   string
	Dob        string
	Ein        string
}

type ProfileEntity interface {
	ProfileEntity()
}

func (*ProfilePerson) ProfileEntity() {}

func (*ProfileCompany) ProfileEntity() {}

type ProfilePerson struct {
	Person *Person
}

type ProfileCompany struct {
	Company *Company
}

type CompanyType int32

const (
	UNKNOWN         CompanyType = 0
	CORPORATION     CompanyType = 1
	PARTNERSHIP     CompanyType = 2
	GOVERNMENT      CompanyType = 3
	NONPROFIT       CompanyType = 4
	PUBLIC_COMPANY  CompanyType = 5
	PRIVATE_COMPANY CompanyType = 6
)

type Company struct {
	Id           string
	BusinessType CompanyType
	Name         string
	Tin          string
}

type ProfileType = int32

const (
	PERSON  ProfileType = 0
	COMPANY ProfileType = 1
)

type AddressType = int32

const (
	RESIDENTIAL AddressType = 0
	COMMERCIAL  AddressType = 1
)

type PhoneType int32

const (
	HOME  PhoneType = 0
	CELL  PhoneType = 1
	WORK  PhoneType = 2
	OTHER PhoneType = 3
)

type Address struct {
	AddressType AddressType
	Line1       string
	Line2       string
	City        string
	State       string
	PostalCode  string
	CountryCode string
}

type Phone struct {
	PhoneType   PhoneType
	CountryCode string
	Number      string
	Ext         string
}
