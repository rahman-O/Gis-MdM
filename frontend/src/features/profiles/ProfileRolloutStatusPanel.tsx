import { useCallback, useEffect, useState } from 'react'
import { Button } from '@/shared/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shared/ui/table'
import {
  listRolloutDevices,
  recomputeRollout,
  type DeviceRolloutRow,
} from '@/features/profiles/profileRolloutService'

interface Props {
  profileId: number
}

const statusClass: Record<string, string> = {
  pending: 'text-amber-600',
  installed: 'text-emerald-600',
  partial: 'text-orange-600',
  failed: 'text-red-600',
}

export function ProfileRolloutStatusPanel({ profileId }: Props) {
  const [rows, setRows] = useState<DeviceRolloutRow[]>([])
  const [loading, setLoading] = useState(true)

  const load = useCallback(async () => {
    try {
      const page = await listRolloutDevices(profileId, { pageSize: 50 })
      setRows(page.items ?? [])
    } catch {
      setRows([])
    } finally {
      setLoading(false)
    }
  }, [profileId])

  useEffect(() => {
    void load()
    const id = window.setInterval(() => void load(), 60_000)
    return () => window.clearInterval(id)
  }, [load])

  return (
    <div className="space-y-3">
      <div className="flex justify-end gap-2">
        <Button type="button" variant="outline" size="sm" onClick={() => void load()}>
          Refresh
        </Button>
        <Button type="button" variant="outline" size="sm" onClick={() => void recomputeRollout(profileId).then(load)}>
          Recompute all
        </Button>
      </div>
      {loading ? (
        <p className="text-sm text-muted-foreground">Loading rollout status…</p>
      ) : (
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Device</TableHead>
              <TableHead>Folder</TableHead>
              <TableHead>Target</TableHead>
              <TableHead>Applied</TableHead>
              <TableHead>Status</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {rows.map((r) => (
              <TableRow key={r.deviceId}>
                <TableCell>{r.deviceName}</TableCell>
                <TableCell>{r.treeNodeName || '—'}</TableCell>
                <TableCell>
                  {r.targetVersionNumber != null ? `v${r.targetVersionNumber}` : '—'}
                </TableCell>
                <TableCell>
                  {r.appliedVersionNumber != null ? `v${r.appliedVersionNumber}` : '—'}
                </TableCell>
                <TableCell className={statusClass[r.status] ?? ''}>
                  {r.status}
                  {r.reason ? ` — ${r.reason}` : ''}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      )}
    </div>
  )
}
