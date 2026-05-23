# UI Contract: Configuration editor (016)

**Surface**: React route `/configurations/:id` (`ConfigurationEditorPage`)  
**Consumers**: tenant administrators with `configurations` permission

---

## Tab structure (required)

| Tab ID | Label (en) | Content component |
|--------|------------|-------------------|
| `common` | General | `ConfigurationCommonTab` — name, description, QR key, groups, connectivity basics |
| `mdm` | MDM | `ConfigurationMdmTab` (new or extracted) — main/content app, launcher, event receiver, kiosk core |
| `applications` | Applications | `ConfigurationApplicationsTab` |
| `restrictions` | Restrictions | `ConfigurationRestrictionsTab` (new) — `restrictions` textarea, lock toggles for restriction-related fields |
| `design` | Design | `ConfigurationDesignTab` |
| `appSettings` | App settings | `ConfigurationAppSettingsTab` — per-app defaults + readonly checkbox |
| `files` | Files | `ConfigurationFilesTab` |

**Removed from page shell**: duplicate inline MDM form block; multi-paragraph `CardDescription` help text.

---

## Field lock UX

- Adjacent **lock** control on supported fields (toggle `policyLocks[fieldKey]`).
- Locked fields: visually distinct (`opacity`, lock icon, `aria-locked="true"`).
- Save sends `policyLocks` in PUT body.
- Validation error banner lists `Tab name → Field label`.

---

## Copy guidelines

- Field label: ≤ 4 words where possible.
- Optional hint: one line, `text-xs text-muted-foreground`, max 80 characters.
- No duplicate explanations across tabs.

---

## Acceptance

| # | Check |
|---|--------|
| U1 | All seven tabs render without horizontal scroll on 1280px |
| U2 | Save error references tab name |
| U3 | Lock on `mainAppId` persists after reload |
| U4 | No “phase 1” or developer placeholder strings visible |
