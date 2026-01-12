import { CreateProjectForm } from "@/features/projects/components/CreateProjectForm";
import { ChevronLeft } from "lucide-react";
import Link from "next/link";

export default function NewProjectPage() {
  return (
    <div className="space-y-8 py-6">
      <div className="space-y-2">
        <Link 
          href="/" 
          className="text-muted-foreground hover:text-primary transition-all flex items-center gap-1 text-sm mb-4 font-bold"
        >
          <ChevronLeft className="w-4 h-4" />
          Back to Dashboard
        </Link>
        <h1 className="text-foreground text-3xl md:text-5xl font-black tracking-tighter">
          Create New Project
        </h1>
        <p className="text-muted-foreground text-lg md:text-xl font-medium">
          Deploy your static site to the cloud in seconds.
        </p>
      </div>

      <div className="bg-card border border-border rounded-2xl p-8 shadow-2xl transition-all">
        <CreateProjectForm />
      </div>
    </div>
  );
}
