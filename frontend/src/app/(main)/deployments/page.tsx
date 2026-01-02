import { DeploymentListHeader } from "@/features/deployments/components/DeploymentListHeader";
import { DeploymentTable } from "@/features/deployments/components/DeploymentTable";

export default function DeploymentsPage() {
  return (
    <div className="space-y-8">
      <DeploymentListHeader />
      <DeploymentTable />
      
      <footer className="mt-12 mb-6 flex justify-center text-xs text-muted-foreground/50">
        <p>© 2024 Statify Inc. Custom Deep Plum Edition.</p>
      </footer>
    </div>
  );
}
