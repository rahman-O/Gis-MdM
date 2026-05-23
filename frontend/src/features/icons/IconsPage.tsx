import { useCallback, useEffect, useRef, useState } from 'react'
import { Loader2 } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { Skeleton } from '@/shared/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/shared/ui/table'
import { useToast } from '@/shared/hooks/use-toast'
import * as iconsService from '@/features/icons/iconsService'

export function IconsPage() {
  const { toast } = useToast()
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [loading, setLoading] = useState(true)
  const [rows, setRows] = useState<iconsService.IconRow[]>([])
  const [search, setSearch] = useState('')
  const [name, setName] = useState('')
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
    if (!name.trim()) {
      toast({ title: 'Icon name is required.', variant: 'destructive' })
      return
    }
    const file = fileInputRef.current?.files?.[0]
    if (!file) {
      toast({ title: 'Choose an image file to upload.', variant: 'destructive' })
      return
    }
    setBusy(true)
    try {
      const { fileId } = await iconsService.uploadIconFile(file)
      await iconsService.saveIcon({ id: undefined, name: name.trim(), fileId })
      toast({ title: 'Icon saved' })
      setName('')
      if (fileInputRef.current) {
        fileInputRef.current.value = ''
      }
      await load()
    } catch (reason: unknown) {
      toast({
        title: 'Save failed',
        variant: 'destructive',
        description: reason instanceof Error ? reason.message : undefined,
      })
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="space-y-4">
      <h1 className="text-xl font-semibold tracking-tight">Icons</h1>
      <p className="text-muted-foreground text-sm">Upload a square PNG and create an icon record (Go `POST /private/icon-files`).</p>

      <div className="flex flex-wrap items-end gap-2 max-w-xl">
        <div className="min-w-[200px] flex-1 space-y-1">
          <label className="text-xs font-medium">Name</label>
          <Input value={name} onChange={(e) => setName(e.target.value)} placeholder="Home" />
        </div>
        <div className="min-w-[200px] flex-1 space-y-1">
          <label className="text-xs font-medium">Image file</label>
          <Input ref={fileInputRef} type="file" accept="image/png,image/jpeg,image/*" />
        </div>
        <Button type="button" onClick={() => void addIcon()} disabled={busy}>
          {busy ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
          Upload &amp; save
        </Button>
      </div>

      <div className="flex gap-2 max-w-md">
        <Input placeholder="Search" value={search} onChange={(e) => setSearch(e.target.value)} />
        <Button type="button" variant="outline" onClick={() => void load()}>
          Search
        </Button>
      </div>

      {loading ? (
        <Skeleton className="h-40 w-full" />
      ) : (
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Name</TableHead>
              <TableHead>File</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {rows.map((row) => (
              <TableRow key={row.id ?? row.name}>
                <TableCell>{row.id ?? '—'}</TableCell>
                <TableCell>{row.name ?? '—'}</TableCell>
                <TableCell>{row.fileName ?? row.fileId ?? '—'}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      )}
    </div>
  )
}
