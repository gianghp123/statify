"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Database, Mail, Lock, ArrowRight } from "lucide-react";
import Link from "next/link";

export function RegisterForm() {
  return (
    <div className="glass-panel w-full max-w-[380px] rounded-xl p-8 flex flex-col gap-6 relative shadow-2xl shadow-black/50 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="flex justify-center mb-2">
        <div className="w-12 h-12 rounded-lg bg-gradient-to-br from-primary to-[#a3e635] flex items-center justify-center shadow-lg shadow-primary/20">
          <Database className="text-primary-foreground w-7 h-7 font-bold" />
        </div>
      </div>
      <div className="text-center flex flex-col gap-2">
        <h1 className="text-white text-2xl font-bold tracking-tight">
          Create your account
        </h1>
        <h2 className="text-muted-foreground text-sm font-medium">
          Start hosting your static sites in seconds.
        </h2>
      </div>

      <Button
        variant="outline"
        className="flex w-full items-center justify-center rounded-lg h-11 px-4 bg-transparent border-primary text-primary hover:bg-primary hover:text-primary-foreground transition-all duration-300 gap-3 font-semibold text-sm shadow-neon group"
      >
        <svg
          className="transition-transform group-hover:scale-110 duration-300"
          fill="currentColor"
          height="20"
          viewBox="0 0 16 16"
          width="20"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.012 8.012 0 0 0 16 8c0-4.42-3.58-8-8-8z"></path>
        </svg>
        <span>Continue with GitHub</span>
      </Button>

      <div className="relative flex py-1 items-center">
        <div className="flex-grow border-t border-white/10"></div>
        <span className="flex-shrink-0 mx-4 text-muted-foreground text-xs font-medium uppercase tracking-wider">
          Or register with email
        </span>
        <div className="flex-grow border-t border-white/10"></div>
      </div>

      <form className="flex flex-col gap-4" onSubmit={(e) => e.preventDefault()}>
        <div className="flex flex-col gap-1.5 group">
          <Label className="text-white text-sm font-medium pl-1">
            Email Address
          </Label>
          <div className="relative">
            <span className="absolute inset-y-0 left-0 flex items-center pl-3 text-muted-foreground group-focus-within:text-primary transition-colors duration-300">
              <Mail className="w-5 h-5" />
            </span>
            <Input
              className="w-full rounded-lg bg-input border-transparent focus:border-primary focus:ring-1 focus:ring-primary text-white placeholder-white/20 h-11 pl-10 pr-4 text-sm transition-all duration-300 shadow-sm"
              placeholder="name@example.com"
              type="email"
            />
          </div>
        </div>
        <div className="flex flex-col gap-1.5 group">
          <Label className="text-white text-sm font-medium pl-1">
            Password
          </Label>
          <div className="relative">
            <span className="absolute inset-y-0 left-0 flex items-center pl-3 text-muted-foreground group-focus-within:text-primary transition-colors duration-300">
              <Lock className="w-5 h-5" />
            </span>
            <Input
              className="w-full rounded-lg bg-input border-transparent focus:border-primary focus:ring-1 focus:ring-primary text-white placeholder-white/20 h-11 pl-10 pr-4 text-sm transition-all duration-300 shadow-sm"
              placeholder="••••••••"
              type="password"
            />
          </div>
        </div>
        <Button
          className="mt-2 w-full h-11 rounded-lg bg-primary text-primary-foreground text-sm font-bold tracking-wide hover:bg-[#d4ff33] active:scale-[0.98] transition-all duration-200 flex items-center justify-center gap-2 shadow-[0_0_15px_rgba(204,255,0,0.3)]"
          type="button"
        >
          <span>Get Started</span>
          <ArrowRight className="w-[18px] h-[18px]" />
        </Button>
      </form>

      <div className="text-center pt-2">
        <p className="text-muted-foreground text-sm">
          Already have an account?{" "}
          <Link
            className="text-primary hover:text-primary/80 font-medium transition-colors ml-1 hover:underline decoration-primary/50 underline-offset-4"
            href="/login"
          >
            Log in
          </Link>
        </p>
      </div>
    </div>
  );
}
