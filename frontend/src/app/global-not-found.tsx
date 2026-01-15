import { ThemeProvider } from "@/components/theme-provider";
import { Button } from "@/components/ui/button";
import { FileQuestion, Home } from "lucide-react";
import type { Metadata } from 'next';
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: 'Statify - 404 Global Not Found',
  description: 'The requested resource could not be found globally.',
}

export default function GlobalNotFound() {
  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          <div className="flex min-h-screen w-full flex-col items-center justify-center bg-background p-8 text-center">
            {/* Pulsing Icon Container */}
            <div className="relative mb-8">
              <div className="absolute inset-0 animate-pulse rounded-3xl bg-primary/10 blur-xl" />
              <div className="relative flex h-24 w-24 items-center justify-center rounded-3xl bg-card border border-border shadow-neon backdrop-blur-xl">
                <FileQuestion className="h-12 w-12 text-primary" />
              </div>
            </div>

            <h1 className="mb-3 text-4xl md:text-5xl font-bold tracking-tight text-foreground">
              404
            </h1>
            <h2 className="mb-4 text-xl md:text-2xl font-semibold text-muted-foreground">
              Resource Not Found
            </h2>
            <p className="mb-10 max-w-md text-base text-muted-foreground/80 leading-relaxed">
              We couldn't locate this resource anywhere in our system. It might have been permanently removed or the link is broken.
            </p>

            <div className="flex flex-col sm:flex-row gap-4">
              <Button asChild className="shadow-neon px-8">
                <a href="/">
                  <Home className="mr-2 h-4 w-4" />
                  Return to Home
                </a>
              </Button>
            </div>

            <div className="mt-16 text-xs text-muted-foreground/40 font-mono tracking-widest uppercase">
              Statify 2026
            </div>
          </div>
        </ThemeProvider>
      </body>
    </html>
  );
}
