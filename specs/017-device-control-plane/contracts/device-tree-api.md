# Contract: Device Tree API

**Feature**: `017-device-control-plane` | **Audience**: Admin React client

**Base**: `/rest/private` | **Envelope**: `{ "status": "OK"|"ERROR", "message"?, "data"? }`

## Endpoints

### `GET /device-tree`

Returns full tree for current `customerId`.

**Response `data`**:

```json
{
  "nodes": [
    {
      "id": 1,
      "parentId": null,
      "name": "All devices",
      "sortOrder": 0,
      "path": "/1/",
      "depth": 0,
      "deviceCount": 42
    }
  ],
  "rootId": 1
}
```

### `POST /device-tree/nodes`

Create folder.

**Body**: `{ "parentId": 1, "name": "Warehouse A", "sortOrder": 0 }`

**Errors**: `error.device_tree.duplicate_name`, `error.device_tree.invalid_parent`

### `PUT /device-tree/nodes/:id`

Rename or reorder.

**Body**: `{ "name"?, "sortOrder"?, "parentId"? }`

**Errors**: `error.device_tree.cycle` when move would create loop

### `POST /device-tree/nodes/:id/delete`

Mandatory relocation when subtree has devices.

**Body**: `{ "targetNodeId": 2 }`

**Behavior**: Move all devices in subtree to target; delete node (and empty children per product rule).

### `POST /devices/:id/move-tree`

Move single device.

**Body**: `{ "treeNodeId": 5 }`

### `GET /devices?treeNodeId=:id&includeDescendants=true`

List devices filtered by tree selection (extends existing devices list).

## Permissions

Reuse `devices` read/write; new permission `device_tree` optional alias to `devices` in v1.

## Parity

New surface — document in `serverBackendGo/docs/parity/device-control-plane.md`.
