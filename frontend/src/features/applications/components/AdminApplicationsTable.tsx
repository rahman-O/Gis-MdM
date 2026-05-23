import { Pencil, Share2, Trash2 } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/shared/ui/table'
import type { Application } from '@/features/applications/model/types'

interface Props {
  applications: Application[]
  canEdit: boolean
  onEdit: (app: Application) => void
  onDelete: (app: Application) => void
  onTurnCommon: (app: Application) => void
}

export function AdminApplicationsTable({
  applications,
  canEdit,
  onEdit,
  onDelete,
  onTurnCommon,
}: Props) {
  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Organization</TableHead>
            <TableHead>Package</TableHead>
            <TableHead>Name</TableHead>
            <TableHead>Version</TableHead>
            <TableHead>URL</TableHead>
            <TableHead>Icon</TableHead>
            <TableHead className="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {applications.length === 0 ? (
            <TableRow>
              <TableCell colSpan={7} className="h-24 text-center text-muted-foreground">
                No applications to display.
              </TableCell>
            </TableRow>
          ) : (
            applications.map((app) => {
              const isCommon = Boolean(app.commonApplication)
              return (
                <TableRow key={app.id ?? `${app.customerId}-${app.pkg}-${app.version}`}>
                  <TableCell>{isCommon ? '—' : (app.customerName ?? '—')}</TableCell>
                  <TableCell>{app.pkg ?? '—'}</TableCell>
                  <TableCell className="font-medium">{app.name ?? '—'}</TableCell>
                  <TableCell>{app.version ?? app.latestVersionText ?? '—'}</TableCell>
                  <TableCell className="max-w-[240px] truncate">{app.url ?? '—'}</TableCell>
                  <TableCell>{app.showIcon ? '✓' : ''}</TableCell>
                  <TableCell className="text-right">
                    {!canEdit ? (
                      <span className="text-muted-foreground">—</span>
                    ) : isCommon ? (
                      <div className="flex justify-end gap-1">
                        <Button type="button" size="icon" variant="ghost" title="Edit shared application" onClick={() => onEdit(app)}>
                          <Pencil className="h-4 w-4" />
                        </Button>
                        <Button
                          type="button"
                          size="icon"
                          variant="ghost"
                          title="Delete shared application"
                          className="text-destructive hover:text-destructive"
                          onClick={() => onDelete(app)}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    ) : (
                      <div className="flex justify-end">
                        <Button type="button" size="icon" variant="ghost" title="Turn into shared application" onClick={() => onTurnCommon(app)}>
                          <Share2 className="h-4 w-4" />
                        </Button>
                      </div>
                    )}
                  </TableCell>
                </TableRow>
              )
            })
          )}
        </TableBody>
      </Table>
    </div>
  )
}
