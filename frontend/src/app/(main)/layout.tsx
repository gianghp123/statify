import { AppSidebar } from "@/components/app-sidebar";
import { Footer } from "@/components/footer";
import { SidebarProvider } from "@/components/ui/sidebar";
import { logout } from "@/features/auth/services/auth.actions";
import { getCurrentUser } from "@/features/auth/services/auth.get";
import { Bell, Search } from "lucide-react";

export default async function MainLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const userRes = await getCurrentUser()
  if (!userRes.success || !userRes.data) {
    await logout
  }

  const user = userRes.success ? userRes.data : undefined

  return (
    <SidebarProvider defaultOpen={true}>
      <div className="flex min-h-screen bg-background text-foreground w-full">
        <AppSidebar user={user} />

        <div className="flex-1 flex flex-col min-w-0 overflow-hidden">
          <header className="h-16 border-b border-border/50 px-8 flex items-center justify-between bg-background/60 backdrop-blur-md sticky top-0 z-10 transition-all">
            <div className="flex items-center gap-6 flex-1">
              <div className="hidden md:flex items-center gap-2 text-sm text-muted-foreground">
                <span>Teams</span>
                <span className="w-4 h-4 opacity-50">/</span>
                <span className="text-foreground font-semibold">Personal</span>
              </div>
              <div className="relative max-w-md w-full ml-4">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground w-4 h-4" />
                <input
                  className="w-full bg-muted/30 border border-border/50 rounded-lg pl-10 pr-4 py-2.5 text-sm text-foreground focus:outline-none focus:border-primary/50 focus:ring-1 focus:ring-primary/50 transition-all placeholder:text-muted-foreground/60"
                  placeholder="Search projects..."
                  type="text"
                />
              </div>
            </div>
            <div className="flex items-center gap-4">
              {/* <button className="text-muted-foreground hover:text-foreground transition-colors relative p-2 rounded-lg hover:bg-muted/50 border border-transparent hover:border-border/30">
                <Bell className="w-5 h-5" />
                <span className="absolute top-2 right-2 size-2 bg-destructive rounded-full border-2 border-background animate-pulse"></span>
              </button> */}
              <a
                className="text-sm font-bold text-primary hover:text-primary/80 transition-all drop-shadow-[0_0_8px_var(--neon-brand-glow)]"
                href="#"
              >
                Feedback
              </a>
              <a
                className="text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
                href="#"
              >
                Docs
              </a>
            </div>
          </header>

          <main className="flex-1 overflow-y-auto p-8 bg-muted/10 backdrop-blur-[2px]">
            <div className="max-w-6xl mx-auto">
              {children}
            </div>
          </main>
          <Footer />
        </div>
      </div>
    </SidebarProvider>
  );
}
