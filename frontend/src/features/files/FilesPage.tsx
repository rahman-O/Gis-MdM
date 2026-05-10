import { useCallback, useEffect, useState } from 'react'
import { AlertCircle, Loader2, Trash2 } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Skeleton } from '@/shared/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/shared/ui/table'
import { useToast } from '@/shared/hooks/use-toast'
import * as filesService from '@/features/files/filesService'

export function FilesPage() {
  const { toast } = useToast()
  const [rows, setRows] = useState<filesService.FileRecord[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [deletingId, setDeletingId] = useState<number | null>(null)

  const load = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      setRows(await filesService.searchAllFiles())
    } catch (reason: unknown) {
      setError(reason instanceof Error ? reason.message : 'Failed to load files.')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void load()
  }, [load])

  async function removeRow(row: filesService.FileRecord) {
    if (row.id == null) return
    setDeletingId(row.id)
    try {
      await filesService.removeFile({ id: row.id, filePath: row.filePath ?? undefined })
      toast({ title: 'File deleted' })
      await load()
    } catch (reason: unknown) {
      toast({
        title: 'Delete failed',
        variant: 'destructive',
        description: reason instanceof Error ? reason.message : undefined,
      })
    } finally {
      setDeletingId(null)
    }
  }

  if (loading) {
    return (
      <div className="space-y-3">
        <h1 className="text-xl font-semibold tracking-tight">Files</h1>
        <Skeleton className="h-32 w-full" />
      </div>
    )
  }

  if (error) {
    return (
      <div className="flex items-center gap-2 rounded-md border border-destructive/50 bg-destructive/10 p-3 text-sm">
        <AlertCircle className="text-destructive h-4 w-4 shrink-0" />
        <span className="flex-1">{error}</span>
        <Button size="sm" variant="outline" onClick={() => void load()}>
          Retry
        </Button>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      <div>
        <h1 className="text-xl font-semibold tracking-tight">Files</h1>
        <p className="text-muted-foreground text-sm">Library referenced by configurations and icons.</p>
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Path / URL</TableHead>
            <TableHead>Note</TableHead>
            <TableHead className="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {rows.map((row) => (
            <TableRow key={String(row.id ?? row.url)}>
              <TableCell className="max-w-xl truncate">{row.url || row.filePath || '—'}</TableCell>
              <TableCell className="text-muted-foreground text-sm">{row.description || '—'}</TableCell>
              <TableCell className="text-right">
                <Button
                  type="button"
                  variant="destructive"
                  size="sm"
                  disabled={row.id == null || deletingId === row.id}
                  onClick={() => void removeRow(row)}
                >
                  {deletingId === row.id ? <Loader2 className="mr-2 h-3 w-3 animate-spin" /> : <Trash2 className="h-3 w-3" />}
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      {rows.length === 0 ? <p className="text-muted-foreground text-sm">No files found.</p> : null}
    </div>
  )
}
