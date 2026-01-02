import { ProjectCard } from "./ProjectCard";
import { mockProjects } from "../../../lib/mock";

export function ProjectList() {
  return (
    <div className="grid grid-cols-1 gap-4">
      {mockProjects.map((project) => (
        <ProjectCard key={project.id} project={project} />
      ))}
    </div>
  );
}
