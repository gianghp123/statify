"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Database, Mail, Lock, ArrowRight, Loader2 } from "lucide-react";
import Link from "next/link";
import { register } from "../services/auth.actions";
import { useAction } from "@/hooks/use-action";
import { useRouter } from "next/navigation";
import { startTransition, useState } from "react";
import { toast } from "sonner";

export function RegisterForm() {
  const router = useRouter();
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const { execute, isPending } = useAction(register, {
    onSuccess: () => {
      toast.success("Registration successful! Please login.");
      router.push("/login");
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    startTransition(() => {
      execute({ username, email, password });
    });
  };

  return (
    <div className="glass-panel w-full max-w-[380px] rounded-2xl p-8 flex flex-col gap-8 relative shadow-2xl dark:shadow-black/50 animate-in fade-in slide-in-from-bottom-6 duration-700">
      <div className="flex justify-center -mb-2">
        <div className="w-14 h-14 rounded-2xl bg-linear-to-br from-primary to-[#a3e635] flex items-center justify-center shadow-xl shadow-primary/20 rotate-3 group-hover:rotate-0 transition-transform">
          <Database className="text-primary-foreground w-8 h-8 font-black" />
        </div>
      </div>
      <div className="text-center flex flex-col gap-2">
        <h1 className="text-foreground text-3xl font-black tracking-tighter">
          Join Statify
        </h1>
        <h2 className="text-muted-foreground text-sm font-bold uppercase tracking-widest text-[10px]">
          The modern platform for static sites
        </h2>
      </div>

      <Button
        variant="outline"
        className="flex w-full items-center justify-center rounded-xl h-12 px-4 bg-background border-border text-foreground hover:bg-muted/50 hover:border-primary transition-all duration-300 gap-3 font-bold text-sm shadow-sm group"
      >
        <svg
          className="transition-transform group-hover:scale-110 duration-200"
          fill="currentColor"
          height="20"
          viewBox="0 0 16 16"
          width="20"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.012 8.012 0 0 0 16 8c0-4.42-3.58-8-8-8z"></path>
        </svg>
        <span className="group-hover:text-primary transition-colors">Continue with GitHub</span>
      </Button>

      <div className="relative flex py-1 items-center">
        <div className="grow border-t border-border"></div>
        <span className="shrink-0 mx-4 text-muted-foreground text-[10px] font-black uppercase tracking-[0.2em]">
          Or email
        </span>
        <div className="grow border-t border-border"></div>
      </div>

      <form className="flex flex-col gap-5" onSubmit={handleSubmit}>
        <div className="flex flex-col gap-2 group">
          <Label htmlFor="username" className="text-foreground text-sm font-bold pl-1 uppercase tracking-wider text-[11px]">
            Username
          </Label>
          <div className="relative">
            <Input
              id="username"
              className="w-full rounded-xl bg-background border-border focus:border-primary focus:ring-1 focus:ring-primary text-foreground placeholder:text-muted-foreground/30 h-12 px-4 text-base font-medium transition-all duration-300 shadow-sm"
              placeholder="johndoe"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              disabled={isPending}
              required
            />
          </div>
        </div>
        <div className="flex flex-col gap-2 group">
          <Label className="text-foreground text-sm font-bold pl-1 uppercase tracking-wider text-[11px]">
            Email Address
          </Label>
          <div className="relative">
            <span className="absolute inset-y-0 left-0 flex items-center pl-4 text-muted-foreground/50 group-focus-within:text-primary transition-colors duration-300">
              <Mail className="w-5 h-5" />
            </span>
            <Input
              className="w-full rounded-xl bg-background border-border focus:border-primary focus:ring-1 focus:ring-primary text-foreground placeholder:text-muted-foreground/30 h-12 pl-12 pr-4 text-base font-medium transition-all duration-300 shadow-sm"
              placeholder="name@example.com"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              disabled={isPending}
              required
            />
          </div>
        </div>
        <div className="flex flex-col gap-2 group">
          <Label className="text-foreground text-sm font-bold pl-1 uppercase tracking-wider text-[11px]">
            Password
          </Label>
          <div className="relative">
            <span className="absolute inset-y-0 left-0 flex items-center pl-4 text-muted-foreground/50 group-focus-within:text-primary transition-colors duration-300">
              <Lock className="w-5 h-5" />
            </span>
            <Input
              className="w-full rounded-xl bg-background border-border focus:border-primary focus:ring-1 focus:ring-primary text-foreground placeholder:text-muted-foreground/30 h-12 pl-12 pr-4 text-base font-medium transition-all duration-300 shadow-sm"
              placeholder="••••••••"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              disabled={isPending}
              required
            />
          </div>
        </div>
        <Button
          className="mt-4 w-full h-12 rounded-xl bg-primary text-primary-foreground text-sm font-black tracking-widest uppercase hover:bg-primary/90 active:scale-[0.98] transition-all duration-200 flex items-center justify-center gap-2 shadow-neon-brand hover:shadow-neon"
          type="submit"
          disabled={isPending}
        >
          {isPending ? (
            <Loader2 className="w-5 h-5 animate-spin" />
          ) : (
            <>
              <span>Get Started</span>
              <ArrowRight className="w-5 h-5" />
            </>
          )}
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
