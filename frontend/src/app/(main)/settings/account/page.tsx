"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";

export default function AccountSettings() {
  return (
    <div className="max-w-4xl space-y-8">
      <div>
        <h1 className="text-foreground tracking-tight text-3xl font-bold">Account Settings</h1>
        <p className="text-muted-foreground mt-1">
          Manage your personal information and account preferences.
        </p>
      </div>

      <Card className="glass-panel border-border/40">
        <CardHeader>
          <CardTitle>Public Profile</CardTitle>
          <CardDescription>
            How you appear to other members.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="flex items-center gap-6">
            <Avatar className="size-20 border-2 border-primary/20">
              <AvatarImage src="" />
              <AvatarFallback className="bg-primary/10 text-primary text-xl font-bold">JD</AvatarFallback>
            </Avatar>
            <div className="space-y-2">
              <Button variant="outline" size="sm" className="border-border/40">Change Avatar</Button>
              <p className="text-xs text-muted-foreground">JPG, GIF or PNG. Max size of 800K</p>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="username">Username</Label>
              <Input id="username" defaultValue="johndoe" className="bg-background/50 border-border/40" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="email">Email Address</Label>
              <Input id="email" type="email" defaultValue="john@example.com" className="bg-background/50 border-border/40" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="fullname">Full Name</Label>
              <Input id="fullname" defaultValue="John Doe" className="bg-background/50 border-border/40" />
            </div>
            <div className="space-y-2">
              <Label htmlFor="location">Location</Label>
              <Input id="location" placeholder="San Francisco, CA" className="bg-background/50 border-border/40" />
            </div>
          </div>
        </CardContent>
        <CardFooter className="justify-end border-t border-border/40 pt-6">
          <Button className="shadow-neon">Update Profile</Button>
        </CardFooter>
      </Card>

      <Card className="glass-panel border-border/40">
        <CardHeader>
          <CardTitle>Preferences</CardTitle>
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
    </div>
  );
}
