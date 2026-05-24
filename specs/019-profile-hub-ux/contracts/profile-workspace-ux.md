# Contract: Profile Workspace UX

**Feature**: `019-profile-hub-ux` | **Audience**: Frontend implementers, QA

**Not a REST contract** — defines layout, navigation, and interaction rules for the admin UI.

---

## Shell

| Property | Desktop | Mobile (`< md`) |
|----------|---------|-----------------|
| Container | Dialog ~96vw × 94vh, `bg-black/40` overlay | Full-screen Sheet |
| Close | Header Close + Esc | Header + swipe (if supported) |
| Scroll | Main content area only; header + sidebar fixed | Same; sidebar → drawer |

---

## Layout regions

```text
┌─────────────────────────────────────────────────────────────┐
│ COCKPIT HEADER (fixed)                                       │
│ Name · Health · Lifecycle · Published vN · Assigned · Actions│
├──────────┬──────────────────────────────────────────────────┤
│ SIDEBAR  │ MAIN CONTENT                                      │
│ Overview │ (section-specific)                                │
│ Assign…  │                                                   │
│ Rollout  │                                                   │
│ Versions │                                                   │
│ Editor   │                                                   │
│ Activity │                                                   │
└──────────┴──────────────────────────────────────────────────┘
```

---

## Cockpit header actions

| Action | Behavior |
|--------|----------|
| **Edit** | Sets `section=editor`; if no draft, offer fork-from-published |
| **Publish** | Opens `secondaryPanel=publish-impact` (side sheet, NOT nested dialog) |
| **Close** | If `editorDirty`, confirm; else close workspace |

Publish button states:

- Disabled: `!canPublish` or validation errors
- Enabled: draft has changes and passes client validation

---

## Sections

| Section | Mode | Content |
|---------|------|---------|
| **Overview** | Read | 6 cards grid (Status, Assignment, Rollout, Apps, Kiosk, Last publish) — **no inputs** |
| **Assignments** | Read/Write | 018 assignment panel + `TreePreview` |
| **Rollout** | Read | 018 rollout table + filters |
| **Versions** | Read/Write | 018 version list + fork |
| **Editor** | Write | Warning bar + sticky save + configuration tabs |
| **Activity** | Read | Timeline list |

Default section on open: **Overview**.

---

## Read vs Edit visual tokens

| Token | Read | Edit (Editor section) |
|-------|------|------------------------|
| Page background | `bg-muted/30` | `bg-background` + left border accent |
| Inputs | None in Overview | Full forms in Editor |
| Banner | None | Amber warning: production policy change |

---

## Secondary panels (allowed)

| Panel | Trigger | Component |
|-------|---------|-----------|
| Publish impact | Header Publish | Right `Sheet` ~400px |
| Assignment confirm (≥50 devices) | Assign save | Inline in Assignments or right Sheet |

**Forbidden**: `Dialog` opened while workspace Dialog/Sheet is open.

---

## Profiles list (Control Radar)

Columns/badges minimum:

- Name (click opens workspace)
- Health chip
- Badges row: `No Assignment`, `Disabled`, `Draft Changes`, `Rollout Issues`, `Stale` (max 3 visible + "+N")

Optional: row quick actions menu (P2).

---

## Deep linking

| URL | Behavior |
|-----|----------|
| `/profiles?open={id}` | Open workspace Overview |
| `/profiles?open={id}&section=editor` | Open Editor |
| `/profiles/{id}/edit` (legacy) | Redirect to `?open={id}&section=editor` |

---

## Create profile flow

1. List → New profile form (existing or inline)
2. On success → open workspace `?open={id}&wizard=assign`
3. If no published version → highlight publish CTA in header before assign
4. On assign success → `section=overview` + success toast

---

## Keyboard (P2)

| Key | Action |
|-----|--------|
| `Esc` | Close workspace (with dirty check) |
| `E` | Go to Editor |
| `Ctrl/Cmd+S` | Save draft (Editor only) |

---

## 5-second rule (QA)

From Overview only, user must answer within 5s:

1. Is a version published?
2. Is profile assigned to tree?
3. Is profile disabled?
4. Any rollout failures?

Test with 10 admins; pass if ≥9/10 correct.
