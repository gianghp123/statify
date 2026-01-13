"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { Bell, Mail, Smartphone, Monitor } from "lucide-react";

export default function NotificationsSettings() {
  const notificationGroups = [
    {
      title: "Email Notifications",
      description: "Get updates delivered right to your inbox.",
      icon: Mail,
      items: [
        { label: "Deployment Success", description: "Notify me when a project is successfully deployed." },
        { label: "Deployment Failed", description: "Notify me immediately on deployment errors." },
        { label: "Monthly Reports", description: "A summary of your projects performance." },
      ]
    },
    {
      title: "Push Notifications",
      description: "Stay updated with real-time browser alerts.",
      icon: Monitor,
      items: [
        { label: "Critical Alerts", description: "High-priority system notifications." },
        { label: "New Features", description: "Updates about new Statify tools and features." },
      ]
    }
  ];

  return (
    <div className="max-w-4xl space-y-8">
      <div>
        <h1 className="text-foreground tracking-tight text-3xl font-bold">Notifications</h1>
        <p className="text-muted-foreground mt-1">
          Choose how you want to be notified about your project activities.
        </p>
      </div>

      {notificationGroups.map((group, idx) => (
        <Card key={idx} className="glass-panel border-border/40">
          <CardHeader>
            <div className="flex items-center gap-3">
              <div className="size-10 rounded-lg bg-primary/10 flex items-center justify-center text-primary">
                <group.icon className="size-5" />
              </div>
              <div>
                <CardTitle>{group.title}</CardTitle>
                <CardDescription>{group.description}</CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent className="space-y-6 pt-2">
            {group.items.map((item, itemIdx) => (
              <div key={itemIdx} className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label className="text-sm font-bold text-foreground">{item.label}</Label>
                  <p className="text-xs text-muted-foreground">{item.description}</p>
                </div>
                <Switch defaultChecked={itemIdx < 2} />
              </div>
            ))}
          </CardContent>
        </Card>
      ))}

      <Card className="glass-panel border-border/40">
        <CardHeader>
          <div className="flex items-center gap-3">
            <div className="size-10 rounded-lg bg-primary/10 flex items-center justify-center text-primary">
              <Smartphone className="size-5" />
            </div>
            <div>
              <CardTitle>Mobile App</CardTitle>
              <CardDescription>Configure notifications for your mobile device.</CardDescription>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">
            No mobile devices linked to this account. Download the Statify app to get started.
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
