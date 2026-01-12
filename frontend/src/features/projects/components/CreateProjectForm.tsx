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
        <div className="space-y-2">
          <Label htmlFor="name" className="text-white font-semibold">Project Name</Label>
          <Input
            id="name"
            placeholder="my-awesome-site"
            value={projectName}
            onChange={(e) => setProjectName(e.target.value)}
            className="bg-card border-border text-white focus:ring-primary h-12 text-lg"
            required
            disabled={isPending}
          />
          <p className="text-sm text-muted-foreground flex items-center gap-2">
            <Globe className="w-4 h-4 text-primary" />
            Your site will be available at: 
            <span className="text-white font-mono bg-white/5 px-2 py-0.5 rounded border border-white/10">
              {subdomain}.statify.app
            </span>
          </p>
        </div>
      </div>

      <div className="flex items-center gap-4 pt-4">
        <Button 
          type="button" 
          variant="ghost" 
          onClick={() => router.back()}
          className="text-muted-foreground hover:text-white"
          disabled={isPending}
        >
          Cancel
        </Button>
        <Button 
          type="submit" 
          disabled={!projectName || isPending}
          className="bg-primary text-primary-foreground font-bold px-8 shadow-neon hover:shadow-neon-strong transition-all"
        >
          {isPending ? (
            <>
              <Loader2 className="w-4 h-4 mr-2 animate-spin" />
              Creating Project...
            </>
          ) : (
            "Create Project"
          )}
        </Button>
      </div>
    </form>
  );
}
