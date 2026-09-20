package dependency

import addresssvc "be/internal/services/address"

func NewAddressService(infra *Infra) *addresssvc.Service {
	return addresssvc.NewService(
		newCountryRepository(infra.DB),
		newDivisionRepository(infra.DB),
	)
}
