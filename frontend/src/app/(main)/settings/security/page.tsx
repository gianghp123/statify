"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { Shield, Smartphone, Key, Globe } from "lucide-react";

export default function SecuritySettings() {
  return (
    <div className="max-w-4xl space-y-8">
      <div>
        <h1 className="text-foreground tracking-tight text-3xl font-bold">Security</h1>
        <p className="text-muted-foreground mt-1">
          Manage your account security and authentication methods.
        </p>
      </div>

      <Card className="glass-panel border-border/40">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Key className="size-5 text-primary" />
            Change Password
          </CardTitle>
          <CardDescription>
            Ensure your account is using a long, random password to stay secure.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="current">Current Password</Label>
            <Input id="current" type="password" placeholder="••••••••" className="bg-background/50 border-border/40" />
          </div>
          <div className="space-y-2">
            <Label htmlFor="new">New Password</Label>
            <Input id="new" type="password" placeholder="••••••••" className="bg-background/50 border-border/40" />
          </div>
          <div className="space-y-2">
            <Label htmlFor="confirm">Confirm New Password</Label>
            <Input id="confirm" type="password" placeholder="••••••••" className="bg-background/50 border-border/40" />
          </div>
        </CardContent>
        <CardFooter className="justify-end border-t border-border/40 pt-6">
          <Button className="shadow-neon">Update Password</Button>
        </CardFooter>
      </Card>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card className="glass-panel border-border/40">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <Smartphone className="size-5 text-primary" />
              Two-Factor Auth
            </CardTitle>
            <CardDescription className="text-xs">
              Add an extra layer of security.
            </CardDescription>
          </CardHeader>
          <CardContent className="flex items-center justify-between">
            <span className="text-sm text-muted-foreground font-medium">Currently Disabled</span>
            <Switch />
          </CardContent>
        </Card>

        <Card className="glass-panel border-border/40">
          <CardHeader>
            <CardTitle className="flex items-center gap-2 text-lg">
              <Shield className="size-5 text-primary" />
              Recovery Codes
            </CardTitle>
            <CardDescription className="text-xs">
              Backup access methods.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button variant="outline" size="sm" className="w-full border-border/40">Generate New Codes</Button>
          </CardContent>
        </Card>
      </div>

      <Card className="glass-panel border-border/40">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Globe className="size-5 text-primary" />
            Active Sessions
          </CardTitle>
          <CardDescription>
            Places where you are logged in.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center justify-between p-4 rounded-lg bg-accent/30 border border-border/40">
            <div className="flex items-center gap-3">
              <div className="size-10 rounded-full bg-primary/10 flex items-center justify-center text-primary">
                <Globe className="size-5" />
              </div>
              <div>
                <p className="text-sm font-medium text-foreground">Chrome on MacOS</p>
                <p className="text-xs text-muted-foreground">San Francisco, USA • Active now</p>
              </div>
            </div>
            <Button variant="ghost" size="sm" className="text-muted-foreground hover:text-destructive">Revoke</Button>
          </div>
          <div className="flex items-center justify-between p-4 rounded-lg bg-accent/30 border border-border/40">
            <div className="flex items-center gap-3">
              <div className="size-10 rounded-full bg-primary/10 flex items-center justify-center text-primary">
                <Smartphone className="size-5" />
              </div>
              <div>
                <p className="text-sm font-medium text-foreground">Statify Mobile (iPhone 15)</p>
                <p className="text-xs text-muted-foreground">Saigon, Vietnam • 2 days ago</p>
              </div>
            </div>
            <Button variant="ghost" size="sm" className="text-muted-foreground hover:text-destructive">Revoke</Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
