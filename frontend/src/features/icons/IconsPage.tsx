import { useCallback, useEffect, useState } from 'react'
import { Loader2 } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { Skeleton } from '@/shared/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/shared/ui/table'
import { useToast } from '@/shared/hooks/use-toast'
import * as iconsService from '@/features/icons/iconsService'

export function IconsPage() {
  const { toast } = useToast()
  const [loading, setLoading] = useState(true)
  const [rows, setRows] = useState<iconsService.IconRow[]>([])
  const [search, setSearch] = useState('')
  const [name, setName] = useState('')
  const [fileId, setFileId] = useState('')
  const [busy, setBusy] = useState(false)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      setRows(await iconsService.listIcons(search.trim() || undefined))
    } catch {
      toast({ title: 'Failed to load icons', variant: 'destructive' })
    } finally {
      setLoading(false)
    }
  }, [search, toast])

  useEffect(() => {
    void load()
  }, [load])

  async function addIcon() {
    setBusy(true)
    try {
      const fid = Number(fileId.trim())
      if (!name.trim() || !fid) {
        toast({ title: 'Name and numeric file ID are required.', variant: 'destructive' })
        return
      }
      await iconsService.saveIcon({ id: undefined, name: name.trim(), fileId: fid })
      toast({ title: 'Icon saved' })
      setName('')
      setFileId('')
      await load()
    } catch {
      toast({ title: 'Save failed', variant: 'destructive' })
    } finally {
      setBusy(false)
    }
  }

  async function remove(id: number) {
    setBusy(true)
    try {
      await iconsService.deleteIcon(id)
      toast({ title: 'Icon deleted' })
      await load()
    } catch {
      toast({ title: 'Delete failed', variant: 'destructive' })
    } finally {
      setBusy(false)
    }
  }

  if (loading) {
    return (
      <div className="space-y-3">
        <Skeleton className="h-24 w-full" />
      </div>
    )
  }

  return (
    <div className="space-y-5">
      <div>
        <h1 className="text-xl font-semibold tracking-tight">Icons</h1>
        <p className="text-muted-foreground text-sm">Reference uploaded file IDs.</p>
      </div>
      <div className="flex max-w-xl flex-wrap items-end gap-2">
        <div className="flex-1 min-w-[120px]">
          <label className="mb-1 block text-xs text-muted-foreground">Filter value</label>
          <Input value={search} onChange={(e) => setSearch(e.target.value)} />
        </div>
        <Button type="button" variant="secondary" onClick={() => void load()}>
          Search
        </Button>
      </div>
      <div className="flex max-w-xl flex-wrap items-end gap-2 rounded border p-3">
        <div className="flex-1">
          <label className="mb-1 block text-xs text-muted-foreground">New name</label>
          <Input value={name} onChange={(e) => setName(e.target.value)} />
        </div>
        <div className="w-28">
          <label className="mb-1 block text-xs text-muted-foreground">File ID</label>
          <Input value={fileId} onChange={(e) => setFileId(e.target.value)} />
        </div>
        <Button type="button" onClick={() => void addIcon()} disabled={busy}>
          {busy ? <Loader2 className="h-4 w-4 animate-spin" /> : 'Add'}
        </Button>
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Name</TableHead>
            <TableHead>File</TableHead>
            <TableHead className="text-right"></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {rows.map((r) => (
            <TableRow key={String(r.id)}>
              <TableCell>{r.name}</TableCell>
              <TableCell className="text-muted-foreground text-sm">
                #{r.fileId} {r.fileName ? `· ${r.fileName}` : ''}
              </TableCell>
              <TableCell className="text-right">
                <Button type="button" variant="destructive" size="sm" disabled={!r.id || busy} onClick={() => remove(r.id!)}>
                  Delete
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
