package handlers

import (
	"be/common/httpx"
	addressdto "be/internal/dto/address"
	addresssvc "be/internal/services/address"

	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	addresses *addresssvc.Service
}

func NewAddressHandler(addresses *addresssvc.Service) *AddressHandler {
	return &AddressHandler{addresses: addresses}
}

func (h *AddressHandler) ListCountries(c *gin.Context) {
	result, err := h.addresses.ListCountries(c.Request.Context())
	httpx.OK(c, result, err)
}

func (h *AddressHandler) ListDivisions(c *gin.Context) {
	var q addressdto.DivisionListQuery
	_ = c.ShouldBindQuery(&q)
	result, err := h.addresses.ListDivisions(c.Request.Context(), q)
	httpx.OK(c, result, err)
}
