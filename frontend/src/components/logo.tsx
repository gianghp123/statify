import { Bolt } from "lucide-react";

export function Logo() {
  return <div className="flex items-center gap-2">
    <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center shadow-[0_0_10px_var(--neon-brand-glow)]">
      <Bolt className="text-primary-foreground w-5 h-5 fill-current" />
    </div>
    <span className="text-xl font-bold tracking-tight text-foreground">Statify</span>
  </div>
}