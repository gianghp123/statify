import { ProjectOverviewMetricsDTO } from "@/features/analytics/dtos/response/analytic-metrics.dto";
import { Activity, BarChart3, Clock, Users } from "lucide-react";


interface AnalyticsStatsProps {
  metrics: ProjectOverviewMetricsDTO;
}

export function AnalyticsStats({ metrics }: AnalyticsStatsProps) {
  const stats = [
    {
      title: "Total Requests",
      value: metrics.totalRequests.toLocaleString(),
      icon: Users,
    },
    {
      title: "Bandwidth Used",
      value: `${metrics.totalBandwidth.toFixed(2)} MB`,
      icon: BarChart3,
    },
    {
      title: "Avg. Response Time",
      value: `${metrics.avgResponseMs.toFixed(0)}ms`,
      icon: Clock,
    },
    {
      title: "Error Rate",
      value: `${metrics.errorRatePercent.toFixed(2)}%`,
      icon: Activity,
    },
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      {stats.map((stat) => (
        <div key={stat.title} className="bg-card rounded-xl p-5 border border-border shadow-lg shadow-black/20">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-sm font-medium text-muted-foreground">{stat.title}</h3>
            <div className="size-8 rounded-lg bg-primary/10 flex items-center justify-center text-primary">
              <stat.icon className="w-4 h-4" />
            </div>
          </div>
          <div className="flex items-end justify-between">
            <span className="text-2xl font-bold text-foreground">{stat.value}</span>
          </div>
        </div>
      ))}
    </div>
  );
}
