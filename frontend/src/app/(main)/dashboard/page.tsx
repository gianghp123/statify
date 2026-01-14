import { Button } from "@/components/ui/button";
import { logout } from "@/features/auth/services/auth.actions";
import { getCurrentUser } from "@/features/auth/services/auth.get";
import { ProjectList } from "@/features/projects/components/ProjectList";
import { StatCards } from "@/features/projects/components/StatCards";
import { getProjects } from "@/features/projects/services/project.get";
import { DeploymentStatus } from "@/lib/enums/deployment-status.enum";
import { Plus } from "lucide-react";
import Link from "next/link";

export default async function DashboardPage({
  searchParams,
}: {
  searchParams?: { [key: string]: string | string[] | undefined };
}) {
  const status = (await searchParams)?.status as DeploymentStatus || '';

  const [projectsRes, userRes] = await Promise.all([
    getProjects(1, 10, status),
    getCurrentUser(),
  ]);

  const projects = projectsRes.success ? projectsRes.data || [] : [];
  if (!userRes.success || !userRes.data?.username) {
    await logout()
  }
  const userName = userRes.data?.username

  return (
    <div className="space-y-8">
      <div className="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
          <h2 className="text-3xl font-bold text-foreground mb-2 tracking-tight">
            Welcome back, {userName}
          </h2>
          <p className="text-muted-foreground">
            Here is an overview of your projects and deployments.
          </p>
        </div>
        <Link href="/projects/new">
          <Button className="flex items-center gap-2 bg-primary text-primary-foreground px-5 py-2.5 rounded-lg font-bold text-sm shadow-neon-primary hover:shadow-neon hover:-translate-y-0.5 transition-all">
            <Plus className="w-5 h-5" />
            Create New Project
          </Button>
        </Link>
      </div>

      <StatCards projects={projects} />

      <div className="flex flex-wrap items-center gap-3 pb-2 border-b border-border">
        <Button asChild className="bg-white/10 text-foreground px-3 py-1.5 rounded-lg text-sm font-medium hover:bg-white/20 transition-colors border border-white/5">
          <Link href="/dashboard">All Projects</Link>
        </Button>
        <Button asChild className="bg-white/10 text-muted-foreground px-3 py-1.5 rounded-lg text-sm font-medium hover:text-foreground hover:bg-white/5 transition-colors">
          <Link href={`/dashboard?status=${DeploymentStatus.LIVE}`}>Live</Link>
        </Button>
        <Button asChild className="bg-white/10 text-muted-foreground px-3 py-1.5 rounded-lg text-sm font-medium hover:text-foreground hover:bg-white/5 transition-colors">
          <Link href={`/dashboard?status=${DeploymentStatus.READY}`}>Ready</Link>
        </Button>
        <Button asChild className="bg-white/10 text-muted-foreground px-3 py-1.5 rounded-lg text-sm font-medium hover:text-foreground hover:bg-white/5 transition-colors">
          <Link href={`/dashboard?status=${DeploymentStatus.FAILED}`}>Failed</Link>
        </Button>
      </div>

      <ProjectList projects={projects} />

      <footer className="mt-12 mb-6 flex justify-center text-xs text-muted-foreground/50">
        <p>© 2026 Statify Inc. Custom Deep Plum Edition.</p>
      </footer>
    </div>
  );
}
