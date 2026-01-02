import { Rocket, Plus, Search } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export function DeploymentListHeader() {
  return (
    <div className="flex flex-col gap-8">
      <div className="flex flex-col md:flex-row md:items-end justify-between gap-4">
        <div className="flex flex-col gap-2">
          <h1 className="text-white text-3xl md:text-4xl font-black tracking-tight">Deployments</h1>
          <p className="text-muted-foreground text-base">Global view of your static site build history and status.</p>
        </div>
        <Button className="flex items-center justify-center gap-2 bg-primary hover:bg-primary-hover text-black px-5 py-2.5 rounded-lg font-bold transition-colors shadow-lg shadow-primary/25">
          <Plus className="w-5 h-5" />
          New Deployment
        </Button>
      </div>

      <div className="flex flex-col gap-4">
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center border-b border-border gap-4 pb-0">
          <div className="flex gap-1 overflow-x-auto w-full md:w-auto no-scrollbar">
            <div className="relative flex flex-col items-center justify-center px-4 py-3 text-white group cursor-pointer">
              <span className="text-sm font-bold tracking-wide text-primary">All</span>
              <div className="absolute bottom-0 h-0.5 w-full bg-primary rounded-t-full shadow-[0_0_10px_rgba(204,255,0,0.5)]"></div>
            </div>
            <div className="relative flex flex-col items-center justify-center px-4 py-3 text-muted-foreground hover:text-white transition-colors cursor-pointer">
              <span className="text-sm font-semibold tracking-wide">Production</span>
              <div className="absolute bottom-0 h-0.5 w-full bg-transparent group-hover:bg-border transition-colors"></div>
            </div>
            <div className="relative flex flex-col items-center justify-center px-4 py-3 text-muted-foreground hover:text-white transition-colors cursor-pointer">
              <span className="text-sm font-semibold tracking-wide">Preview</span>
              <div className="absolute bottom-0 h-0.5 w-full bg-transparent group-hover:bg-border transition-colors"></div>
            </div>
          </div>
          <div className="w-full md:w-96 pb-2 md:pb-2">
            <div className="relative group">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <Search className="w-4 h-4 text-muted-foreground group-focus-within:text-primary transition-colors" />
              </div>
              <Input
                className="block w-full rounded-lg bg-card border-border text-white placeholder:text-muted-foreground/50 focus:border-primary focus:ring-1 focus:ring-primary sm:text-sm pl-10 py-2.5 shadow-sm transition-all"
                placeholder="Search by commit, hash, or branch..."
                type="text"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
