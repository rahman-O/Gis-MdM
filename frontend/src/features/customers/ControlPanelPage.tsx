import { useEffect, useState } from 'react'
import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/shared/ui/table'
import { useToast } from '@/shared/hooks/use-toast'
import * as customersService from '@/features/customers/customersService'

export function ControlPanelPage() {
  const { toast } = useToast()
  const [q, setQ] = useState('')
  const [rows, setRows] = useState<customersService.CustomerRow[]>([])
  const [err, setErr] = useState<string | null>(null)

  useEffect(() => {
    customersService
      .searchCustomers({ currentPage: 1, pageSize: 100, searchValue: q.trim() })
      .then((p) => setRows(customersService.unwrapCustomerRows(p)))
      .catch((reason: unknown) => setErr(reason instanceof Error ? reason.message : String(reason)))
  }, [q])

  async function loginAs(customerId: number) {
    try {
      const user = await customersService.impersonateCustomer(customerId)
      customersService.hydrateSessionAfterImpersonation(user)
      toast({
        title: 'Impersonated (if backend returned tokens)',
        description: user.authToken ? 'Session refreshed from response.' : 'You may still need the legacy servlet session.',
      })
      window.location.href = '/dashboard'
    } catch {
      toast({ title: 'Impersonation failed', variant: 'destructive' })
    }
  }

  return (
    <div className="space-y-4">
      <h1 className="text-xl font-semibold tracking-tight">Control panel</h1>
      <Input placeholder="Search…" value={q} onChange={(e) => setQ(e.target.value)} className="max-w-sm" />
      {err ? <p className="text-destructive text-sm">{err}</p> : null}
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>ID</TableHead>
            <TableHead>Name</TableHead>
            <TableHead></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {rows.map((r) => (
            <TableRow key={String(r.id)}>
              <TableCell>{r.id}</TableCell>
              <TableCell>{r.name ?? r.email ?? '—'}</TableCell>
              <TableCell className="text-right">
                <Button type="button" size="sm" variant="secondary" disabled={!r.id} onClick={() => loginAs(r.id!)}>
                  Open as admin
                </Button>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}
