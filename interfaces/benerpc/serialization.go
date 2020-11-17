package benerpc

import (
	grpcClient "github.com/lanemets/claim-funnel/external/claim/gen"
	"github.com/lanemets/claim-funnel/model"
)

func toClaimModel(claim *grpcClient.Claim) *model.Claim {
	return &model.Claim{
		Email:        claim.Email,
		Amount:       claim.Amount,
		CurrencyCode: claim.CurrencyCode,
		Description:  claim.Description,
		Status:       claim.Status.String(),
	}
}

func fromClaimModel(claim *model.Claim) *grpcClient.Claim {
	return &grpcClient.Claim{
		Email:                 claim.Email,
		Amount:                claim.Amount,
		CurrencyCode:          claim.CurrencyCode,
		ClientReferenceNumber: claim.ClientReferenceNumber,
		Description:           claim.Description,
		Status:                grpcClient.Claim_StatusType(ClaimStatusTypeReversed[claim.Status]),
	}
}

func toProfileModel(profile *grpcClient.Profile) *model.Profile {
	return &model.Profile{
		ExternalID:  profile.ExternalId,
		ProfileType: profile.ProfileType.String(),
		Person:      model.Person{},
		Company:     model.Company{},
		Address:     model.Address{},
		Phone:       model.Phone{},
	}
}

func toPersonModel(person *grpcClient.Person) *model.Person {
	return &model.Person{
		FirstName:  person.FirstName,
		MiddleName: person.MiddleName,
		LastName:   person.LastName,
		Dob:        person.Dob,
		Ein:        person.Ein,
	}
}

func fromPersonModel(person *model.Person) *grpcClient.Person {
	return &grpcClient.Person{
		FirstName:  person.FirstName,
		MiddleName: person.MiddleName,
		LastName:   person.LastName,
		Dob:        person.Dob,
		Ein:        person.Ein,
	}
}

func fromProfileModel(profile *model.Profile) *grpcClient.Profile {

	profileType := grpcClient.Profile_ProfileType(ProfileTypeReversed[profile.ProfileType])

	//TODO: refactor
	if profileType == 0 {
		return &grpcClient.Profile{
			ExternalId:  profile.ExternalID,
			ProfileType: profileType,
			Entity:      &grpcClient.Profile_Person{Person: fromPersonModel(&profile.Person)},
			Address:     fromAddressModel(&profile.Address),
			Phone:       fromPhoneModel(&profile.Phone),
		}
	} else {
		return &grpcClient.Profile{
			ExternalId:  profile.ExternalID,
			ProfileType: profileType,
			Entity:      &grpcClient.Profile_Company{Company: fromCompanyModel(&profile.Company)},
			Address:     fromAddressModel(&profile.Address),
			Phone:       fromPhoneModel(&profile.Phone),
		}
	}
}

func fromPhoneModel(phone *model.Phone) *grpcClient.Phone {
	return &grpcClient.Phone{
		PhoneType:   grpcClient.Phone_PhoneType(PhoneTypeReversed[phone.PhoneType]),
		CountryCode: phone.CountryCode,
		Number:      phone.Number,
		Ext:         phone.Ext,
	}
}

func fromAddressModel(address *model.Address) *grpcClient.Address {
	return &grpcClient.Address{
		AddressType: grpcClient.Address_AddressType(AddressTypeReversed[address.AddressType]),
		Line1:       address.Line1,
		Line2:       address.Line2,
		City:        address.City,
		State:       address.State,
		PostalCode:  address.PostalCode,
		CountryCode: address.CountryCode,
	}
}

func toCompanyModel(company *grpcClient.Company) *model.Company {
	return &model.Company{
		BusinessType: company.BusinessType.String(),
		Name:         company.Name,
		Tin:          company.Tin,
	}
}

func fromCompanyModel(company *model.Company) *grpcClient.Company {
	return &grpcClient.Company{
		BusinessType: grpcClient.Company_CompanyType(CompanyTypeReversed[company.BusinessType]),
		Name:         company.Name,
		Tin:          company.Tin,
	}
}

var (
	PhoneType = map[int32]string{
		0: "HOME",
		1: "CELL",
		2: "WORK",
		3: "OTHER",
	}
	PhoneTypeReversed = map[string]int32{
		"HOME":  0,
		"CELL":  1,
		"WORK":  2,
		"OTHER": 3,
	}
)

var (
	AddressType = map[int32]string{
		0: "RESIDENTIAL",
		1: "COMMERCIAL",
	}
	AddressTypeReversed = map[string]int32{
		"RESIDENTIAL": 0,
		"COMMERCIAL":  1,
	}
)

var (
	ProfileType = map[int32]string{
		0: "PERSON",
		1: "COMPANY",
	}
	ProfileTypeReversed = map[string]int32{
		"PERSON":  0,
		"COMPANY": 1,
	}
)

var (
	ClaimStatusType = map[int32]string{
		0: "UNKNOWN",
		1: "CLAIM_CREATED",
		2: "CLAIM_ACKNOWLEDGED",
		3: "CLAIM_CANCELLED",
		4: "CLAIM_EXPIRED",
		5: "PAYMENT_IN_PROGRESS",
		6: "PAYMENT_REJECTED",
		7: "PAYMENT_FAILED",
		8: "CLAIMED",
		9: "CLAIM_PAID",
	}
	ClaimStatusTypeReversed = map[string]int32{
		"UNKNOWN":             0,
		"CLAIM_CREATED":       1,
		"CLAIM_ACKNOWLEDGED":  2,
		"CLAIM_CANCELLED":     3,
		"CLAIM_EXPIRED":       4,
		"PAYMENT_IN_PROGRESS": 5,
		"PAYMENT_REJECTED":    6,
		"PAYMENT_FAILED":      7,
		"CLAIMED":             8,
		"CLAIM_PAID":          9,
	}
)

var (
	CompanyType = map[int32]string{
		0: "UNKNOWN",
		1: "CORPORATION",
		2: "PARTNERSHIP",
		3: "GOVERNMENT",
		4: "NONPROFIT",
		5: "PUBLIC_COMPANY",
		6: "PRIVATE_COMPANY",
	}
	CompanyTypeReversed = map[string]int32{
		"UNKNOWN":         0,
		"CORPORATION":     1,
		"PARTNERSHIP":     2,
		"GOVERNMENT":      3,
		"NONPROFIT":       4,
		"PUBLIC_COMPANY":  5,
		"PRIVATE_COMPANY": 6,
	}
)
