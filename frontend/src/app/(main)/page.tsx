import { Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { StatCards } from "@/features/projects/components/StatCards";
import { ProjectList } from "@/features/projects/components/ProjectList";
import Link from "next/link";

export default function DashboardPage() {
  return (
    <div className="space-y-8">
      <div className="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div>
          <h2 className="text-3xl font-bold text-white mb-2 tracking-tight">
            Welcome back, Alex
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

      <StatCards />

      <div className="flex flex-wrap items-center gap-3 pb-2 border-b border-border">
        <button className="bg-white/10 text-white px-3 py-1.5 rounded-lg text-sm font-medium hover:bg-white/20 transition-colors border border-white/5">
          All Projects
        </button>
        <button className="text-muted-foreground px-3 py-1.5 rounded-lg text-sm font-medium hover:text-white hover:bg-white/5 transition-colors">
          Live
        </button>
        <button className="text-muted-foreground px-3 py-1.5 rounded-lg text-sm font-medium hover:text-white hover:bg-white/5 transition-colors">
          Building
        </button>
        <button className="text-muted-foreground px-3 py-1.5 rounded-lg text-sm font-medium hover:text-white hover:bg-white/5 transition-colors">
          Offline
        </button>
      </div>

      <ProjectList />

      <footer className="mt-12 mb-6 flex justify-center text-xs text-muted-foreground/50">
        <p>© 2024 Statify Inc. Custom Deep Plum Edition.</p>
      </footer>
    </div>
  );
}
