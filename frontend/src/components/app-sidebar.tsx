"use client";

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  SidebarGroup,
  SidebarGroupContent,
} from "@/components/ui/sidebar";
import {
  LayoutDashboard,
  Folder,
  Rocket,
  BarChart3,
  Settings,
  Bolt,
  ChevronDown,
} from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";

const items = [
  {
    title: "Dashboard",
    url: "/dashboard",
    icon: LayoutDashboard,
  },
  // {
  //   title: "Projects",
  //   url: "/projects",
  //   icon: Folder,
  // },
  {
    title: "Deployments",
    url: "/deployments",
    icon: Rocket,
  },
  {
    title: "Settings",
    url: "/settings",
    icon: Settings,
  },
];

import { ThemeToggle } from "@/components/theme-toggle";

export function AppSidebar() {
  const pathname = usePathname();

  return (
    <Sidebar className="bg-sidebar border-r border-sidebar-border shadow-2xl dark:shadow-black/50">
      <SidebarHeader className="p-6">
        <div className="flex items-center justify-between">
          <Link href="/" className="flex items-center gap-3">
            <div className="size-8 rounded bg-linear-to-br from-primary to-[#aacc00] dark:to-[#aacc00] flex items-center justify-center text-primary-foreground shadow-[0_0_10px_rgba(204,255,0,0.4)]">
              <Bolt className="w-5 h-5 fill-current" />
            </div>
            <h1 className="text-foreground text-xl font-bold tracking-tight">Statify</h1>
          </Link>
          <ThemeToggle />
        </div>
      </SidebarHeader>

      <SidebarContent className="px-4 py-2">
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu className="gap-1">
              {items.map((item) => {
                const isActive = pathname === item.url;
                return (
                  <SidebarMenuItem key={item.title}>
                    <SidebarMenuButton
                      asChild
                      isActive={isActive}
                      className={`group flex items-center gap-3 px-3 py-2.5 rounded-lg transition-all border border-transparent ${
                        isActive
                          ? "bg-sidebar-accent text-foreground border-border/50 shadow-sm"
                          : "text-sidebar-foreground hover:text-foreground hover:bg-sidebar-accent"
                      }`}
                    >
                      <Link href={item.url}>
                        <item.icon
                          className={`w-5 h-5 ${
                            isActive ? "text-primary drop-shadow-[0_0_8px_rgba(204,255,0,0.4)]" : ""
                          }`}
                        />
                        <span className="text-sm font-medium">{item.title}</span>
                        {isActive && (
                          <div className="ml-auto size-2 rounded-full bg-primary shadow-[0_0_8px_rgba(204,255,0,0.6)]" />
                        )}
                      </Link>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                );
              })}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter className="p-4 border-t border-sidebar-border mt-auto">
        <div className="flex items-center gap-3 px-2 py-2 rounded-lg hover:bg-sidebar-accent cursor-pointer transition-colors group">
          <div
            className="size-8 rounded-full bg-cover bg-center bg-primary/20 flex items-center justify-center text-primary text-xs font-bold border border-primary/40"
          >
            AM
          </div>
          <div className="flex flex-col overflow-hidden">
            <p className="text-sm font-medium text-foreground truncate">Alex Morgan</p>
            <p className="text-xs text-sidebar-foreground truncate">alex@statify.app</p>
          </div>
          <ChevronDown className="ml-auto text-sidebar-foreground w-4 h-4" />
        </div>
      </SidebarFooter>
    </Sidebar>
  );
}
