# API Contract: Plugin platform

**Base path**: `/rest/plugin/main`  
**Envelope**: Headwind `{ status, message?, data? }`  
**Java reference**: `com.hmdm.plugin.rest.PluginResource`, `com.hmdm.plugin.persistence.PluginDAO`

## GET `/private/available`

**Auth**: Bearer JWT (session compatible)  
**Permission**: authenticated user with customer context

**Response data**: `Plugin[]` — plugins enabled for build, not globally disabled, for current customer catalog.

## GET `/private/active`

**Auth**: Bearer JWT  
**Response data**: `Plugin[]` — available minus `pluginsDisabled` for customer.

## GET `/public/registered`

**Auth**: none  
**Response data**: `Plugin[]` — all plugins registered in deployment (build-filtered).

## POST `/private/disabled`

**Auth**: Bearer JWT  
**Permission**: `plugins_customer_access_management`  
**Body**: JSON array of plugin IDs to disable, e.g. `[2, 5]`

**Response**: `OK` empty data  
**Side effect**: Replace customer disabled set; invalidate plugin status cache.

**Errors**: `error.permission.denied`

## React consumer

`frontend/src/features/plugins/pluginService.ts` — paths above (no `/rest` prefix in client; base URL includes it).
