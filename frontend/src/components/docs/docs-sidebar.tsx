'use client'
import Link from "next/link";
import { cn } from "@/lib/utils";

const sidebarItems = [
  {
    title: "Getting Started",
    items: [
      { title: "Introduction", href: "/documentation" },
      { title: "Quickstart Guide", href: "#" },
      { title: "Architecture", href: "#" },
    ],
  },
  {
    title: "Frameworks",
    items: [
      { title: "React & Next.js", href: "#" },
      { title: "Vue & Nuxt", href: "#" },
      { title: "Astro", href: "#" },
      { title: "SvelteKit", href: "#" },
    ],
  },
  {
    title: "Deployments",
    items: [
      { title: "Git Integration", href: "#" },
      { title: "Command Line Interface", href: "#" },
      { title: "Preview Deployments", href: "#" },
      { title: "Environment Variables", href: "#" },
    ],
  },
  {
    title: "Edge Platform",
    items: [
      { title: "Global Edge Network", href: "#" },
      { title: "Edge Functions", href: "#" },
      { title: "Edge Middleware", href: "#" },
    ],
  },
];

export function DocsSidebar() {
  return (
    <aside className="hidden lg:block w-64 shrink-0 sticky top-16 h-[calc(100vh-4rem)] pt-8 pb-12 overflow-y-auto">
      <nav className="space-y-8">
        {sidebarItems.map((section) => (
          <div key={section.title}>
            <h5 className="mb-3 text-xs font-semibold uppercase tracking-wider text-foreground">
              {section.title}
            </h5>
            <ul className="space-y-2 border-l border-border ml-1">
              {section.items.map((item, index) => (
                <li key={index}>
                  <Link
                    href={item.href}
                    className={cn(
                      "block pl-4 -ml-px border-l transition-colors",
                      item.href === "/documentation"
                        ? "border-primary text-primary font-medium"
                        : "border-transparent hover:border-muted-foreground/50 text-muted-foreground hover:text-foreground"
                    )}
                  >
                    {item.title}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        ))}
      </nav>
    </aside>
  );
}
