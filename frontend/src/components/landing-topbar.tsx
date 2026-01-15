'use client'

import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu"
import { UserDto } from "@/features/users/dtos/response/user.response.dto"
import { Bolt } from "lucide-react"
import Link from "next/link"
import { ThemeToggle } from "./theme-toggle"
import { Button } from "./ui/button"
import { logout } from "@/features/auth/services/auth.actions"
import { Logo } from "./logo"


export function LandingTopbar({ user }: { user: UserDto | undefined }) {
  {/* Navigation */ }
  return <nav className="fixed top-0 w-full z-50 border-b border-border/50 bg-background/80 backdrop-blur-md">
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div className="flex justify-between items-center h-16">
        <Logo />
        <div className="hidden md:flex items-center gap-8">
          <Link href="#features" className="text-sm text-muted-foreground hover:text-primary transition-colors">Features</Link>
          <Link href="/documentation" className="text-sm text-muted-foreground hover:text-primary transition-colors">Docs</Link>
          <Link href="#workflow" className="text-sm text-muted-foreground hover:text-primary transition-colors">Workflow</Link>
          <Link href="#pricing" className="text-sm text-muted-foreground hover:text-primary transition-colors">Pricing</Link>
        </div>

        <div className="flex items-center gap-4">

          {user ? (
            <div className="flex items-center gap-4">
              <div className="pl-4">
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

                    <DropdownMenuItem asChild><Link
                      href="/dashboard"
                      className="hidden md:flex items-center gap-2 text-sm font-medium hover:text-primary transition-colors"
                    >
                      Dashboard
                    </Link>
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
            </div>
          ) : (
            <>
              <Link href="/login" className="text-sm font-medium text-foreground hover:text-primary transition-colors">
                Login
              </Link>
              <Button asChild>
                <Link href="/register" className="bg-primary text-primary-foreground hover:bg-primary/90 rounded-xl font-bold shadow-neon">
                  Sign Up Now
                </Link>
              </Button>
            </>
          )}
          <div className="w-9 h-9 flex items-center justify-center">
            <ThemeToggle />
          </div>
        </div>
      </div>
    </div>
  </nav>
} 