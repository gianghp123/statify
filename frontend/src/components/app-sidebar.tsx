"use client";

import { useState } from "react"; // Added for modal state
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import {
  Bolt,
  ChevronDown,
  LayoutDashboard,
  Rocket,
  Settings,
  User,
  LogOut,
} from "lucide-react";
import Link from "next/link";
import { usePathname } from "next/navigation";

// --- Shadcn UI Imports ---
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";

import { ThemeToggle } from "@/components/theme-toggle";
import { UserDto } from "@/features/users/dtos/response/user.response.dto";
import { logout } from "@/features/auth/services/auth.actions";

const items = [
  { title: "Dashboard", url: "/dashboard", icon: LayoutDashboard },
  { title: "Deployments", url: "/deployments", icon: Rocket },
  { title: "Settings", url: "/settings", icon: Settings },
];

export function AppSidebar({ user }: { user: UserDto | undefined }) {
  const pathname = usePathname();
  const [showLogoutDialog, setShowLogoutDialog] = useState(false);

  // Get initials for the avatar (e.g., "John Doe" -> "JD")
  const initials = user?.username
    ? user.username.split(" ").map((n) => n[0]).join("").toUpperCase().slice(0, 2)
    : "??";

  return (
    <>
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
                        className={`group flex items-center gap-3 px-3 py-2.5 rounded-lg transition-all border border-transparent ${isActive
                          ? "bg-sidebar-accent text-foreground border-border/50 shadow-sm"
                          : "text-sidebar-foreground hover:text-foreground hover:bg-sidebar-accent"
                          }`}
                      >
                        <Link href={item.url}>
                          <item.icon
                            className={`w-5 h-5 ${isActive ? "text-primary drop-shadow-[0_0_8px_rgba(204,255,0,0.4)]" : ""
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
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <div className="flex items-center gap-3 px-2 py-2 rounded-lg hover:bg-sidebar-accent cursor-pointer transition-colors group outline-none">
                <div className="size-8 rounded-full bg-primary/20 flex items-center justify-center text-primary text-xs font-bold border border-primary/40 shrink-0">
                  {initials}
                </div>
                <div className="flex flex-col overflow-hidden text-left">
                  <p className="text-sm font-medium text-foreground truncate">
                    {user?.username || "Guest"}
                  </p>
                  <p className="text-xs text-sidebar-foreground truncate">
                    {user?.email || "No email"}
                  </p>
                </div>
                <ChevronDown className="ml-auto text-sidebar-foreground w-4 h-4 shrink-0" />
              </div>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" side="top" className="w-56 mb-2">
              <DropdownMenuLabel>My Account</DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem asChild>
                <Link href="/settings/account" className="flex items-center cursor-pointer">
                  <User className="mr-2 h-4 w-4" />
                  <span>Profile</span>
                </Link>
              </DropdownMenuItem>
              <DropdownMenuItem asChild>
                <Link href="/settings" className="flex items-center cursor-pointer">
                  <Settings className="mr-2 h-4 w-4" />
                  <span>Settings</span>
                </Link>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                onSelect={() => setShowLogoutDialog(true)}
                className="text-destructive focus:text-destructive cursor-pointer"
              >
                <LogOut className="mr-2 h-4 w-4" />
                <span>Log out</span>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </SidebarFooter>
      </Sidebar>

      {/* Logout Dialog */}
      <AlertDialog open={showLogoutDialog} onOpenChange={setShowLogoutDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This will end your current session. You will need to log back in to access your dashboard.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={logout}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              Log out
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}