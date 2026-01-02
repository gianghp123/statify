import { CreateProjectForm } from "@/features/projects/components/CreateProjectForm";
import { ChevronLeft } from "lucide-react";
import Link from "next/link";

export default function NewProjectPage() {
  return (
    <div className="space-y-8 py-6">
      <div className="space-y-2">
        <Link 
          href="/" 
          className="text-muted-foreground hover:text-primary transition-colors flex items-center gap-1 text-sm mb-4"
        >
          <ChevronLeft className="w-4 h-4" />
          Back to Dashboard
        </Link>
        <h1 className="text-white text-3xl md:text-4xl font-black tracking-tight">
          Create New Project
        </h1>
        <p className="text-muted-foreground text-lg">
          Deploy your static site to the cloud in seconds.
        </p>
      </div>

      <div className="bg-card/30 border border-border/50 rounded-2xl p-8 backdrop-blur-sm shadow-2xl">
        <CreateProjectForm />
      </div>
    </div>
  );
}
