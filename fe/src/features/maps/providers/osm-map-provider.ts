import {
  DEFAULT_MAP_VIEWPORT,
  type LatLng,
  type MapBasemap,
  type MapPin,
  type MapProvider,
  type MapScreenPoint,
  type MapViewport,
} from './map-provider'

type LeafletClickEvent = {
  latlng: { lat: number; lng: number }
  containerPoint: { x: number; y: number }
}

type LeafletModule = {
  map: (...args: never[]) => LeafletMap
  marker: (...args: never[]) => LeafletMarker
  divIcon: (...args: never[]) => unknown
  tileLayer: (
    url: string,
    options?: { attribution?: string; maxZoom?: number },
  ) => LeafletTileLayer
  latLngBounds: (
    positions: Array<[number, number]>,
  ) => unknown
  DomEvent?: { stopPropagation: (event: unknown) => void }
  default?: LeafletModule
}

type LeafletMap = {
  setView: (center: [number, number], zoom?: number) => LeafletMap
  setZoom: (zoom: number) => LeafletMap
  fitBounds: (bounds: unknown, options?: unknown) => LeafletMap
  latLngToContainerPoint: (latlng: { lat: number; lng: number }) => {
    x: number
    y: number
  }
  on: (event: string, handler: (event: LeafletClickEvent) => void) => LeafletMap
  off: (event: string, handler: (event: LeafletClickEvent) => void) => LeafletMap
  remove: () => void
  zoomControl?: { setPosition: (position: string) => void }
}

type LeafletTileLayer = {
  addTo: (map: LeafletMap) => LeafletTileLayer
  remove: () => void
}

type LeafletMarker = {
  setLatLng: (latlng: [number, number]) => LeafletMarker
  bindTooltip: (text: string) => LeafletMarker
  on: (event: string, handler: (event?: unknown) => void) => LeafletMarker
  addTo: (map: LeafletMap) => LeafletMarker
  remove: () => void
}

const BASEMAPS: Record<
  MapBasemap,
  { url: string; attribution: string; maxZoom: number }
> = {
  street: {
    url: 'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
    attribution: '&copy; OpenStreetMap contributors',
    maxZoom: 19,
  },
  satellite: {
    url: 'https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}',
    attribution: 'Tiles &copy; Esri',
    maxZoom: 19,
  },
}

const PIN_STYLE_ID = 'maps-osm-pin-styles'

function ensurePinStyles() {
  if (typeof document === 'undefined') return
  if (document.getElementById(PIN_STYLE_ID)) return
  const style = document.createElement('style')
  style.id = PIN_STYLE_ID
  style.textContent = `
    .leaflet-container,
    .leaflet-control,
    .leaflet-tooltip,
    .leaflet-popup-content {
      font-family: var(--font-sans), ui-sans-serif, system-ui, sans-serif;
    }
    .leaflet-container.maps-map-pick,
    .leaflet-container.maps-map-pick .leaflet-interactive {
      cursor: crosshair;
    }
    .maps-osm-pin { background: transparent; border: none; }
    .maps-osm-pin-badge {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 28px;
      height: 28px;
      border-radius: 9999px;
      border: 2px solid var(--background);
      box-shadow: 0 1px 4px rgb(0 0 0 / 0.35);
      color: #fff;
    }
    .maps-osm-pin-roomRental { background: #4f46e5; }
    .maps-osm-pin-restaurant { background: #dc2626; }
    .maps-osm-pin-hotel { background: #0284c7; }
    .maps-osm-pin-eatery { background: #d97706; }
    .maps-osm-pin-uncategorized { background: #6b7280; }
    .maps-osm-pin-user-dot {
      display: block;
      width: 16px;
      height: 16px;
      border-radius: 9999px;
      background: var(--foreground);
      border: 3px solid var(--background);
      box-shadow: 0 0 0 6px color-mix(in oklch, var(--foreground) 25%, transparent);
    }
    .maps-osm-pin-draft-dot {
      display: block;
      width: 16px;
      height: 16px;
      border-radius: 9999px;
      background: var(--primary);
      border: 2px dashed var(--background);
      box-shadow: 0 0 0 5px color-mix(in oklch, var(--primary) 35%, transparent);
    }
    .maps-osm-pin-location {
      background: #0f766e;
      border-radius: 6px;
    }
  `
  document.head.appendChild(style)
}

function categoryIconSvg(paths: string) {
  return `<svg viewBox="0 0 24 24" width="14" height="14" aria-hidden="true" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">${paths}</svg>`
}

const CATEGORY_PIN_SVG: Record<string, string> = {
  roomRental: categoryIconSvg(
    '<path d="M15 21v-8a1 1 0 0 0-1-1h-4a1 1 0 0 0-1 1v8"/><path d="M3 10a2 2 0 0 1 .709-1.528l7-5.999a2 2 0 0 1 2.582 0l7 5.999A2 2 0 0 1 21 10v9a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>',
  ),
  restaurant: categoryIconSvg(
    '<path d="M3 2v7c0 1.1.9 2 2 2h4a2 2 0 0 0 2-2V2"/><path d="M7 2v20"/><path d="M21 15V2a5 5 0 0 0-5 5v6c0 1.1.9 2 2 2h3Zm0 0v7"/>',
  ),
  hotel: categoryIconSvg(
    '<path d="M2 4v16"/><path d="M2 8h18a2 2 0 0 1 2 2v10"/><path d="M2 17h20"/><path d="M6 8v9"/>',
  ),
  eatery: categoryIconSvg(
    '<path d="M10 2v2"/><path d="M14 2v2"/><path d="M16 8a1 1 0 0 1 1 1v8a4 4 0 0 1-4 4H7a4 4 0 0 1-4-4V9a1 1 0 0 1 1-1h14a4 4 0 1 1 0 8h-1"/><path d="M6 2v2"/>',
  ),
  uncategorized: categoryIconSvg(
    '<path d="M20 10c0 4.993-5.539 10.193-7.399 11.799a1 1 0 0 1-1.202 0C9.539 20.193 4 14.993 4 10a8 8 0 0 1 16 0"/><circle cx="12" cy="10" r="3"/>',
  ),
}

function pinMarkup(pin: MapPin) {
  const variant =
    pin.variant ?? (pin.interactive === false ? 'user' : 'default')
  if (variant === 'user') {
    return {
      html: '<span class="maps-osm-pin-user-dot"></span>',
      iconSize: [16, 16] as [number, number],
      iconAnchor: [8, 8] as [number, number],
    }
  }
  if (variant === 'draft') {
    return {
      html: '<span class="maps-osm-pin-draft-dot"></span>',
      iconSize: [16, 16] as [number, number],
      iconAnchor: [8, 8] as [number, number],
    }
  }
  if (variant === 'location') {
    return {
      html: `<span class="maps-osm-pin-badge maps-osm-pin-location">${categoryIconSvg(
        '<path d="M6 22V4a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v18Z"/><path d="M6 12h4"/><path d="M6 16h4"/><path d="M14 12h4"/><path d="M14 16h4"/><path d="M10 22v-4h4v4"/>',
      )}</span>`,
      iconSize: [28, 28] as [number, number],
      iconAnchor: [14, 14] as [number, number],
    }
  }
  const category = CATEGORY_PIN_SVG[pin.category ?? '']
    ? pin.category
    : 'uncategorized'
  return {
    html: `<span class="maps-osm-pin-badge maps-osm-pin-${category}">${CATEGORY_PIN_SVG[category ?? 'uncategorized']}</span>`,
    iconSize: [28, 28] as [number, number],
    iconAnchor: [14, 14] as [number, number],
  }
}

/**
 * Live OpenStreetMap adapter. Leaflet and OSM tile URLs stay in this file.
 */
export function createOsmMapProvider(): MapProvider {
  let map: LeafletMap | null = null
  let Lref: LeafletModule | null = null
  const markers = new Map<string, LeafletMarker>()
  let clickHandler: ((pinId: string) => void) | null = null
  let hoverHandler: ((pinId: string | null, screen?: MapScreenPoint) => void) | null =
    null
  let mapClickHandler: ((point: LatLng, screen: MapScreenPoint) => void) | null =
    null
  let pendingPins: MapPin[] = []
  let mountGeneration = 0
  let mapClickListener: ((event: LeafletClickEvent) => void) | null = null
  let tiles: LeafletTileLayer | null = null
  let currentBasemap: MapBasemap = 'street'

  function addPinNow(pin: MapPin) {
    if (!map || !Lref) return
    const existing = markers.get(pin.id)
    if (existing) {
      existing.setLatLng([pin.position.lat, pin.position.lng])
      return
    }
    const icon = pinMarkup(pin)
    const marker = Lref.marker([pin.position.lat, pin.position.lng], {
      icon: Lref.divIcon({
        className: 'maps-osm-pin',
        html: icon.html,
        iconSize: icon.iconSize,
        iconAnchor: icon.iconAnchor,
      }),
      title: pin.title,
      keyboard: pin.interactive !== false,
    })
    if (pin.interactive !== false) {
      marker.on('click', (event) => {
        Lref?.DomEvent?.stopPropagation(event)
        hoverHandler?.(null)
        clickHandler?.(pin.id)
      })
      marker.on('mouseover', (raw) => {
        const event = raw as LeafletClickEvent | undefined
        const point =
          event?.containerPoint ??
          (event?.latlng
            ? map?.latLngToContainerPoint(event.latlng)
            : undefined)
        if (!point) return
        hoverHandler?.(pin.id, { x: point.x, y: point.y })
      })
      marker.on('mouseout', () => {
        hoverHandler?.(null)
      })
    } else {
      marker.on('click', (event) => {
        Lref?.DomEvent?.stopPropagation(event)
      })
    }
    marker.addTo(map)
    markers.set(pin.id, marker)
  }

  function applyBasemap(basemap: MapBasemap) {
    currentBasemap = basemap
    if (!map || !Lref) return
    tiles?.remove()
    const spec = BASEMAPS[basemap]
    tiles = Lref.tileLayer(spec.url, {
      attribution: spec.attribution,
      maxZoom: spec.maxZoom,
    }).addTo(map)
  }

  async function mountAsync(container: HTMLElement, viewport?: MapViewport) {
    const generation = ++mountGeneration
    const leaflet = (await import('leaflet')) as LeafletModule & {
      default: LeafletModule
    }
    await import('leaflet/dist/leaflet.css')
    if (generation !== mountGeneration) return

    const L = leaflet.default ?? leaflet
    Lref = L
    ensurePinStyles()

    const view = viewport ?? DEFAULT_MAP_VIEWPORT
    map = L.map(container, { zoomControl: true }).setView(
      [view.center.lat, view.center.lng],
      view.zoom,
    )
    map.zoomControl?.setPosition('bottomright')
    mapClickListener = (event) => {
      mapClickHandler?.(
        { lat: event.latlng.lat, lng: event.latlng.lng },
        { x: event.containerPoint.x, y: event.containerPoint.y },
      )
    }
    map.on('click', mapClickListener)
    applyBasemap(currentBasemap)

    for (const pin of pendingPins) {
      addPinNow(pin)
    }
    pendingPins = []
  }

  return {
    mount(container, viewport) {
      void mountAsync(container, viewport)
    },
    unmount() {
      mountGeneration += 1
      markers.forEach((marker) => {
        marker.remove()
      })
      markers.clear()
      pendingPins = []
      if (map) {
        if (mapClickListener) {
          map.off('click', mapClickListener)
        }
        map.remove()
        map = null
      }
      mapClickListener = null
      tiles = null
      Lref = null
    },
    setCenter(center) {
      map?.setView([center.lat, center.lng])
    },
    setZoom(zoom) {
      map?.setZoom(zoom)
    },
    setBasemap(basemap) {
      applyBasemap(basemap)
    },
    addPin(pin) {
      if (!map || !Lref) {
        const index = pendingPins.findIndex((item) => item.id === pin.id)
        if (index >= 0) {
          pendingPins[index] = pin
        } else {
          pendingPins.push(pin)
        }
        return
      }
      addPinNow(pin)
    },
    updatePin(pin) {
      const marker = markers.get(pin.id)
      if (!marker) {
        this.addPin(pin)
        return
      }
      marker.setLatLng([pin.position.lat, pin.position.lng])
    },
    removePin(id) {
      pendingPins = pendingPins.filter((pin) => pin.id !== id)
      const marker = markers.get(id)
      if (marker) {
        marker.remove()
        markers.delete(id)
      }
    },
    clearPins() {
      pendingPins = []
      markers.forEach((marker) => {
        marker.remove()
      })
      markers.clear()
    },
    onPinClick(handler) {
      clickHandler = handler
    },
    onPinHover(handler) {
      hoverHandler = handler
    },
    onMapClick(handler) {
      mapClickHandler = handler
    },
    fitBounds(positions: LatLng[]) {
      if (!map || !Lref || positions.length === 0) return
      if (positions.length === 1) {
        map.setView([positions[0].lat, positions[0].lng], 16)
        return
      }
      const bounds = Lref.latLngBounds(
        positions.map((position) => [position.lat, position.lng] as [number, number]),
      )
      map.fitBounds(bounds, { padding: [40, 40], maxZoom: 16 })
    },
  }
}
