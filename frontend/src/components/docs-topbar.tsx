'use client'

import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { logout } from "@/features/auth/services/auth.actions"
import { UserDto } from "@/features/users/dtos/response/user.response.dto"
import { cn } from "@/lib/utils"
import { MessageSquare } from "lucide-react"
import Link from "next/link"
import { Logo } from "./logo"
import { ThemeToggle } from "./theme-toggle"
import { Button } from "./ui/button"

import { usePathname } from "next/navigation"

export function DocsTopbar({ user }: { user: UserDto | undefined }) {
  const pathname = usePathname();

  return <header className="sticky top-0 z-50 w-full border-b border-border/50 bg-background/80 backdrop-blur-md">
    <div className="max-w-[1440px] mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
      <div className="flex items-center gap-8">
        <Logo />
        <nav className="hidden md:flex items-center gap-6 text-sm font-medium text-muted-foreground">
          <Link
            href="/documentation"
            className={cn(
              "transition-colors hover:text-foreground",
              pathname === "/documentation" ? "text-primary text-bold" : "text-muted-foreground"
            )}
          >
            Docs
          </Link>
          <Link
            href="/documentation/api-reference"
            className={cn(
              "transition-colors hover:text-foreground",
              pathname === "/documentation/api-reference" ? "text-primary text-bold" : "text-muted-foreground"
            )}
          >
            API Reference
          </Link>
        </nav>
      </div>

      <div className="flex items-center gap-4">

        {/* Theme Toggle */}
        <div className="w-9 h-9 flex items-center justify-center">
          <ThemeToggle />
        </div>

        {/* Feedback Link */}
        <Link href="#" className="hidden md:flex items-center gap-2 text-sm font-medium text-muted-foreground hover:text-foreground transition-colors">
          <MessageSquare className="w-4 h-4" />
          Feedback
        </Link>

        {/* User Menu or Auth Buttons */}
        {user ? (
          <div className="pl-4 border-l border-border">
            <DropdownMenu modal={false}>
              <DropdownMenuTrigger asChild>
                <button className="flex items-center gap-3 outline-none group">
                  <div className="flex-col items-end hidden sm:flex text-right">
                    <span className="text-sm font-bold leading-none group-hover:text-primary transition-colors">
                      {user.username}
                    </span>
                    <span className="text-xs text-muted-foreground">{user.role}</span>
                  </div>

                  <div className="h-10 w-10 rounded-xl bg-linear-to-br from-primary to-accent flex items-center justify-center text-primary-foreground font-bold shadow-neon border-2 border-background ring-1 ring-border transition-transform active:scale-95">
                    {user.username.charAt(0).toUpperCase()}
                  </div>
                </button>
              </DropdownMenuTrigger>

              <DropdownMenuContent align="end" className="w-56 mt-2">
                <DropdownMenuLabel className="sm:hidden">
                  <div className="flex flex-col">
                    <span className="text-sm font-bold">{user.username}</span>
                    <span className="text-xs text-muted-foreground font-normal">{user.role}</span>
                  </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator className="sm:hidden" />

                <DropdownMenuItem asChild>
                  <Link href="/dashboard">Dashboard</Link>
                </DropdownMenuItem>

                <DropdownMenuItem asChild>
                  <Link href="/settings/account">View Profile</Link>
                </DropdownMenuItem>
                <DropdownMenuItem asChild>
                  <Link href="/settings">Account Settings</Link>
                </DropdownMenuItem>

                <DropdownMenuSeparator />

                <DropdownMenuItem
                  className="text-error focus:text-error cursor-pointer"
                  onClick={logout}
                >
                  Sign Out
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </div>
        ) : (
          <>
            <Link href="/login" className="text-sm font-medium text-foreground hover:text-primary transition-colors">
              Login
            </Link>
            <Button asChild size="sm">
              <Link href="/register" className="bg-primary text-primary-foreground hover:bg-primary/90 rounded-xl font-bold shadow-neon">
                Sign Up
              </Link>
            </Button>
          </>
        )}
      </div>
    </div>
  </header>
}
