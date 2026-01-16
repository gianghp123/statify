'use client'
import { InfoTooltip } from "@/components/info-tooltip"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

export function PreferenceSettingCard() {
  return <Card className="glass-panel border-border/40">
    <CardHeader>
      <CardTitle className="flex items-center gap-2">Preferences <InfoTooltip content="This section is on development. Please come back later." /></CardTitle>
      <CardDescription>
        Customize your account experience.
      </CardDescription>
    </CardHeader>
    <CardContent className="space-y-4">
      <div className="flex items-center justify-between p-4 rounded-lg bg-accent/30 border border-border/40">
        <div>
          <p className="text-sm font-medium text-foreground">Language</p>
          <p className="text-xs text-muted-foreground">Select your preferred language.</p>
        </div>
        <Button variant="ghost" className="text-primary hover:text-primary hover:bg-primary/10">English (US)</Button>
      </div>
      <div className="flex items-center justify-between p-4 rounded-lg bg-accent/30 border border-border/40">
        <div>
          <p className="text-sm font-medium text-foreground">Timezone</p>
          <p className="text-xs text-muted-foreground">Used for logs and deployment timestamps.</p>
        </div>
        <Button variant="ghost" className="text-primary hover:text-primary hover:bg-primary/10">(GMT-08:00) Pacific Time</Button>
      </div>
    </CardContent>
  </Card>
} 