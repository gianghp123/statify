import { Folder, BarChart3, MemoryStick, TrendingUp } from "lucide-react";

export function StatCards() {
  const stats = [
    {
      title: "Total Projects",
      value: "12",
      change: "+1",
      icon: Folder,
      increasing: true,
    },
    {
      title: "Total Bandwidth",
      value: "45",
      unit: "GB",
      change: "+12%",
      icon: BarChart3,
      increasing: true,
    },
    {
      title: "Current Builds",
      value: "2",
      description: "Processing",
      icon: MemoryStick,
      increasing: false,
    },
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
      {stats.map((stat) => (
        <div
          key={stat.title}
          className="bg-card rounded-xl p-5 border border-border hover:border-white/10 transition-colors shadow-lg shadow-black/20"
        >
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-sm font-medium text-muted-foreground">{stat.title}</h3>
            <stat.icon className="text-primary bg-primary/10 p-1.5 rounded-lg w-8 h-8" />
          </div>
          <div className="flex items-end gap-3">
            <span className="text-3xl font-bold text-white">
              {stat.value}
              {stat.unit && (
                <span className="text-lg font-normal text-muted-foreground">
                  {stat.unit}
                </span>
              )}
            </span>
            {stat.change && (
              <span className="text-sm text-primary font-medium mb-1 flex items-center">
                <TrendingUp className="w-4 h-4 mr-1" />
                {stat.change}
              </span>
            )}
            {stat.description && (
              <span className="text-sm text-muted-foreground font-medium mb-1">
                {stat.description}
              </span>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}
