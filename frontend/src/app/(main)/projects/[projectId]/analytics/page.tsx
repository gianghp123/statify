import { BreadScrum } from "@/components/bread-scrum";
import { AnalyticsPage } from "@/features/analytics/components/AnalyticsPage"
import { getAnalyticMetrics } from "@/features/analytics/services/analytic.get";
import { getProjectById } from "@/features/projects/services/project.get";
import { ChevronRight, Folder } from "lucide-react";
import Link from "next/link";
import { notFound } from "next/navigation";

interface PageProps {
  params: Promise<{ projectId: string }>;
}

export default async function Page({ params }: PageProps) {
  const { projectId } = await params;
  const projectIdNum = Number(projectId)
  const [initialData, projectRes] = await Promise.all([
    getAnalyticMetrics(projectIdNum),
    getProjectById(projectIdNum)
  ]);
  if (!initialData.success || !initialData.data || !projectRes.success || !projectRes.data) {
    return notFound()
  }

  const projectBreadcrumbItems = [
    { name: "Projects", href: "/dashboard" },
    { name: projectRes.data.name, href: `/projects/${projectIdNum}` },
    { name: "Analytics", href: `/projects/${projectIdNum}/analytics`, isCurrent: true },
  ]

  return (
    <div className="space-y-8 animate-in fade-in duration-500">
      {/* Breadcrumbs */}
      <BreadScrum items={projectBreadcrumbItems} />
      <AnalyticsPage initialData={initialData.data} project={projectRes.data} />
    </div>
  )
}