"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Globe, Plus, MoreHorizontal, ShieldCheck, ExternalLink } from "lucide-react";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";

export default function DomainsSettings() {
  const domains = [
    { name: "myproject.statify.app", type: "System", status: "Active", primary: true },
    { name: "www.myawesomeproject.com", type: "Custom", status: "Verified", primary: false },
    { name: "blog.myawesomeproject.com", type: "Custom", status: "Pending", primary: false },
  ];

  return (
    <div className="max-w-4xl space-y-8">
      <div>
        <h1 className="text-foreground tracking-tight text-3xl font-bold">Domains</h1>
        <p className="text-muted-foreground mt-1">
          Manage your custom domains and subdomains.
        </p>
      </div>

      <Card className="glass-panel border-border/40">
        <CardHeader>
          <CardTitle>Add Custom Domain</CardTitle>
          <CardDescription>
            Connect your own domain name to your project.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex gap-2">
            <Input placeholder="example.com" className="bg-background/50 border-border/40" />
            <Button className="shadow-neon shrink-0">
              <Plus className="size-4 mr-2" />
              Add Domain
            </Button>
          </div>
        </CardContent>
      </Card>

      <div className="space-y-4">
        {domains.map((domain) => (
          <Card key={domain.name} className="glass-panel border-border/40 overflow-hidden">
            <div className="flex items-center justify-between p-6">
              <div className="flex items-center gap-4">
                <div className="size-10 rounded-full bg-primary/10 flex items-center justify-center text-primary">
                  <Globe className="size-5" />
                </div>
                <div>
                  <div className="flex items-center gap-2">
                    <span className="font-bold text-foreground">{domain.name}</span>
                    {domain.primary && (
                      <Badge className="bg-primary text-primary-foreground text-[10px] uppercase font-bold px-1.5 h-4">Primary</Badge>
                    )}
                  </div>
                  <p className="text-xs text-muted-foreground">{domain.type} Domain</p>
                </div>
              </div>

              <div className="flex items-center gap-4">
                <div className="hidden md:flex flex-col items-end">
                  <div className="flex items-center gap-1.5">
                    <div className={`size-1.5 rounded-full ${domain.status === 'Pending' ? 'bg-warning' : 'bg-success shadow-[0_0_8px_var(--success)]'}`} />
                    <span className={`text-xs font-bold ${domain.status === 'Pending' ? 'text-warning' : 'text-success'}`}>{domain.status}</span>
                  </div>
                  {domain.status === 'Verified' && (
                    <span className="text-[10px] text-muted-foreground flex items-center gap-0.5">
                      <ShieldCheck className="size-3" /> SSL Active
                    </span>
                  )}
                </div>

                <div className="flex items-center gap-2">
                  <Button variant="ghost" size="icon" className="hover:text-primary">
                    <ExternalLink className="size-4" />
                  </Button>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="icon" className="hover:text-primary">
                        <MoreHorizontal className="size-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem>Set as Primary</DropdownMenuItem>
                      <DropdownMenuItem>View DNS Settings</DropdownMenuItem>
                      <DropdownMenuItem className="text-destructive focus:text-destructive">Remove Domain</DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              </div>
            </div>
            {domain.status === 'Pending' && (
              <div className="bg-warning/5 border-t border-warning/10 p-4">
                <p className="text-xs text-warning leading-relaxed">
                  <strong>Action Required:</strong> Point your domain's CNAME record to <code>cname.statify.app</code> to complete verification.
                </p>
              </div>
            )}
          </Card>
        ))}
      </div>
    </div>
  );
}
