"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { AnalyticsCharts } from "@/features/analytics/components/AnalyticsCharts";
import { AnalyticsStats } from "@/features/analytics/components/AnalyticsStats";
import { AnalyticMetricsDto, AnalyticMetricsSchema } from "@/features/analytics/dtos/response/analytic-metrics.dto";
import { streamAnalytics } from "@/features/analytics/services/analytic.sse";
import { ProjectDto } from "@/features/projects/dtos/response/project.response.dto";
import {
  Calendar
} from "lucide-react";
import { useEffect, useState } from "react";
import { toast } from "sonner";


interface AnalyticsPageProps {
  initialData: AnalyticMetricsDto;
  project: ProjectDto;
}

export function AnalyticsPage({ initialData, project }: AnalyticsPageProps) {
  const [analyticsData, setAnalyticsData] = useState<AnalyticMetricsDto>(initialData);

  useEffect(() => {
    const eventSource = streamAnalytics(project.id)

    eventSource.onmessage = (event) => {
      const data = AnalyticMetricsSchema.parse(JSON.parse(event.data)) as AnalyticMetricsDto
      setAnalyticsData(data)
    }

    eventSource.onerror = (error) => {
      toast.error("Failed to stream analytics")
    }

    return () => {
      eventSource.close()
    }

  }, [])

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
            <span className="text-primary">{project.name.toUpperCase()}</span> - INSIGHTS<span className="text-primary">.</span>
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
          {/* <Button size="sm" variant="default" className="rounded-full px-4 font-bold shadow-lg shadow-primary/20">
            <Download className="w-4 h-4 mr-2" /> Export
          </Button> */}
        </div>
      </div>

      {/* Primary Stats - Highlighting Static Strengths */}
      <AnalyticsStats metrics={analyticsData.projectOverviewMetrics} />

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
            <AnalyticsCharts data={analyticsData.timeSeriesPoints} />
          </CardContent>
        </Card>
      </div>
    </div>
  );
}