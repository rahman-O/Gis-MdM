import { MoreHorizontal } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/shared/ui/dropdown-menu'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/shared/ui/table'
import type { Application } from '@/features/applications/model/types'

interface Props {
  applications: Application[]
  onEdit: (app: Application) => void
  onDelete: (app: Application) => void
  onVersions: (app: Application) => void
  onConfigurations: (app: Application) => void
}

export function ApplicationsTable({
  applications,
  onEdit,
  onDelete,
  onVersions,
  onConfigurations,
}: Props) {
  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Package</TableHead>
            <TableHead>Name</TableHead>
            <TableHead>Version</TableHead>
            <TableHead>URL</TableHead>
            <TableHead>Type</TableHead>
            <TableHead>Icon</TableHead>
            <TableHead className="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {applications.length === 0 ? (
            <TableRow>
              <TableCell colSpan={7} className="h-24 text-center text-muted-foreground">
                No applications found.
              </TableCell>
            </TableRow>
          ) : applications.map((app) => (
            <TableRow key={app.id ?? `${app.pkg}-${app.version}`}>
              <TableCell>{app.pkg ?? '—'}</TableCell>
              <TableCell className="font-medium">{app.name ?? '—'}</TableCell>
              <TableCell>{app.version ?? app.latestVersionText ?? '—'}</TableCell>
              <TableCell className="max-w-[240px] truncate">{app.url ?? '—'}</TableCell>
              <TableCell>{String(app.type ?? 'app')}</TableCell>
              <TableCell>{app.showIcon ? 'Yes' : 'No'}</TableCell>
              <TableCell className="text-right">
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button size="icon" variant="ghost"><MoreHorizontal className="h-4 w-4" /></Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuItem onSelect={() => onConfigurations(app)}>Configurations</DropdownMenuItem>
                    <DropdownMenuItem onSelect={() => onVersions(app)}>Versions</DropdownMenuItem>
                    <DropdownMenuItem onSelect={() => onEdit(app)}>Edit</DropdownMenuItem>
                    <DropdownMenuItem
                      onSelect={() => onDelete(app)}
                      className="text-destructive focus:text-destructive"
                    >
                      Delete
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
