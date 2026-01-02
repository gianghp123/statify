import { SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar"
import { AppSidebar } from "@/components/app-sidebar"
import { Search, Bell } from "lucide-react"

export default function MainLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <SidebarProvider defaultOpen={true}>
      <div className="flex min-h-screen bg-background text-foreground w-full">
        <AppSidebar />
        
        <div className="flex-1 flex flex-col min-w-0 overflow-hidden">
          <header className="h-16 border-b border-white/5 px-8 flex items-center justify-between bg-sidebar/60 backdrop-blur-md sticky top-0 z-10">
            <div className="flex items-center gap-6 flex-1">
              <div className="hidden md:flex items-center gap-2 text-sm text-muted-foreground">
                <span>Teams</span>
                <span className="w-4 h-4 opacity-50">/</span>
                <span className="text-white font-medium">Personal</span>
              </div>
              <div className="relative max-w-md w-full ml-4">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground w-4 h-4" />
                <input
                  className="w-full bg-card/50 border border-white/5 rounded-lg pl-10 pr-4 py-2 text-sm text-white focus:outline-none focus:border-primary/50 focus:ring-1 focus:ring-primary/50 transition-all placeholder:text-muted-foreground/50"
                  placeholder="Search projects..."
                  type="text"
                />
              </div>
            </div>
            <div className="flex items-center gap-4">
              <button className="text-muted-foreground hover:text-white transition-colors relative p-1">
                <Bell className="w-5 h-5" />
                <span className="absolute top-1 right-1 size-2 bg-red-500 rounded-full border-2 border-sidebar"></span>
              </button>
              <a
                className="text-sm font-medium text-primary hover:text-primary/80 transition-colors shadow-[0_0_10px_rgba(204,255,0,0.1)]"
                href="#"
              >
                Feedback
              </a>
              <a
                className="text-sm font-medium text-muted-foreground hover:text-white transition-colors"
                href="#"
              >
                Docs
              </a>
            </div>
          </header>
          
          <main className="flex-1 overflow-y-auto p-8">
            <div className="max-w-6xl mx-auto">
              {children}
            </div>
          </main>
        </div>
      </div>
    </SidebarProvider>
  );
}
