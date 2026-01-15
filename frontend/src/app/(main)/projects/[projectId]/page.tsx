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
  const pagination = deploymentsRes.success ? deploymentsRes.pagination || { totalCount: 0, page: 1, limit: 10 } : { totalCount: 0, page: 1, limit: 10 };

  return (
    <div className="space-y-8 animate-in fade-in duration-500">
      {/* Breadcrumbs */}
      <nav className="flex items-center gap-2 text-sm mb-6 border-b border-border pb-4">
        <Link
          href="/dashboard"
          className="text-muted-foreground hover:text-primary transition-all flex items-center gap-1.5 font-bold"
        >
          <Folder className="w-4 h-4" />
          Projects
        </Link>
        <ChevronRight className="w-4 h-4 text-muted-foreground/30" />
        <span className="text-foreground font-black px-2 py-0.5 bg-muted/50 rounded border border-border/50">{project.name}</span>
      </nav>

      <ProjectHeader project={project} />
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-10">
        <div className="lg:col-span-2 space-y-12">
          <ProjectOverview project={project} />
          <div className="pt-4">
            <DeploymentHistoryTable initialDeployments={deployments} projectId={project.id} pagination={pagination} />
          </div>
        </div>
      </div>
    </div>
  );
}
