import { Button } from '@/shared/ui/button'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'
import { Checkbox } from '@/shared/ui/checkbox'
import type { Configuration, ConfigurationFile } from '@/features/configurations/types'

interface Props {
  configuration: Configuration
  onChange: (next: Configuration) => void
}

export function ConfigurationFilesTab({ configuration, onChange }: Props) {
  const files = Array.isArray(configuration.files) ? configuration.files : []

  const updateFile = (index: number, patch: Partial<ConfigurationFile>) => {
    const next = [...files]
    next[index] = { ...next[index], ...patch }
    onChange({ ...configuration, files: next })
  }

  const addFile = () => {
    onChange({
      ...configuration,
      files: [...files, { path: '', externalUrl: '', url: '', remove: false }],
    })
  }

  const removeFile = (index: number) => {
    onChange({
      ...configuration,
      files: files.filter((_, idx) => idx !== index),
    })
  }

  return (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label>Default file path</Label>
        <Input
          value={String(configuration.defaultFilePath ?? '/')}
          onChange={(e) => onChange({ ...configuration, defaultFilePath: e.target.value })}
        />
      </div>

      <Button type="button" variant="outline" onClick={addFile}>
        Add file link
      </Button>

      {files.length === 0 ? (
        <p className="text-sm text-muted-foreground">No files linked.</p>
      ) : files.map((file, index) => (
        <div key={`file-${index}`} className="rounded-md border p-3">
          <div className="grid gap-3 md:grid-cols-2">
            <div className="space-y-2">
              <Label>Path</Label>
              <Input
                value={String(file.path ?? '')}
                onChange={(e) => updateFile(index, { path: e.target.value })}
              />
            </div>
            <div className="space-y-2">
              <Label>External URL</Label>
              <Input
                value={String(file.externalUrl ?? '')}
                onChange={(e) => updateFile(index, { externalUrl: e.target.value })}
              />
            </div>
            <div className="space-y-2 md:col-span-2">
              <Label>Resolved URL</Label>
              <Input
                value={String(file.url ?? '')}
                onChange={(e) => updateFile(index, { url: e.target.value })}
              />
            </div>
          </div>
          <div className="mt-3 flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Checkbox
                checked={Boolean(file.remove)}
                onCheckedChange={(checked) => updateFile(index, { remove: checked === true })}
              />
              <Label>Mark for removal</Label>
            </div>
            <Button type="button" size="sm" variant="destructive" onClick={() => removeFile(index)}>
              Remove
            </Button>
          </div>
        </div>
      ))}
    </div>
  )
}
