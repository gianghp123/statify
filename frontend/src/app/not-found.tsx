import Link from "next/link";
import { Button } from "@/components/ui/button";
import { MoveLeft, AlertCircle } from "lucide-react";

export default function NotFound() {
  return (
    <div className="flex h-screen w-full flex-col items-center justify-center bg-background p-4 text-center">
      <div className="mx-auto flex h-full w-full max-w-md flex-col items-center justify-center gap-6">
        {/* Animated Icon Container */}
        <div className="relative flex items-center justify-center">
          <div className="absolute inset-0 animate-pulse rounded-full bg-destructive/20 blur-xl" />
          <div className="relative flex h-24 w-24 items-center justify-center rounded-full bg-background/50 ring-1 ring-border shadow-neon backdrop-blur-xl transition-transform hover:scale-105">
            <AlertCircle className="h-12 w-12 text-destructive animate-bounce-slow" />
          </div>
        </div>

        <div className="space-y-2">
          <h1 className="text-4xl font-bold tracking-tighter sm:text-5xl text-foreground">
            404 Not Found
          </h1>
          <p className="text-lg text-muted-foreground">
            We couldn't find the page you were looking for. It might have been removed, renamed, or doesn't exist.
          </p>
        </div>

        <div className="flex w-full flex-col gap-2 sm:flex-row sm:justify-center">
          <Button asChild size="lg" className="shadow-neon group gap-2">
            <Link href="/">
              <MoveLeft className="h-4 w-4 transition-transform group-hover:-translate-x-1" />
              Back to Home
            </Link>
          </Button>
          <Button asChild variant="outline" size="lg" className="gap-2">
            <Link href="/dashboard">
              Go to Dashboard
            </Link>
          </Button>
        </div>

        <div className="mt-8 text-xs text-muted-foreground/50">
          Error Code: 404 • Statify
        </div>
      </div>
    </div>
  );
}
