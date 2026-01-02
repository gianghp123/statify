import { ExternalLink, Settings, ArrowUpRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import Link from "next/link";
import { ProjectDto } from "../dtos/response/project.response.dto";

interface ProjectHeaderProps {
  project: ProjectDto;
}

export function ProjectHeader({ project }: ProjectHeaderProps) {
  return (
    <div className="flex flex-col md:flex-row md:items-start justify-between gap-6 mb-10 border-b border-border pb-8">
      <div className="flex flex-col gap-2">
        <div className="flex items-center gap-3">
          <h1 className="text-3xl md:text-4xl font-bold text-white tracking-tight">
            {project.name}
          </h1>
          <Badge className="bg-primary/10 text-primary border-primary/20 hover:bg-primary/20 transition-colors">
            Production
          </Badge>
        </div>
        <Link
          href={`https://${project.url}`}
          target="_blank"
          className="text-muted-foreground hover:text-primary hover:underline flex items-center gap-1.5 transition-colors text-base group"
        >
          {project.url}
          <ExternalLink className="w-4 h-4 opacity-50 group-hover:opacity-100 transition-opacity" />
        </Link>
      </div>
      <div className="flex items-center gap-3">
        <Button variant="outline" className="h-10 px-4 rounded-lg border-border bg-card text-foreground font-medium text-sm hover:bg-accent hover:text-white transition-all flex items-center gap-2">
          <Settings className="w-4 h-4" />
          Settings
        </Button>
        <Button className="h-10 px-5 rounded-lg bg-primary text-primary-foreground font-bold text-sm shadow-neon-primary hover:shadow-neon transition-all flex items-center gap-2">
          Visit Site
          <ArrowUpRight className="w-4 h-4" />
        </Button>
      </div>
    </div>
  );
}
