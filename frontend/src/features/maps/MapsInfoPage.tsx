/**
 * Legacy `hmdmMap` factory exists in the Angular bundle, but no shipped view wires a map tab into the main chrome.
 * This page documents the gap; if product requires live maps, wire Leaflet to the same endpoints as `map.service.js`.
 */
export function MapsInfoPage() {
  return (
    <div className="max-w-2xl space-y-3">
      <h1 className="text-xl font-semibold tracking-tight">Device maps</h1>
      <p className="text-muted-foreground text-sm">
        The legacy SPA includes map utilities, yet the default templates never mount a dedicated map route. Nothing was
        migrated to React until a concrete production use case confirms which legacy controller payload is authoritative.
      </p>
      <p className="text-muted-foreground text-sm">
        When needed, prefer reusing `POST /rest/private/devices/search` (with coordinates if returned) and Leaflet with
        the same attribution policy as your deployment.
      </p>
    </div>
  )
}
