import { Checkbox } from '@/shared/ui/checkbox'
import { Input } from '@/shared/ui/input'
import { Label } from '@/shared/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/shared/ui/select'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/shared/ui/card'
import type { Configuration } from '@/features/configurations/types'

interface Props {
  configuration: Configuration
  onChange: (next: Configuration) => void
}

export function ConfigurationDesignTab({ configuration, onChange }: Props) {
  // Safe helper to validate color for preview style
  const safeColor = (colorStr: string | null | undefined, fallback: string): string => {
    if (!colorStr) return fallback
    const trimmed = colorStr.trim()
    // Simple check to ensure we only apply CSS colors, hex, rgb, etc.
    if (trimmed.startsWith('#') || trimmed.startsWith('rgb') || trimmed.startsWith('hsl') || /^[a-z]+$/i.test(trimmed)) {
      return trimmed
    }
    return fallback
  }

  return (
    <div className="grid gap-6 md:grid-cols-2">
      {/* CARD 1: Style & Visual Theme */}
      <Card className="shadow-sm border">
        <CardHeader className="bg-muted/15 border-b py-3 px-4">
          <CardTitle className="text-sm font-bold uppercase tracking-wider text-muted-foreground">Style & Color Theme</CardTitle>
          <CardDescription className="text-xs mt-0.5">Control corporate themes, backgrounds, and element colors.</CardDescription>
        </CardHeader>
        <CardContent className="p-4 grid gap-4">
          <div className="flex items-center space-x-2.5 rounded-md border p-3 hover:bg-muted/20 transition-all duration-200">
            <Checkbox
              id="useDefaultDesignSettings"
              checked={Boolean(configuration.useDefaultDesignSettings)}
              onCheckedChange={(checked) =>
                onChange({ ...configuration, useDefaultDesignSettings: checked === true })
              }
            />
            <Label htmlFor="useDefaultDesignSettings" className="cursor-pointer text-xs font-bold uppercase tracking-wider text-muted-foreground flex-1">
              Use default design settings
            </Label>
          </div>

          {!Boolean(configuration.useDefaultDesignSettings) && (
            <>
              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-1.5">
                  <Label htmlFor="backgroundColor" className="text-xs font-semibold text-muted-foreground uppercase">Background Color</Label>
                  <div className="relative flex items-center">
                    <Input
                      id="backgroundColor"
                      placeholder="#FFFFFF"
                      value={String(configuration.backgroundColor ?? '')}
                      onChange={(e) => onChange({ ...configuration, backgroundColor: e.target.value })}
                      className="pl-9"
                    />
                    <div
                      className="absolute left-3 w-4 h-4 rounded-full border shadow-sm transition-colors duration-200"
                      style={{
                        backgroundColor: safeColor(configuration.backgroundColor, '#ffffff'),
                        borderColor: 'rgba(0,0,0,0.15)'
                      }}
                    />
                  </div>
                </div>

                <div className="space-y-1.5">
                  <Label htmlFor="textColor" className="text-xs font-semibold text-muted-foreground uppercase">Text Color</Label>
                  <div className="relative flex items-center">
                    <Input
                      id="textColor"
                      placeholder="#000000"
                      value={String(configuration.textColor ?? '')}
                      onChange={(e) => onChange({ ...configuration, textColor: e.target.value })}
                      className="pl-9"
                    />
                    <div
                      className="absolute left-3 w-4 h-4 rounded-full border shadow-sm transition-colors duration-200"
                      style={{
                        backgroundColor: safeColor(configuration.textColor, '#000000'),
                        borderColor: 'rgba(0,0,0,0.15)'
                      }}
                    />
                  </div>
                </div>
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="backgroundImageUrl" className="text-xs font-semibold text-muted-foreground uppercase">Background Image URL</Label>
                <Input
                  id="backgroundImageUrl"
                  placeholder="https://example.com/background.png"
                  value={String(configuration.backgroundImageUrl ?? '')}
                  onChange={(e) => onChange({ ...configuration, backgroundImageUrl: e.target.value })}
                />
              </div>

              <div className="flex items-center space-x-2.5 rounded-md border p-3 hover:bg-muted/20 transition-all duration-200">
                <Checkbox
                  id="displayStatus"
                  checked={Boolean(configuration.displayStatus)}
                  onCheckedChange={(checked) => onChange({ ...configuration, displayStatus: checked === true })}
                />
                <Label htmlFor="displayStatus" className="cursor-pointer text-xs font-bold uppercase tracking-wider text-muted-foreground flex-1">
                  Display status information
                </Label>
              </div>
            </>
          )}
          {Boolean(configuration.useDefaultDesignSettings) && (
            <div className="py-6 text-center text-xs text-muted-foreground border border-dashed rounded-md">
              Default system styles are currently enforced. Uncheck above to customize.
            </div>
          )}
        </CardContent>
      </Card>

      {/* CARD 2: Layout & Display Options */}
      <Card className="shadow-sm border">
        <CardHeader className="bg-muted/15 border-b py-3 px-4">
          <CardTitle className="text-sm font-bold uppercase tracking-wider text-muted-foreground">Layout & Displays</CardTitle>
          <CardDescription className="text-xs mt-0.5">Control icon size, orientations, and custom launcher header displays.</CardDescription>
        </CardHeader>
        <CardContent className="p-4 grid gap-4">
          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-1.5">
              <Label htmlFor="iconSize" className="text-xs font-semibold text-muted-foreground uppercase">Launcher Icon Size</Label>
              <Select
                value={String(configuration.iconSize ?? 'SMALL')}
                onValueChange={(value) => onChange({ ...configuration, iconSize: value })}
              >
                <SelectTrigger id="iconSize"><SelectValue placeholder="Select size" /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="SMALL">Small Icons</SelectItem>
                  <SelectItem value="MEDIUM">Medium Icons</SelectItem>
                  <SelectItem value="LARGE">Large Icons</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-1.5">
              <Label htmlFor="orientation" className="text-xs font-semibold text-muted-foreground uppercase">Screen Orientation</Label>
              <Select
                value={String(configuration.orientation ?? 'AUTO')}
                onValueChange={(value) => onChange({ ...configuration, orientation: value })}
              >
                <SelectTrigger id="orientation"><SelectValue placeholder="Select orientation" /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="AUTO">Auto Rotate</SelectItem>
                  <SelectItem value="PORTRAIT">Locked Portrait</SelectItem>
                  <SelectItem value="LANDSCAPE">Locked Landscape</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="border-t pt-4 grid gap-4">
            <div className="space-y-1.5">
              <Label htmlFor="desktopHeader" className="text-xs font-semibold text-muted-foreground uppercase">Launcher Header Mode</Label>
              <Select
                value={String(configuration.desktopHeader ?? 'NO_HEADER')}
                onValueChange={(value) => onChange({ ...configuration, desktopHeader: value })}
              >
                <SelectTrigger id="desktopHeader"><SelectValue placeholder="Select mode" /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="NO_HEADER">No header</SelectItem>
                  <SelectItem value="DEVICE_ID">Show Device ID</SelectItem>
                  <SelectItem value="DESCRIPTION">Show Description</SelectItem>
                  <SelectItem value="TEMPLATE">Use Custom Template</SelectItem>
                  <SelectItem value="CUSTOM1">Show Custom Field 1</SelectItem>
                  <SelectItem value="CUSTOM2">Show Custom Field 2</SelectItem>
                  <SelectItem value="CUSTOM3">Show Custom Field 3</SelectItem>
                </SelectContent>
              </Select>
            </div>

            {String(configuration.desktopHeader ?? 'NO_HEADER') === 'TEMPLATE' && (
              <div className="space-y-1.5">
                <Label htmlFor="desktopHeaderTemplate" className="text-xs font-semibold text-muted-foreground uppercase">Header Custom Template</Label>
                <Input
                  id="desktopHeaderTemplate"
                  placeholder="e.g. Device ID: ${deviceId}"
                  value={String(configuration.desktopHeaderTemplate ?? configuration.desktopHeaderText ?? '')}
                  onChange={(e) => onChange({ ...configuration, desktopHeaderTemplate: e.target.value })}
                />
              </div>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

