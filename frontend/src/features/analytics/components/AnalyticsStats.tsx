import { BarChart3, TrendingUp, Users, Clock, ArrowUpRight, ArrowDownRight } from "lucide-react";

export function AnalyticsStats() {
  const stats = [
    {
      title: "Total Visitors",
      value: "48.2k",
      change: "+12.5%",
      increasing: true,
      icon: Users,
    },
    {
      title: "Bandwidth Used",
      value: "124 GB",
      change: "+18.2%",
      increasing: true,
      icon: BarChart3,
    },
    {
      title: "Avg. Build Time",
      value: "42s",
      change: "-4s",
      increasing: true,
      icon: Clock,
    },
    {
      title: "Uptime",
      value: "99.98%",
      change: "+0.01%",
      increasing: true,
      icon: TrendingUp,
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
            <span className="text-2xl font-bold text-white">{stat.value}</span>
            <div className={`flex items-center text-xs font-medium ${stat.increasing ? 'text-primary' : 'text-red-400'}`}>
              {stat.increasing ? <ArrowUpRight className="w-3 h-3 mr-0.5" /> : <ArrowDownRight className="w-3 h-3 mr-0.5" />}
              {stat.change}
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
