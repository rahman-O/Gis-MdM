import { useCallback, useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { AlertCircle, MoreHorizontal } from 'lucide-react'
import { Button } from '@/shared/ui/button'
import { Skeleton } from '@/shared/ui/skeleton'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/shared/ui/table'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/shared/ui/dropdown-menu'
import * as configurationService from '@/features/configurations/configurationService'
import { hasPermission } from '@/features/auth/permissions'
import type { Configuration } from '@/features/configurations/types'
import { ConfigurationForm } from '@/features/configurations/ConfigurationForm'
import { ConfigurationDeleteDialog } from '@/features/configurations/ConfigurationDeleteDialog'

function typeLabel(type: number | null | undefined): string {
  return configurationService.typeToConfigurationKind(type) === 'COMMON' ? 'COMMON' : 'WORK'
}

function deviceCountCell(count: number | null | undefined): string {
  if (count == null) return '—'
  return String(count)
}

export function ConfigurationsPage() {
  const navigate = useNavigate()
  const [configurations, setConfigurations] = useState<Configuration[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [formOpen, setFormOpen] = useState(false)
  const [formMode, setFormMode] = useState<'create' | 'edit'>('create')
  const [selectedConfig, setSelectedConfig] = useState<Configuration | null>(null)
  const [configToDelete, setConfigToDelete] = useState<Configuration | null>(null)
  const canAdd = hasPermission('add_config')
  const canCopyDelete = hasPermission('copy_config')

  const loadList = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const list = await configurationService.getConfigurations()
      setConfigurations(Array.isArray(list) ? list : [])
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load configurations.')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void loadList()
  }, [loadList])

  const openCreate = () => {
    setFormMode('create')
    setSelectedConfig(null)
    setFormOpen(true)
  }

  const openEdit = (c: Configuration) => {
    if (c.id == null) return
    navigate(`/configurations/${c.id}/edit`)
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Configurations</h1>
          <p className="text-muted-foreground text-sm">MDM configurations for your devices.</p>
        </div>
        <Button type="button" onClick={openCreate} disabled={!canAdd}>
          New configuration
        </Button>
      </div>

      {error ?
        <div
          className="flex flex-wrap items-center gap-3 rounded-lg border border-destructive/50 bg-destructive/10 px-4 py-3 text-sm"
          role="alert"
        >
          <AlertCircle className="h-4 w-4 shrink-0 text-destructive" />
          <span className="flex-1">{error}</span>
          <Button type="button" variant="outline" size="sm" onClick={() => void loadList()}>
            Retry
          </Button>
        </div>
      : null}

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>Description</TableHead>
              <TableHead>Device count</TableHead>
              <TableHead className="w-[60px]">
                <span className="sr-only">Actions</span>
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ?
              Array.from({ length: 5 }).map((_, i) => (
                <TableRow key={`sk-${i}`}>
                  <TableCell colSpan={5}>
                    <Skeleton className="h-9 w-full" />
                  </TableCell>
                </TableRow>
              ))
            : configurations.length === 0 ?
              <TableRow>
                <TableCell colSpan={5} className="text-muted-foreground h-24 text-center">
                  No configurations found
                </TableCell>
              </TableRow>
            : configurations.map((c) => (
                <TableRow key={c.id ?? c.name}>
                  <TableCell className="font-medium">{c.name ?? '—'}</TableCell>
                  <TableCell>{typeLabel(c.type)}</TableCell>
                  <TableCell className="max-w-md truncate">{c.description?.trim() || '—'}</TableCell>
                  <TableCell>{deviceCountCell(c.deviceCount)}</TableCell>
                  <TableCell className="text-right">
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-8 w-8"
                          aria-label={`Actions for ${c.name ?? 'configuration'}`}
                        >
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem onSelect={() => openEdit(c)}>Edit</DropdownMenuItem>
                        <DropdownMenuItem
                          disabled={!canCopyDelete || c.id == null}
                          onSelect={async () => {
                            if (c.id == null) return
                            await configurationService.copyConfiguration({
                              id: c.id,
                              name: `${(c.name ?? 'Configuration').trim() || 'Configuration'} (Copy)`,
                              description: c.description ?? null,
                            })
                            await loadList()
                          }}
                        >
                          Copy
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          disabled={!canCopyDelete}
                          className="text-destructive focus:text-destructive"
                          onSelect={() => setConfigToDelete(c)}
                        >
                          Delete
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              ))
            }
          </TableBody>
        </Table>
      </div>

      <ConfigurationForm
        open={formOpen}
        onOpenChange={setFormOpen}
        mode={formMode}
        initialData={selectedConfig}
        onSuccess={() => void loadList()}
      />

      <ConfigurationDeleteDialog
        configuration={
          configToDelete?.id != null ?
            { id: configToDelete.id, name: configToDelete.name }
          : null
        }
        onConfirm={async () => {
          if (configToDelete?.id == null) throw new Error('Invalid configuration.')
          await configurationService.deleteConfiguration(configToDelete.id)
          await loadList()
        }}
        onCancel={() => setConfigToDelete(null)}
      />
    </div>
  )
}
