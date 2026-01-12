import { Folder, BarChart3, MemoryStick, TrendingUp } from "lucide-react";
import { ProjectDto } from "../dtos/response/project.response.dto";

interface StatCardsProps {
  projects: ProjectDto[];
}

export function StatCards({ projects }: StatCardsProps) {
  const totalProjects = projects.length;
  // For other stats, we might need more APIs, but let's at least make this one real.
  
  const stats = [
    {
      title: "Total Projects",
      value: totalProjects.toString(),
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
      value: "0",
      description: "Idle",
      icon: MemoryStick,
      increasing: false,
    },
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
      {stats.map((stat) => (
        <div
          key={stat.title}
          className="bg-card rounded-2xl p-6 border border-border hover:border-border/80 transition-all shadow-xl dark:shadow-black/30 group hover:scale-[1.02]"
        >
          <div className="flex items-center justify-between mb-5">
            <h3 className="text-xs font-bold text-muted-foreground uppercase tracking-wider">{stat.title}</h3>
            <div className="p-3 rounded-2xl bg-primary/10 group-hover:bg-primary/20 transition-colors shadow-sm border border-primary/10">
              <stat.icon className="text-primary w-6 h-6" />
            </div>
          </div>
          <div className="flex items-end gap-3">
            <span className="text-4xl font-black text-foreground tracking-tighter">
              {stat.value}
              {stat.unit && (
                <span className="text-xl font-bold text-muted-foreground ml-1">
                  {stat.unit}
                </span>
              )}
            </span>
            {stat.change && (
              <span className="text-sm text-primary font-black mb-1.5 flex items-center bg-primary/10 px-2 py-0.5 rounded-full border border-primary/20 shadow-sm shadow-primary/10">
                <TrendingUp className="w-3.5 h-3.5 mr-1" />
                {stat.change}
              </span>
            )}
            {stat.description && (
              <span className="text-xs text-muted-foreground font-black mb-2 flex items-center bg-muted px-2 py-0.5 rounded-full border border-border uppercase tracking-tighter">
                {stat.description}
              </span>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}
