import { ChevronRight } from "lucide-react"
import Link from "next/link"
import React from "react";

interface BreadcrumItem {
  href?: string;
  name: string;
  isCurrent?: boolean;
  icon?: React.ComponentType<{ className?: string }>;
}

export function BreadScrum({ items }: { items: BreadcrumItem[] }) {

  return <nav className="flex items-center gap-2 text-sm mb-6 border-b border-border pb-4">
    {items.map((item, index) => (
      <React.Fragment key={index}>
        {item.icon && <item.icon className="w-4 h-4" />}
        {item.isCurrent ? (
          <span className="text-foreground font-black px-2 py-0.5 bg-muted/50 rounded border border-border/50">{item.name}</span>
        ) : (
          <Link
            href={item.href || "#"}
            className="text-muted-foreground hover:text-primary transition-all flex items-center gap-1.5 font-bold"
          >
            {item.name}
          </Link>
        )}
        {!item.isCurrent && (
          <ChevronRight className="w-4 h-4 text-muted-foreground/30" />
        )}
      </React.Fragment>
    ))}
  </nav>
}
