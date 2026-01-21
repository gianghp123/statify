'use client'
import Link from "next/link";
import { usePathname } from "next/navigation";

const sidebarItems = [
  {
    title: "Getting Started",
    items: [
      { title: "Introduction", href: "/documentation" },
      { title: "Project Creation", href: "/documentation#project-creation" },
      { title: "Deploying via UI", href: "/documentation#deployment" },
    ],
  },
  {
    title: "Management",
    items: [
      { title: "Analytics", href: "/documentation#analytics" },
      { title: "Project Settings", href: "/documentation#settings" },
    ],
  },
  {
    title: "Resources",
    items: [
      { title: "API Reference", href: "/documentation/api-reference" },
    ],
  },
];

export function DocsSidebar() {
  const pathname = usePathname();

  return (
    <aside className="hidden lg:block w-64 shrink-0 sticky top-16 h-[calc(100vh-4rem)] pt-8 pb-12 overflow-y-auto">
      <nav className="space-y-8">
        {sidebarItems.map((section) => (
          <div key={section.title}>
            <h5 className="mb-3 text-xs font-semibold uppercase tracking-wider text-foreground">
              {section.title}
            </h5>
            <ul className="space-y-2 border-l border-border ml-1">
              {section.items.map((item, index) => {

                return (
                  <li key={index}>
                    <Link
                      href={item.href}
                      className={
                        "block pl-4 -ml-px border-l transition-colors border-transparent hover:border-muted-foreground/50 text-muted-foreground hover:text-foreground"
                      }
                    >
                      {item.title}
                    </Link>
                  </li>
                );
              })}
            </ul>
          </div>
        ))}
      </nav>
    </aside>
  );
}
