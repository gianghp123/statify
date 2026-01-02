import { MoreHorizontal, Rocket, FlaskConical, CheckCircle2, History } from "lucide-react";

const getStatusBadge = (status: string) => {
  switch (status) {
    case "Ready":
      return (
        <span className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 px-2.5 py-0.5 text-xs font-medium text-primary border border-primary/20">
          <span className="size-1.5 rounded-full bg-primary"></span>
          Ready
        </span>
      );
    case "Building":
      return (
        <span className="inline-flex items-center gap-1.5 rounded-full bg-yellow-400/10 px-2.5 py-0.5 text-xs font-medium text-yellow-400 border border-yellow-400/20">
          <span className="size-1.5 rounded-full bg-yellow-400 animate-pulse"></span>
          Building
        </span>
      );
    case "Failed":
      return (
        <span className="inline-flex items-center gap-1.5 rounded-full bg-red-400/10 px-2.5 py-0.5 text-xs font-medium text-red-400 border border-red-400/20">
          <span className="size-1.5 rounded-full bg-red-400"></span>
          Failed
        </span>
      );
    default:
      return null;
  }
};

const mockDeployments = [
  {
    id: "5f3a2b1",
    commitMessage: "Update hero copy",
    status: "Ready",
    context: "Production",
    createdAt: "2m ago",
    duration: "45s",
  },
  {
    id: "8c9d1e4",
    commitMessage: "Fix typo in footer",
    status: "Building",
    context: "Preview",
    createdAt: "5m ago",
    duration: "--",
  },
  {
    id: "2a1b5c9",
    commitMessage: "Add analytics script",
    status: "Failed",
    context: "Preview",
    createdAt: "2h ago",
    duration: "1m 12s",
  },
];

export function DeploymentHistoryTable() {
  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center justify-between">
        <h3 className="text-white font-semibold text-lg flex items-center gap-2">
          <History className="text-primary w-5 h-5" />
          Deployment History
        </h3>
      </div>
      <div className="w-full overflow-hidden rounded-xl border border-border bg-card shadow-sm">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead>
              <tr className="border-b border-border bg-background/30 text-xs uppercase tracking-wider text-muted-foreground font-semibold">
                <th className="px-6 py-4">Deployment ID</th>
                <th className="px-6 py-4">Status</th>
                <th className="px-6 py-4">Context</th>
                <th className="px-6 py-4">Created</th>
                <th className="px-6 py-4">Duration</th>
                <th className="px-6 py-4 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {mockDeployments.map((deployment) => (
                <tr key={deployment.id} className="group hover:bg-white/[0.02] transition-colors">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center gap-2">
                      <span className="font-mono text-white text-xs bg-background border border-border px-1.5 py-0.5 rounded">
                        {deployment.id}
                      </span>
                      <span className="text-muted-foreground truncate max-w-[150px] text-xs">
                        {deployment.commitMessage}
                      </span>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    {getStatusBadge(deployment.status)}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-foreground">
                    <div className="flex items-center gap-1.5">
                      {deployment.context === "Production" ? (
                        <Rocket className="w-4 h-4 text-muted-foreground" />
                      ) : (
                        <FlaskConical className="w-4 h-4 text-muted-foreground" />
                      )}
                      {deployment.context}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-muted-foreground">
                    {deployment.createdAt}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-muted-foreground">
                    {deployment.duration}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right">
                    <button className="text-muted-foreground hover:text-white transition-colors p-1 rounded hover:bg-white/5">
                      <MoreHorizontal className="w-4 h-4" />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        <div className="border-t border-border bg-background/20 px-6 py-3 flex items-center justify-between">
          <span className="text-xs text-muted-foreground">
            Showing 1-{mockDeployments.length} of 24 deployments
          </span>
          <div className="flex items-center gap-2">
            <button className="p-1 rounded text-muted-foreground hover:text-white hover:bg-border disabled:opacity-50" disabled>
              <span className="text-xs">Prev</span>
            </button>
            <button className="p-1 rounded text-muted-foreground hover:text-white hover:bg-border">
              <span className="text-xs">Next</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
