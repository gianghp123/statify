import { ChevronRight, Folder } from "lucide-react";
import Link from "next/link";
import { ProjectHeader } from "@/features/projects/components/ProjectHeader";
import { ProjectOverview } from "@/features/projects/components/ProjectOverview";
import { DeploymentHistoryTable } from "@/features/deployments/components/DeploymentHistoryTable";
import { getProjectById } from "@/features/projects/services/project.get";
import { getDeployments } from "@/features/deployments/services/deployment.get";
import { notFound } from "next/navigation";
import { BreadScrum } from "@/components/bread-scrum";


interface ProjectDetailsPageProps {
  params: Promise<{ projectId: string }>;
  searchParams: Promise<{ page?: string }>;
}

export default async function ProjectDetailsPage({ params, searchParams }: ProjectDetailsPageProps) {
  const { projectId } = await params;
  const { page: pageParam } = await searchParams;
  const id = parseInt(projectId);
  const page = pageParam ? parseInt(pageParam) : 1;

  if (isNaN(id)) {
    return notFound();
  }

  const [projectRes, deploymentsRes] = await Promise.all([
    getProjectById(id),
    getDeployments(id, page),
  ]);

  if (!projectRes.success || !projectRes.data) {
    return notFound();
  }

  const project = projectRes.data;
  const deployments = deploymentsRes.success ? deploymentsRes.data || [] : [];
  const pagination = deploymentsRes.success ? deploymentsRes.pagination || { totalCount: 0, page: 1, limit: 10 } : { totalCount: 0, page: 1, limit: 10 };

  const projectBreadcrumbItems = [
    { name: "Projects", href: "/dashboard" },
    { name: project.name, href: `/projects/${projectId}`, isCurrent: true },
  ]

  return (
    <div className="space-y-8 animate-in fade-in duration-500">
      <BreadScrum items={projectBreadcrumbItems} />

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
