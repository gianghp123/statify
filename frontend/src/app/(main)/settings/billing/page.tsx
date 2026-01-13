"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { CreditCard, Check, Zap, Download } from "lucide-react";

export default function BillingSettings() {
  const invoices = [
    { id: "INV-001", date: "Jan 1, 2024", amount: "$19.00", status: "Paid" },
    { id: "INV-002", date: "Dec 1, 2023", amount: "$19.00", status: "Paid" },
    { id: "INV-003", date: "Nov 1, 2023", amount: "$19.00", status: "Paid" },
  ];

  return (
    <div className="max-w-4xl space-y-8">
      <div>
        <h1 className="text-foreground tracking-tight text-3xl font-bold">Billing & Subscription</h1>
        <p className="text-muted-foreground mt-1">
          Manage your subscription plan, payment methods, and billing history.
        </p>
      </div>

      <Card className="glass-panel border-primary/20 relative overflow-hidden">
        <div className="absolute top-0 right-0 p-6 opacity-10">
          <Zap className="size-32 text-primary" />
        </div>
        <CardHeader>
          <div className="flex items-center justify-between">
            <Badge variant="outline" className="text-primary border-primary/30 bg-primary/10">Current Plan</Badge>
            <span className="text-2xl font-bold">$19<span className="text-sm text-muted-foreground font-normal">/mo</span></span>
          </div>
          <CardTitle className="text-2xl pt-2">Pro Developer</CardTitle>
          <CardDescription>
            For serious developers building production applications.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <ul className="grid grid-cols-1 md:grid-cols-2 gap-2 text-sm text-muted-foreground">
            <li className="flex items-center gap-2"><Check className="size-4 text-primary" /> Unlimited Projects</li>
            <li className="flex items-center gap-2"><Check className="size-4 text-primary" /> 500GB Bandwidth</li>
            <li className="flex items-center gap-2"><Check className="size-4 text-primary" /> Custom Domains</li>
            <li className="flex items-center gap-2"><Check className="size-4 text-primary" /> Priority Support</li>
          </ul>
        </CardContent>
        <CardFooter className="border-t border-border/40 pt-6">
          <Button className="w-full md:w-auto shadow-neon">Upgrade Plan</Button>
          <Button variant="ghost" className="ml-0 md:ml-4 mt-2 md:mt-0">Cancel Subscription</Button>
        </CardFooter>
      </Card>

      <Card className="glass-panel border-border/40">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <CreditCard className="size-5 text-primary" />
            Payment Methods
          </CardTitle>
          <CardDescription>
            Your current primary payment method.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-between p-4 rounded-xl bg-accent/30 border border-border/40">
            <div className="flex items-center gap-4">
              <div className="size-12 rounded-lg bg-card border border-border/40 flex items-center justify-center">
                <span className="font-bold text-primary">VISA</span>
              </div>
              <div>
                <p className="text-sm font-bold text-foreground">•••• •••• •••• 4242</p>
                <p className="text-xs text-muted-foreground">Expires 12/26</p>
              </div>
            </div>
            <Button variant="outline" size="sm" className="border-border/40">Update</Button>
          </div>
        </CardContent>
      </Card>

      <Card className="glass-panel border-border/40">
        <CardHeader>
          <CardTitle>Billing History</CardTitle>
          <CardDescription>
            Download your previous invoices.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow className="hover:bg-transparent border-border/40">
                <TableHead>Invoice ID</TableHead>
                <TableHead>Date</TableHead>
                <TableHead>Amount</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="text-right"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {invoices.map((invoice) => (
                <TableRow key={invoice.id} className="border-border/40">
                  <TableCell className="font-medium">{invoice.id}</TableCell>
                  <TableCell className="text-muted-foreground text-xs">{invoice.date}</TableCell>
                  <TableCell className="text-foreground font-medium">{invoice.amount}</TableCell>
                  <TableCell>
                    <Badge variant="secondary" className="bg-success/10 text-success border-success/20 text-[10px] uppercase font-bold tracking-wider">
                      {invoice.status}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-right">
                    <Button variant="ghost" size="icon" className="hover:text-primary">
                      <Download className="size-4" />
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
