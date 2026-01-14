import { Button } from "@/components/ui/button";


export function DeploymentTable() {
  return (
    <div className="flex flex-col overflow-hidden rounded-xl border border-border bg-card/50 backdrop-blur-sm shadow-2xl">
      <table className="w-full text-left border-collapse">
        <thead>
          <tr className="border-b border-border bg-card text-xs uppercase tracking-wider text-muted-foreground">
            <th className="px-6 py-4 font-semibold w-32">Status</th>
            <th className="px-6 py-4 font-semibold">Commit</th>
            <th className="hidden lg:table-cell px-6 py-4 font-semibold">Project</th>
            <th className="hidden md:table-cell px-6 py-4 font-semibold w-24">Author</th>
            <th className="hidden sm:table-cell px-6 py-4 font-semibold w-32">Time</th>
            <th className="px-6 py-4 font-semibold w-24 text-right">Actions</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-border">
          <tr>
            <td colSpan={6} className="px-6 py-10 text-center text-muted-foreground">
              No recent deployments found.
            </td>
          </tr>
        </tbody>
      </table>
      <div className="border-t border-border bg-card px-6 py-4 flex items-center justify-between">
        <p className="text-sm text-muted-foreground">
          Showing <span className="text-foreground font-medium">1-4</span> of <span className="text-foreground font-medium">248</span> deployments
        </p>
        <div className="flex gap-2">
          <Button variant="outline" className="h-8 px-3 text-xs" disabled>Previous</Button>
          <Button variant="outline" className="h-8 px-3 text-xs">Next</Button>
        </div>
      </div>
    </div>
  );
}
