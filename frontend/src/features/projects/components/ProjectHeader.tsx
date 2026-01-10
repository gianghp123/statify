'use client';
import { ExternalLink, Settings, ArrowUpRight, BarChart3, Upload } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import Link from "next/link";
import { ProjectDto } from "../dtos/response/project.response.dto";
import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { UploadDeploymentForm } from "../../deployments/components/UploadDeploymentForm";

interface ProjectHeaderProps {
  project: ProjectDto;
}

export function ProjectHeader({ project }: ProjectHeaderProps) {
  const [isDeployDialogOpen, setIsDeployDialogOpen] = useState(false);

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
        <Button variant="outline" asChild className="h-10 px-4 rounded-lg border-border bg-card text-foreground font-medium text-sm hover:bg-accent hover:text-white transition-all flex items-center gap-2">
          <Link href={`/projects/${project.id}/analytics`}>
            <BarChart3 className="w-4 h-4" />
            Analytics
          </Link>
        </Button>
        <Button variant="outline" className="h-10 px-4 rounded-lg border-border bg-card text-foreground font-medium text-sm hover:bg-accent hover:text-white transition-all flex items-center gap-2">
          <Settings className="w-4 h-4" />
          Settings
        </Button>

        <Dialog open={isDeployDialogOpen} onOpenChange={setIsDeployDialogOpen}>
          <DialogTrigger asChild>
            <Button className="h-10 px-5 rounded-lg bg-accent text-accent-foreground font-bold text-sm shadow-neon-accent hover:shadow-neon hover:bg-accent/90 transition-all flex items-center gap-2">
              <Upload className="w-4 h-4" />
              Deploy
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-xl bg-card border-border">
            <DialogHeader>
              <DialogTitle className="text-2xl font-black text-white">Deploy New Version</DialogTitle>
            </DialogHeader>
            <div className="py-4">
              <UploadDeploymentForm 
                projectId={project.id} 
                onSuccess={() => setIsDeployDialogOpen(false)}
                onCancel={() => setIsDeployDialogOpen(false)}
              />
            </div>
          </DialogContent>
        </Dialog>

        <Button className="h-10 px-5 rounded-lg bg-primary text-primary-foreground font-bold text-sm shadow-neon-brand hover:shadow-neon transition-all flex items-center gap-2" asChild>
          <Link href={`https://${project.url}`} target="_blank">
            Visit Site
            <ArrowUpRight className="w-4 h-4" />
          </Link>
        </Button>
      </div>
    </div>
  );
}
