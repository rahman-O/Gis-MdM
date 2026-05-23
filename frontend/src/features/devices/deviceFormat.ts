export function formatLastSeen(ms: number | null | undefined): string {
  if (ms == null || ms <= 0) {
    return '—'
  }
  try {
    return new Intl.DateTimeFormat(undefined, {
      dateStyle: 'medium',
      timeStyle: 'short',
    }).format(new Date(ms))
  } catch {
    return String(ms)
  }
}
