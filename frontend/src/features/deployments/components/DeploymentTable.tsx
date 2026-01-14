import { Button } from "@/components/ui/button";
import { Pagination } from "@/lib/response/api-response";
import { formatDate } from "@/lib/utils/time.utils";
import { Eye } from "lucide-react";
import Link from "next/link";
import { DeploymentDto } from "../dtos/response/deployment.response.dto";
import { getStatusBadge } from "../utils/get-status-badge";

interface DeploymentTableProps {
  deployments: DeploymentDto[];
  pagination: Pagination;
}

export function DeploymentTable({ deployments, pagination }: DeploymentTableProps) {
  const totalPages = Math.ceil(pagination.totalCount / pagination.limit);

  return (
    <div className="flex flex-col overflow-hidden rounded-xl border border-border bg-card/50 backdrop-blur-sm shadow-2xl">
      <table className="w-full text-left border-collapse">
        <thead>
          <tr className="border-b border-border bg-card text-xs uppercase tracking-wider text-muted-foreground">
            <th className="px-6 py-4 font-semibold w-32">Status</th>
            <th className="px-6 py-4 font-semibold">ID</th>
            <th className="hidden lg:table-cell px-6 py-4 font-semibold">Project ID</th>
            <th className="hidden lg:table-cell px-6 py-4 font-semibold">Project Name</th>
            <th className="hidden lg:table-cell px-6 py-4 font-semibold">Project URL</th>
            <th className="hidden sm:table-cell px-6 py-4 font-semibold w-40">Created At</th>
            <th className="px-6 py-4 font-semibold w-24 text-right">Actions</th>
          </tr>
        </thead>
        <tbody className="divide-y divide-border">
          {deployments.length > 0 ? (
            deployments.map((deployment) => (
              <tr key={deployment.id} className="hover:bg-muted/50 transition-colors">
                <td className="px-6 py-4">
                  {getStatusBadge(deployment.status)}
                </td>
                <td className="px-6 py-4">
                  <div className="flex flex-col">
                    <span className="text-sm font-medium text-foreground">#{deployment.id}</span>
                  </div>
                </td>
                <td className="hidden sm:table-cell px-6 py-4 text-sm text-muted-foreground">
                  {deployment.project?.id}
                </td>
                <td className="hidden sm:table-cell px-6 py-4 text-sm text-muted-foreground">
                  {deployment.project?.name}
                </td>
                <td className="hidden lg:table-cell px-6 py-4 text-sm text-muted-foreground hover:underline hover:text-primary">
                  <Link href={`/projects/${deployment.project?.id}`}>{deployment.project?.subdomain}.{process.env.NEXT_PUBLIC_PROJECT_URL}</Link>
                </td>
                <td className="hidden sm:table-cell px-6 py-4 text-sm text-muted-foreground">
                  {formatDate(deployment.createdAt)}
                </td>
                <td className="px-6 py-4 text-right">
                  <Button variant="ghost" size="sm" className="h-8 text-xs hover:text-primary" asChild>
                    <Link href={`/projects/${deployment.project?.id}`}>
                      <Eye />
                    </Link>
                  </Button>
                </td>
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan={5} className="px-6 py-10 text-center text-muted-foreground">
                No recent deployments found.
              </td>
            </tr>
          )}
        </tbody>
      </table>

      {/* Pagination Footer */}
      <div className="border-t border-border bg-card px-6 py-4 flex items-center justify-between">
        <p className="text-sm text-muted-foreground">
          Showing <span className="text-foreground font-medium">
            {Math.min((pagination.page - 1) * pagination.limit + 1, pagination.totalCount)} - {Math.min(pagination.page * pagination.limit, pagination.totalCount)}
          </span> of <span className="text-foreground font-medium">{pagination.totalCount}</span> deployments
        </p>
        <div className="flex gap-2">
          <Button
            variant="outline"
            className="h-8 px-3 text-xs"
            disabled={pagination.page <= 1}
          >
            Previous
          </Button>
          <Button
            variant="outline"
            className="h-8 px-3 text-xs"
            disabled={pagination.page >= totalPages}
          >
            Next
          </Button>
        </div>
      </div>
    </div>
  );
}