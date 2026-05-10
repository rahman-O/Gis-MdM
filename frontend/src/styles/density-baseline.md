# Density Baseline Snapshot

Captured before compact redesign rollout.

## Key Screens to Compare
- Dashboard (`/dashboard`)
- Devices (`/devices`)
- Applications (`/applications`)
- Configurations (`/configurations`, `/configurations/:id/edit`)

## Current Global Tokens
- Radius: `--radius: 0.5rem`
- Base body style: default browser + Tailwind utility classes only.
- No explicit density class on root/body.

## Current Shared UI Sizing
- `Button` default: `h-10`, `text-sm`, icon `size-4`
- `Button` sm: `h-9`
- `Input`: `h-10`, `text-sm`
- `SelectTrigger`: `h-10`, `text-sm`, chevron `h-4 w-4`
- `Textarea`: `min-h-[80px]`, `text-sm`
- `TableHead`: `h-12`, `px-4`
- `TableCell`: `p-4`
- `Card` header/content/footer paddings: `p-6`

## Current Page Density Notes
- App layout main padding: `p-6` (`AppLayout`)
- Feature pages commonly use `space-y-6`
- Primary titles often use `text-2xl`
- Alert/action bars mostly use `px-4 py-3`

## Visual Regression Checklist
- Control heights should remain >= 36px interactive target.
- Focus rings must stay visible after downsizing.
- Table readability should remain clear after tighter row heights.
