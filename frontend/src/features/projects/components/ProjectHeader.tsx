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

import { ProjectSettingsDialog } from "./ProjectSettingsDialog";

interface ProjectHeaderProps {
  project: ProjectDto;
}

export function ProjectHeader({ project }: ProjectHeaderProps) {
  const [isDeployDialogOpen, setIsDeployDialogOpen] = useState(false);
  const [isSettingsDialogOpen, setIsSettingsDialogOpen] = useState(false);

  return (
    <div className="flex flex-col md:flex-row md:items-start justify-between gap-6 mb-10 border-b border-border pb-8 transition-all">
      <div className="flex flex-col gap-2">
        <div className="flex items-center gap-3">
          <h1 className="text-3xl md:text-4xl font-extrabold text-foreground tracking-tight">
            {project.name}
          </h1>
          <Badge className="bg-primary/10 text-primary border-primary/20 hover:bg-primary/20 transition-all font-bold px-3 py-1 shadow-[0_0_10px_var(--neon-brand-glow)]">
            Production
          </Badge>
        </div>
        <Link
          href={`${project.url}`}
          target="_blank"
          className="text-muted-foreground hover:text-primary hover:underline flex items-center gap-1.5 transition-all text-base group font-medium"
        >
          {project.url}
          <ExternalLink className="w-4 h-4 opacity-50 group-hover:opacity-100 transition-opacity" />
        </Link>
      </div>
      <div className="flex items-center gap-3">
        <Button variant="outline" asChild className="h-10 px-4 rounded-lg border-border bg-background text-foreground font-semibold text-sm hover:bg-accent hover:border-border/80 transition-all flex items-center gap-2 shadow-sm">
          <Link href={`/projects/${project.id}/analytics`}>
            <BarChart3 className="w-4 h-4 text-primary" />
            Analytics
          </Link>
        </Button>

        <Dialog open={isSettingsDialogOpen} onOpenChange={setIsSettingsDialogOpen}>
          <DialogTrigger asChild>
            <Button variant="outline" className="h-10 px-4 rounded-lg border-border bg-background text-foreground font-semibold text-sm hover:bg-accent hover:border-border/80 transition-all flex items-center gap-2 shadow-sm">
              <Settings className="w-4 h-4 text-muted-foreground" />
              Settings
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-md bg-card border-border shadow-2xl">
            <DialogHeader>
              <DialogTitle className="text-2xl font-extrabold text-foreground">Project Settings</DialogTitle>
            </DialogHeader>
            <ProjectSettingsDialog
              project={project}
              onSuccess={() => setIsSettingsDialogOpen(false)}
            />
          </DialogContent>
        </Dialog>

        <Dialog open={isDeployDialogOpen} onOpenChange={setIsDeployDialogOpen}>
          <DialogTrigger asChild>
            <Button className="h-10 px-5 rounded-lg bg-accent text-foreground font-bold text-sm shadow-sm hover:bg-accent/80 border border-border/50 transition-all flex items-center gap-2 hover:scale-[1.02]">
              <Upload className="w-4 h-4 text-primary" />
              Deploy
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-xl bg-card border-border shadow-2xl">
            <DialogHeader>
              <DialogTitle className="text-2xl font-extrabold text-foreground">Deploy New Version</DialogTitle>
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

        <Button className="h-10 px-5 rounded-lg bg-primary text-primary-foreground font-extrabold text-sm shadow-neon-brand hover:bg-primary/90 transition-all flex items-center gap-2 hover:scale-[1.02]" asChild>
          <Link href={`${project.url}`} target="_blank">
            Visit Site
            <ArrowUpRight className="w-4 h-4" />
          </Link>
        </Button>
      </div>
    </div>
  );
}
