"use client";

import { useState } from "react";
import { Globe, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useRouter } from "next/navigation";
import { createProject } from "../services/project.actions";
import { toast } from "sonner";

export function CreateProjectForm() {
  const router = useRouter();
  const [projectName, setProjectName] = useState("");
  const [isPending, setIsPending] = useState(false);

  const subdomain = projectName
    .toLowerCase()
    .replace(/[^a-z0-9]/g, "-")
    .replace(/-+/g, "-")
    .replace(/^-|-$/g, "") || "your-project";

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!projectName) return;

    setIsPending(true);
    try {
      const projectRes = await createProject({
        name: projectName,
        subdomain: subdomain,
      });

      if (!projectRes.success || !projectRes.data) {
        toast.error(projectRes.message || "Failed to create project");
        return;
      }

      const projectId = projectRes.data.id;
      toast.success("Project created successfully!");
      router.push(`/projects/${projectId}`);
    } catch (error) {
      console.error(error);
      toast.error("An unexpected error occurred");
    } finally {
      setIsPending(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-8">
      <div className="space-y-6">
        {/* Project Name */}
        <div className="space-y-4">
          <Label htmlFor="name" className="text-foreground font-bold text-lg">Project Name</Label>
          <Input
            id="name"
            placeholder="my-awesome-site"
            value={projectName}
            onChange={(e) => setProjectName(e.target.value)}
            className="bg-background border-border text-foreground focus:ring-primary h-14 text-xl font-bold rounded-xl shadow-sm transition-all focus:border-primary/50"
            required
            disabled={isPending}
          />
          <div className="flex items-center gap-3 p-4 rounded-xl bg-muted/30 border border-border shadow-inner group transition-all hover:bg-muted/50">
            <Globe className="w-5 h-5 text-primary group-hover:scale-110 transition-transform" />
            <div className="flex flex-col">
              <span className="text-xs text-muted-foreground font-bold uppercase tracking-wider">Deployment URL</span>
              <span className="text-foreground font-mono font-bold">
                {subdomain}.statify.app
              </span>
            </div>
          </div>
        </div>
      </div>

      <div className="flex items-center gap-4 pt-4">
        <Button 
          type="button" 
          variant="ghost" 
          onClick={() => router.back()}
          className="text-muted-foreground hover:text-foreground font-bold h-12 px-6"
          disabled={isPending}
        >
          Cancel
        </Button>
        <Button 
          type="submit" 
          disabled={!projectName || isPending}
          className="bg-primary text-primary-foreground font-black h-12 px-10 rounded-xl shadow-neon-brand hover:scale-[1.02] transition-all flex-1"
        >
          {isPending ? (
            <>
              <Loader2 className="w-5 h-5 mr-2 animate-spin" />
              Creating...
            </>
          ) : (
            "Create Project"
          )}
        </Button>
      </div>
    </form>
  );
}
