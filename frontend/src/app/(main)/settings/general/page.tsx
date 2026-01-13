"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Separator } from "@/components/ui/separator";

export default function GeneralSettings() {
  return (
    <div className="max-w-4xl space-y-8">
      <div>
        <h1 className="text-foreground tracking-tight text-3xl font-bold">General Settings</h1>
        <p className="text-muted-foreground mt-1">
          Update your project details and preferences.
        </p>
      </div>

      <Card className="glass-panel border-border/40">
        <CardHeader>
          <CardTitle>Project Profile</CardTitle>
          <CardDescription>
            This information will be displayed publicly.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="name">Project Name</Label>
            <Input id="name" placeholder="My Awesome Project" className="bg-background/50 border-border/40" />
          </div>
          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>
            <Textarea
              id="description"
              placeholder="A brief description of your project."
              className="min-h-[100px] bg-background/50 border-border/40"
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="theme">Default Theme</Label>
            <Select defaultValue="system">
              <SelectTrigger id="theme" className="bg-background/50 border-border/40">
                <SelectValue placeholder="Select a theme" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="light">Light</SelectItem>
                <SelectItem value="dark">Dark</SelectItem>
                <SelectItem value="system">System</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </CardContent>
        <CardFooter className="justify-end border-t border-border/40 pt-6">
          <Button className="shadow-neon">Save Changes</Button>
        </CardFooter>
      </Card>

      <Card className="border-destructive/20 bg-destructive/5">
        <CardHeader>
          <CardTitle className="text-destructive">Danger Zone</CardTitle>
          <CardDescription className="text-destructive/70">
            Permanently delete this project and all of its data. This action cannot be undone.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-destructive/80 font-medium">
            Deleting this project will stop all deployments and remove all associated data.
          </p>
        </CardContent>
        <CardFooter className="justify-end border-t border-destructive/10 pt-6">
          <Button variant="destructive">Delete Project</Button>
        </CardFooter>
      </Card>
    </div>
  );
}
