"use client";

import { ChevronLeft } from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";

export default function SettingsLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const isRoot = pathname === "/settings";

  return (
    <div className="container mx-auto py-6 space-y-6">
      {!isRoot && (
        <Link
          href="/settings"
          className="inline-flex items-center text-sm text-muted-foreground hover:text-primary transition-colors group"
        >
          <ChevronLeft className="w-4 h-4 mr-1 group-hover:-translate-x-1 transition-transform" />
          Back to Settings
        </Link>
      )}
      {children}
    </div>
  );
}
