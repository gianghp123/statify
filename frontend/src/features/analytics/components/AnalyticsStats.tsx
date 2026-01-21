import { BarChart3, TrendingUp, Users, Clock, ArrowUpRight, ArrowDownRight, Activity } from "lucide-react";
import { ProjectOverviewMetricsDTO } from "@/features/analytics/dtos/response/analytic-metrics.dto";


interface AnalyticsStatsProps {
  metrics: ProjectOverviewMetricsDTO;
}

export function AnalyticsStats({ metrics }: AnalyticsStatsProps) {
  const stats = [
    {
      title: "Total Requests",
      value: metrics.totalRequests.toLocaleString(),
      change: "+12.5%", // TODO: Calculate change if previous data available
      increasing: true,
      icon: Users,
    },
    {
      title: "Bandwidth Used",
      value: `${metrics.totalBandwidth.toFixed(2)} MB`,
      change: "+18.2%",
      increasing: true,
      icon: BarChart3,
    },
    {
      title: "Avg. Response Time",
      value: `${metrics.avgResponseMs.toFixed(0)}ms`,
      change: "-4ms",
      increasing: false, // lower is better for latency, but typically green means good. Adjust logic if needed.
      icon: Clock,
    },
    {
      title: "Error Rate",
      value: `${metrics.errorRatePercent.toFixed(2)}%`,
      change: "+0.01%",
      increasing: false, // lower is better
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
