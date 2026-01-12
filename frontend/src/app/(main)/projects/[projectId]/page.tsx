import { ChevronRight, Folder } from "lucide-react";
import Link from "next/link";
import { ProjectHeader } from "@/features/projects/components/ProjectHeader";
import { ProjectOverview } from "@/features/projects/components/ProjectOverview";
import { DeploymentHistoryTable } from "@/features/deployments/components/DeploymentHistoryTable";
import { getProjectById } from "@/features/projects/services/project.get";
import { getDeployments } from "@/features/deployments/services/deployment.get";
import { notFound } from "next/navigation";

interface ProjectDetailsPageProps {
  params: Promise<{ projectId: string }>;
}

export default async function ProjectDetailsPage({ params }: ProjectDetailsPageProps) {
  const { projectId } = await params;
  const id = parseInt(projectId);

  if (isNaN(id)) {
    return notFound();
  }

  const [projectRes, deploymentsRes] = await Promise.all([
    getProjectById(id),
    getDeployments(id),
  ]);

  if (!projectRes.success || !projectRes.data) {
    return notFound();
  }

  const project = projectRes.data;
  const deployments = deploymentsRes.success ? deploymentsRes.data || [] : [];

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
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <div className="lg:col-span-2 space-y-8">
          <ProjectOverview project={project} />
          <DeploymentHistoryTable deployments={deployments} projectId={project.id} />
        </div>
      </div>

      <footer className="mt-12 mb-6 flex justify-center text-xs text-muted-foreground/50">
        <p>© 2024 Statify Inc. Custom Deep Plum Edition.</p>
      </footer>
    </div>
  );
}
