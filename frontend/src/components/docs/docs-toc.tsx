'use client'

import Link from "next/link";
import { cn } from "@/lib/utils";
import { useEffect, useState } from "react";

interface TOCItem {
  id: string;
  title: string;
}

interface DocsTOCProps {
  items: TOCItem[];
}

export function DocsTOC({ items }: DocsTOCProps) {
  const [activeId, setActiveId] = useState<string>("");

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setActiveId(entry.target.id);
          }
        });
      },
      { rootMargin: "-80px 0px -80% 0px" }
    );

    items.forEach((item) => {
      const element = document.getElementById(item.id);
      if (element) observer.observe(element);
    });

    return () => observer.disconnect();
  }, [items]);

  return (
    <aside className="hidden xl:block w-64 pt-8 sticky top-16 h-[calc(100vh-4rem)] overflow-y-auto self-start">
      <h5 className="text-sm font-semibold mb-4 uppercase tracking-wider text-foreground">On this page</h5>
      <ul className="space-y-3 text-sm border-l border-border">
        {items.map((item) => (
          <li key={item.id}>
            <Link
              href={`#${item.id}`}
              className={cn(
                "block pl-4 -ml-px border-l-2 transition-all",
                activeId === item.id
                  ? "border-primary text-primary font-medium"
                  : "border-transparent hover:border-muted-foreground/50 text-muted-foreground hover:text-foreground"
              )}
            >
              {item.title}
            </Link>
          </li>
        ))}
      </ul>
    </aside>
  );
}
