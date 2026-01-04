"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Bolt, Code, Mail, Lock, Loader2 } from "lucide-react";
import Link from "next/link";
import { login } from "../services/auth.actions";
import { useAction } from "@/hooks/use-action";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { LoginRequestDto } from "../dtos/request/login.request.dto";

export function LoginForm() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const { execute, isPending } = useAction(login, {
    onSuccess: () => {
      router.push("/");
      router.refresh();
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    execute({ email, password });
  };

  return (
    <div className="glass-panel rounded-xl p-8 w-full flex flex-col items-center animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="flex flex-col items-center text-center w-full mb-8">
        <div className="w-12 h-12 rounded-lg bg-primary/10 border border-primary/20 flex items-center justify-center mb-4 text-primary shadow-[0_0_15px_rgba(204,255,0,0.2)]">
          <Bolt className="w-8 h-8" />
        </div>
        <h2 className="text-white text-[22px] font-bold leading-tight tracking-tight">
          Welcome back
        </h2>
        <p className="text-muted-foreground text-sm font-normal mt-2 leading-relaxed">
          Manage your static sites with ease.
        </p>
      </div>

      <div className="w-full mb-6">
        <Button
          variant="outline"
          className="flex w-full items-center justify-center rounded-lg h-12 px-5 bg-secondary border-primary/30 text-white gap-3 text-sm font-bold transition-all duration-300 hover:border-primary hover:shadow-neon group"
        >
          <Code className="w-[22px] h-[22px] text-primary group-hover:scale-110 transition-transform" />
          <span className="group-hover:text-primary transition-colors">
            Continue with GitHub
          </span>
        </Button>
      </div>

      <div className="w-full flex items-center gap-3 mb-6">
        <div className="h-px flex-1 bg-border"></div>
        <p className="text-muted-foreground/70 text-xs font-medium uppercase tracking-wider">
          or continue with email
        </p>
        <div className="h-px flex-1 bg-border"></div>
      </div>

      <form className="w-full flex flex-col gap-5" onSubmit={handleSubmit}>
        <div className="flex flex-col w-full group">
          <Label htmlFor="email" className="text-white text-sm font-medium pb-2">
            Email Address
          </Label>
          <div className="relative">
            <Input
              id="email"
              className="flex w-full rounded-lg text-white border-white/5 bg-input h-11 px-4 pl-10 text-sm font-normal placeholder:text-muted-foreground/30 focus:border-primary focus:ring-1 focus:ring-primary focus:outline-none transition-all shadow-inner"
              placeholder="user@example.com"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              disabled={isPending}
            />
            <Mail className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground/50 w-5 h-5 transition-colors group-focus-within:text-primary" />
          </div>
        </div>

        <div className="flex flex-col w-full group">
          <div className="flex justify-between items-center pb-2">
            <Label htmlFor="password" className="text-white text-sm font-medium">Password</Label>
            <Link
              className="text-primary text-xs hover:text-[#eeffad] transition-colors"
              href="#"
            >
              Forgot password?
            </Link>
          </div>
          <div className="relative">
            <Input
              id="password"
              className="flex w-full rounded-lg text-white border-white/5 bg-input h-11 px-4 pl-10 text-sm font-normal placeholder:text-muted-foreground/30 focus:border-primary focus:ring-1 focus:ring-primary focus:outline-none transition-all shadow-inner"
              placeholder="••••••••"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              disabled={isPending}
            />
            <Lock className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground/50 w-5 h-5 transition-colors group-focus-within:text-primary" />
          </div>
        </div>

        <Button 
          type="submit"
          disabled={isPending}
          className="mt-2 w-full rounded-lg h-11 bg-primary hover:bg-[#d4ff33] text-black text-sm font-bold transition-all shadow-lg shadow-primary/20 hover:shadow-primary/40 active:scale-[0.98]"
        >
          {isPending ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Signing In...
            </>
          ) : (
            "Sign In"
          )}
        </Button>
      </form>

      <div className="mt-6 text-center">
        <p className="text-muted-foreground text-sm">
          Don't have an account?{" "}
          <Link
            className="text-primary font-medium hover:underline hover:text-primary/80"
            href="/register"
          >
            Sign up
          </Link>
        </p>
      </div>
    </div>
  );
}
