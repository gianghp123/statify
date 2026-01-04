import { ProjectCard } from "./ProjectCard";
import { ProjectDto } from "../dtos/response/project.response.dto";

interface ProjectListProps {
  projects: ProjectDto[];
}

export function ProjectList({ projects }: ProjectListProps) {
  if (!projects || projects.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-20 bg-card/50 rounded-xl border border-dashed border-border text-center">
        <p className="text-muted-foreground mb-4">No projects found.</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 gap-4">
      {projects.map((project) => (
        <ProjectCard key={project.id} project={project} />
      ))}
    </div>
  );
}
