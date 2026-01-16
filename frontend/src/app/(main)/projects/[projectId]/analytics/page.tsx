"use client";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { AnalyticsCharts } from "@/features/analytics/components/AnalyticsCharts";
import { ContentTable } from "@/features/analytics/components/ContentTable";
import {
  BarChart3,
  Calendar,
  Clock,
  Download,
  Globe,
  Zap
} from "lucide-react";
import { use } from "react";

interface AnalyticsPageProps {
  params: Promise<{ projectId: string }>;
}

export default function AnalyticsPage({ params }: AnalyticsPageProps) {
  const { projectId } = use(params);

  return (
    <div className="max-w-[1400px] mx-auto space-y-6 p-4 md:p-8">
      {/* Header Section - Sleeker & More Functional */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-6">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <span className="flex h-2 w-2 rounded-full bg-emerald-500 animate-pulse" />
            <span className="text-[10px] font-bold uppercase tracking-widest text-muted-foreground">Live Snapshot</span>
          </div>
          <h1 className="text-3xl font-extrabold tracking-tight text-foreground italic">
            INSIGHTS<span className="text-primary">.</span>
          </h1>
        </div>

        <div className="flex items-center gap-3">
          <div className="hidden lg:flex flex-col items-end mr-4">
            <span className="text-[10px] text-muted-foreground uppercase font-semibold">Last Build</span>
            <span className="text-xs font-mono">2 mins ago (v1.0.4)</span>
          </div>
          <div className="flex items-center bg-secondary/50 border border-border rounded-full px-4 py-1.5 shadow-sm">
            <Calendar className="w-3.5 h-3.5 text-primary mr-2" />
            <span className="text-xs font-semibold">Oct 01 — Oct 31</span>
          </div>
          <Button size="sm" variant="default" className="rounded-full px-4 font-bold shadow-lg shadow-primary/20">
            <Download className="w-4 h-4 mr-2" /> Export
          </Button>
        </div>
      </div>

      {/* Primary Stats - Highlighting Static Strengths */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatMiniCard label="Total Visits" value="42.5k" change="+12%" icon={<Globe className="w-4 h-4 text-blue-500" />} />
        <StatMiniCard label="Avg. Load Time" value="0.8s" change="-0.2s" icon={<Zap className="w-4 h-4 text-yellow-500" />} />
        <StatMiniCard label="Bounce Rate" value="24.2%" change="+2%" icon={<BarChart3 className="w-4 h-4 text-purple-500" />} />
        <StatMiniCard label="Session Duration" value="3m 12s" change="+45s" icon={<Clock className="w-4 h-4 text-emerald-500" />} />
      </div>

      {/* Main Content Grid */}
      <div className="flex flex-col gap-6">
        {/* Main Chart - Takes 2/3 width */}
        <Card className="lg:col-span-2 bg-card/50 backdrop-blur-sm border-border shadow-xl overflow-hidden">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-7">
            <div>
              <CardTitle className="text-lg font-bold">Traffic Distribution</CardTitle>
              <p className="text-xs text-muted-foreground">Daily unique visitors and page views.</p>
            </div>
            <div className="flex bg-muted rounded-md p-1 border border-border">
              <button className="px-3 py-1 text-[10px] font-bold uppercase bg-background rounded shadow-sm">1W</button>
              <button className="px-3 py-1 text-[10px] font-bold uppercase text-muted-foreground">1M</button>
              <button className="px-3 py-1 text-[10px] font-bold uppercase text-muted-foreground">1Y</button>
            </div>
          </CardHeader>
          <CardContent>
            <AnalyticsCharts />
          </CardContent>
        </Card>

        {/* Top Pages - Simplified for Static Sites */}
        <Card className="bg-card/50 backdrop-blur-sm border-border shadow-xl">
          <CardHeader>
            <CardTitle className="text-lg font-bold">Top Destinations</CardTitle>
          </CardHeader>
          <CardContent>
            <ContentTable />
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

function StatMiniCard({ label, value, change, icon }: { label: string, value: string, change: string, icon: React.ReactNode }) {
  const isPositive = change.startsWith('+');
  return (
    <div className="p-5 rounded-2xl bg-card border border-border shadow-sm hover:border-primary/50 transition-colors">
      <div className="flex items-center justify-between mb-3">
        <div className="p-2 rounded-lg bg-secondary/50 border border-border">{icon}</div>
        <span className={`text-[10px] font-bold px-2 py-0.5 rounded-full ${isPositive ? 'bg-emerald-500/10 text-emerald-500' : 'bg-red-500/10 text-red-500'}`}>
          {change}
        </span>
      </div>
      <div>
        <p className="text-xs text-muted-foreground font-medium">{label}</p>
        <h4 className="text-2xl font-bold tracking-tight">{value}</h4>
      </div>
    </div>
  );
}