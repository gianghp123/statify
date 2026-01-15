
import Link from "next/link";
import { Bolt } from "lucide-react";

export function Footer() {
  return <footer className="py-12 border-t border-border/50 bg-background/50">
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div className="flex flex-col md:flex-row justify-between items-center gap-8">
        <div className="flex items-center gap-2">
          <div className="w-6 h-6 bg-primary rounded flex items-center justify-center shadow-[0_0_8px_var(--neon-brand-glow)]">
            <Bolt size={14} className="text-primary-foreground fill-current" />
          </div>
          <span className="font-bold text-foreground">Statify</span>
        </div>
        <div className="flex gap-8 text-sm text-muted-foreground">
          <Link className="hover:text-primary transition-colors" href="#">Privacy</Link>
          <Link className="hover:text-primary transition-colors" href="#">Terms</Link>
          <Link className="hover:text-primary transition-colors" href="#">Status</Link>
          <Link className="hover:text-primary transition-colors" href="#">Contact</Link>
        </div>
        <div className="text-sm text-muted-foreground">
          © {new Date().getFullYear()} Statify Inc.
        </div>
      </div>
    </div>
  </footer>
}