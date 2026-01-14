import { DeploymentListHeader } from "@/features/deployments/components/DeploymentListHeader";
import { DeploymentTable } from "@/features/deployments/components/DeploymentTable";
import { getGlobalDeploymentHistory } from "@/features/deployments/services/deployment.get";

export default async function DeploymentsPage() {
  const deploymentRes = await getGlobalDeploymentHistory();
  const deployments = deploymentRes.success && deploymentRes.data ? deploymentRes.data : [];
  const pagination = deploymentRes.success && deploymentRes.pagination ? deploymentRes.pagination : {
    page: 1,
    limit: 10,
    totalCount: 0
  };

  console.log(deploymentRes);

  return (
    <div className="space-y-8">
      <DeploymentListHeader />
      <DeploymentTable deployments={deployments} pagination={pagination} />

      <footer className="mt-12 mb-6 flex justify-center text-xs text-muted-foreground/50">
        <p>© 2026 Statify. All rights reserved.</p>
      </footer>
    </div>
  );
}
