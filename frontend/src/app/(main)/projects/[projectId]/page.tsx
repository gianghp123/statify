"use client";

import { use } from "react";
import { ChevronRight, Folder } from "lucide-react";
import Link from "next/link";
import { ProjectHeader } from "@/features/projects/components/ProjectHeader";
import { ProjectOverview } from "@/features/projects/components/ProjectOverview";
import { DeploymentHistoryTable } from "@/features/deployments/components/DeploymentHistoryTable";
import { mockProjects } from "@/lib/mock";

interface ProjectDetailsPageProps {
  params: Promise<{ projectId: string }>;
}

export default function ProjectDetailsPage({ params }: ProjectDetailsPageProps) {
  const { projectId } = use(params);
  
  // Find project from mock data (or default to the first one)
  const project = mockProjects.find(p => String(p.id) === projectId) || mockProjects[0];

  return (
    <div className="space-y-6">
      {/* Breadcrumbs */}
      <nav className="flex items-center gap-2 text-sm mb-6">
        <Link 
          href="/" 
          className="text-muted-foreground hover:text-primary transition-colors flex items-center gap-1"
        >
          <Folder className="w-4 h-4" />
          Projects
        </Link>
        <ChevronRight className="w-4 h-4 text-muted-foreground/50" />
        <span className="text-white font-medium">{project.name}</span>
      </nav>

      <ProjectHeader project={project} />
      <ProjectOverview project={project} />
      <DeploymentHistoryTable />

      <footer className="mt-12 mb-6 flex justify-center text-xs text-muted-foreground/50">
        <p>© 2024 Statify Inc. Custom Deep Plum Edition.</p>
      </footer>
    </div>
  );
}
