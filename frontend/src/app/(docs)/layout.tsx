import { DocsSidebar } from "@/components/docs/docs-sidebar";
import { DocsTopbar } from "@/components/docs-topbar";
import { getCurrentUser } from "@/features/auth/services/auth.get";
import { Footer } from "@/components/footer";

export default async function DocsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const userRes = await getCurrentUser();
  const user = userRes.success ? userRes.data : undefined;

  return (
    <>
      <DocsTopbar user={user} />
      <div className="max-w-[1440px] mx-auto px-4 sm:px-6 lg:px-8 flex gap-8">
        <DocsSidebar />
        {children}
      </div>
      <Footer />
    </>
  );
}
