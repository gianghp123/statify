"use client";

import { use } from "react";
import { Download, Calendar } from "lucide-react";
import { Button } from "@/components/ui/button";
import { AnalyticsStats } from "@/features/analytics/components/AnalyticsStats";
import { AnalyticsCharts } from "@/features/analytics/components/AnalyticsCharts";
import { ContentTable } from "@/features/analytics/components/ContentTable";

interface AnalyticsPageProps {
  params: Promise<{ projectId: string }>;
}

export default function AnalyticsPage({ params }: AnalyticsPageProps) {
  const { projectId } = use(params);
  
  return (
    <div className="space-y-8">
      <div className="flex flex-col gap-2">
        <div className="flex flex-col md:flex-row md:items-end justify-between gap-4 mt-2">
          <div>
            <h1 className="text-white tracking-tight text-3xl md:text-4xl font-bold leading-tight">
              Website Analytics
            </h1>
            <p className="text-muted-foreground mt-1 text-sm md:text-base">
              Real-time data monitoring and performance tracking.
            </p>
          </div>
          <div className="flex items-center gap-2 bg-card border border-border rounded-lg p-1 pr-4 shadow-xl">
             <Button variant="ghost" size="sm" className="hover:bg-primary/10 hover:text-primary">
               &larr;
             </Button>
             <div className="flex items-center gap-2 text-white text-sm font-medium px-2">
               <Calendar className="w-4 h-4 text-primary" />
               <span>Oct 01, 2023 - Oct 31, 2023</span>
             </div>
             <Button variant="ghost" size="sm" className="hover:bg-primary/10 hover:text-primary">
               &rarr;
             </Button>
          </div>
        </div>
      </div>

      <AnalyticsStats />
      
      <div className="glass-card rounded-xl p-6 border border-border shadow-xl space-y-6">
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div className="flex items-center gap-4">
            <h3 className="text-white text-lg font-bold">Traffic Overview</h3>
            <div className="flex bg-background/50 rounded-lg p-1 border border-border">
              <button className="px-3 py-1 rounded text-xs font-bold bg-primary text-primary-foreground shadow-neon">Visits</button>
              <button className="px-3 py-1 rounded text-xs font-medium text-muted-foreground hover:text-white hover:bg-white/5 transition-colors">Bandwidth</button>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <div className="flex bg-background/50 rounded-lg p-1 border border-border">
              <button className="px-3 py-1 rounded text-xs font-medium text-muted-foreground hover:text-white hover:bg-white/5 transition-colors">1D</button>
              <button className="px-3 py-1 rounded text-xs font-bold bg-primary/20 text-primary border border-primary/50">1W</button>
              <button className="px-3 py-1 rounded text-xs font-medium text-muted-foreground hover:text-white hover:bg-white/5 transition-colors">1M</button>
            </div>
            <Button variant="outline" className="h-8 px-3 text-xs flex items-center gap-2 border-primary/50 text-primary hover:bg-primary/10 transition-all font-bold">
              <Download className="w-3.5 h-3.5" />
              Export
            </Button>
          </div>
        </div>
        
        <AnalyticsCharts />
      </div>

      <ContentTable />

      <footer className="mt-12 mb-6 flex justify-center text-xs text-muted-foreground/50">
        <p>© 2024 Statify Inc. Custom Deep Plum Edition.</p>
      </footer>
    </div>
  );
}
