package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"be/common/httpx"
	"be/common/response"
	mapsdto "be/internal/dto/maps"
	locationsvc "be/internal/services/location"
	mapssvc "be/internal/services/maps"
	ingestsvc "be/internal/services/maps/ingest"
	sourcesvc "be/internal/services/maps/source"
	placesvc "be/internal/services/place"
)

type MapsHandler struct {
	places    *placesvc.Service
	locations *locationsvc.Service
	sources   *sourcesvc.Service
	ingest    *ingestsvc.Service
	search    *mapssvc.Service
}

func NewMapsHandler(
	places *placesvc.Service,
	locations *locationsvc.Service,
	sources *sourcesvc.Service,
	ingest *ingestsvc.Service,
	search *mapssvc.Service,
) *MapsHandler {
	return &MapsHandler{
		places:    places,
		locations: locations,
		sources:   sources,
		ingest:    ingest,
		search:    search,
	}
}

func (h *MapsHandler) Search(c *gin.Context) {
	var q mapsdto.SearchQuery
	_ = c.ShouldBindQuery(&q)
	result, err := h.search.Search(c.Request.Context(), q)
	httpx.OK(c, result, err)
}

func (h *MapsHandler) ListPlaces(c *gin.Context) {
	var q mapsdto.PlaceListQuery
	_ = c.ShouldBindQuery(&q)
	result, err := h.places.List(c.Request.Context(), q)
	httpx.OK(c, result, err)
}

func (h *MapsHandler) GetPlace(c *gin.Context) {
	result, err := h.places.Get(c.Request.Context(), c.Param("id"))
	httpx.OK(c, result, err)
}

func (h *MapsHandler) PatchPlace(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	var body mapsdto.PatchPlaceRequest
	if !httpx.BindJSON(c, &body) {
		return
	}
	result, err := h.places.Patch(c.Request.Context(), c.Param("id"), userID, body)
	httpx.OK(c, result, err)
}

func (h *MapsHandler) CreatePlace(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	var body mapsdto.CreatePlaceRequest
	if !httpx.BindJSON(c, &body) {
		return
	}
	result, err := h.places.Create(c.Request.Context(), userID, body)
	httpx.Created(c, result, err)
}

func (h *MapsHandler) DeletePlace(c *gin.Context) {
	httpx.NoContent(c, h.places.Delete(c.Request.Context(), c.Param("id")))
}

func (h *MapsHandler) ListLocations(c *gin.Context) {
	var q mapsdto.LocationListQuery
	_ = c.ShouldBindQuery(&q)
	result, err := h.locations.List(c.Request.Context(), q)
	httpx.OK(c, result, err)
}

func (h *MapsHandler) GetLocation(c *gin.Context) {
	result, err := h.locations.Get(c.Request.Context(), c.Param("id"))
	httpx.OK(c, result, err)
}

func (h *MapsHandler) CreateLocation(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	var body mapsdto.CreateLocationRequest
	if !httpx.BindJSON(c, &body) {
		return
	}
	result, err := h.locations.Create(c.Request.Context(), userID, body)
	httpx.Created(c, result, err)
}

func (h *MapsHandler) PatchLocation(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	var body mapsdto.PatchLocationRequest
	if !httpx.BindJSON(c, &body) {
		return
	}
	result, err := h.locations.Patch(c.Request.Context(), c.Param("id"), userID, body)
	httpx.OK(c, result, err)
}

func (h *MapsHandler) DeleteLocation(c *gin.Context) {
	httpx.NoContent(c, h.locations.Delete(c.Request.Context(), c.Param("id")))
}

func (h *MapsHandler) ListCategories(c *gin.Context) {
	result, err := h.places.ListCategories(c.Request.Context())
	httpx.OK(c, result, err)
}

func (h *MapsHandler) ListSources(c *gin.Context) {
	var q mapsdto.SourceListQuery
	_ = c.ShouldBindQuery(&q)
	result, err := h.sources.List(c.Request.Context(), q)
	httpx.OK(c, result, err)
}

func (h *MapsHandler) GetSource(c *gin.Context) {
	result, err := h.sources.Get(c.Request.Context(), c.Param("id"))
	httpx.OK(c, result, err)
}

func (h *MapsHandler) ProbeSource(c *gin.Context) {
	var body mapsdto.ProbeSourceRequest
	if !httpx.BindJSON(c, &body) {
		return
	}
	result, err := h.sources.Probe(c.Request.Context(), body)
	httpx.OK(c, result, err)
}

func (h *MapsHandler) CreateSource(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	var body mapsdto.CreateSourceRequest
	if !httpx.BindJSON(c, &body) {
		return
	}
	result, err := h.sources.Create(c.Request.Context(), userID, body)
	httpx.Created(c, result, err)
}

func (h *MapsHandler) PatchSource(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	var body mapsdto.PatchSourceRequest
	if !httpx.BindJSON(c, &body) {
		return
	}
	result, err := h.sources.Patch(c.Request.Context(), c.Param("id"), userID, body)
	httpx.OK(c, result, err)
}

func (h *MapsHandler) DeleteSource(c *gin.Context) {
	httpx.NoContent(c, h.sources.Delete(c.Request.Context(), c.Param("id")))
}

func (h *MapsHandler) StartIngest(c *gin.Context) {
	userID, ok := httpx.RequireUserID(c)
	if !ok {
		return
	}
	run, sync, err := h.ingest.Start(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, ingestsvc.ErrIngestInProgress) {
			response.Error(c, http.StatusConflict, "ingest already in progress")
			return
		}
		response.HandleError(c, err)
		return
	}
	if sync {
		httpx.OK(c, run, nil)
		return
	}
	response.JSON(c, http.StatusAccepted, run)
}

func (h *MapsHandler) LatestIngest(c *gin.Context) {
	result, err := h.ingest.Latest(c.Request.Context())
	httpx.OK(c, result, err)
}

func (h *MapsHandler) GetIngest(c *gin.Context) {
	result, err := h.ingest.Get(c.Request.Context(), c.Param("id"))
	httpx.OK(c, result, err)
}
