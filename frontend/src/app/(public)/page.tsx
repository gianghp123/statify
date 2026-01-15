import { LandingTopbar } from "@/components/landing-topbar";
import { Button } from "@/components/ui/button";
import { getCurrentUser } from "@/features/auth/services/auth.get";
import {
  BarChart3,
  Bolt,
  BookOpen,
  Eye,
  GitBranch,
  Globe,
  Megaphone,
  MoreVertical,
  Rocket,
  User
} from "lucide-react";
import Link from "next/link";
import { Footer } from "react-day-picker";

export default async function LandingPage() {
  const userRes = await getCurrentUser();
  const user = userRes.success ? userRes.data : undefined;

  return (
    <div className="flex flex-col min-h-screen bg-background">
      <LandingTopbar user={user} />

      {/* Hero Section */}
      <section className="pt-32 pb-20 relative overflow-hidden">
        <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-[600px] bg-primary/5 blur-[120px] rounded-full" />
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative">
          <div className="text-center max-w-3xl mx-auto">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-primary/10 border border-primary/20 text-primary text-xs font-bold mb-8 animate-fade-in shadow-[0_0_15px_var(--neon-brand-glow)]">
              <span className="w-2 h-2 rounded-full bg-primary animate-pulse" />
              NOW IN PUBLIC BETA
            </div>
            <h1 className="text-5xl md:text-7xl font-extrabold tracking-tight mb-6 bg-linear-to-b from-foreground to-foreground/60 bg-clip-text text-transparent">
              Static Hosting for the <br />
              <span className="text-primary italic">Modern Web.</span>
            </h1>
            <p className="text-lg md:text-xl text-muted-foreground mb-10 leading-relaxed">
              The fastest way to deploy your static sites. Global edge network,
              automatic CI/CD, and deep analytics built-in.
            </p>
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
              <Button asChild size="lg" className="bg-primary text-primary-foreground hover:bg-primary/90 px-8 py-6 rounded-2xl text-lg font-bold w-full sm:w-auto shadow-neon">
                <Link href="/register">Deploy Now — It's Free</Link>
              </Button>
              <Button asChild variant="outline" size="lg" className="border-border/50 hover:bg-accent/50 px-8 py-6 rounded-2xl text-lg font-bold w-full sm:w-auto text-foreground transition-all">
                <Link href="/documentation">View Documentation</Link>
              </Button>
            </div>

            {/* Dashboard Preview */}
            <div className="mt-20 glass-panel rounded-3xl p-4 border border-border/50 overflow-hidden shadow-2xl dark:shadow-black/50 animate-slide-up">
              <div className="bg-background/40 rounded-2xl border border-border/30 overflow-hidden backdrop-blur-sm">
                <div className="px-6 py-4 border-b border-border/30 flex items-center justify-between bg-muted/20">
                  <div className="flex gap-2">
                    <div className="w-3 h-3 rounded-full bg-destructive/20" />
                    <div className="w-3 h-3 rounded-full bg-warning/20" />
                    <div className="w-3 h-3 rounded-full bg-success/20" />
                  </div>
                  <div className="text-xs text-muted-foreground font-mono">{process.env.NEXT_PUBLIC_PROJECT_URL}/dashboard</div>
                  <div className="w-8" />
                </div>
                <div className="p-6">
                  <div className="flex items-center justify-between mb-8">
                    <h3 className="font-bold text-lg text-foreground">My Projects</h3>
                    <Button variant="outline" size="sm" className="h-8 border-primary/20 text-primary text-xs hover:bg-primary/5">
                      + New Project
                    </Button>
                  </div>
                  <div className="space-y-4">
                    {[1, 2].map((i) => (
                      <div key={i} className="flex items-center justify-between p-4 bg-muted/10 rounded-xl border border-border/20 hover:border-primary/20 transition-all group/item">
                        <div className="flex items-center gap-4">
                          <div className="w-10 h-10 bg-primary/10 rounded-lg flex items-center justify-center text-primary group-hover/item:bg-primary group-hover/item:text-primary-foreground transition-colors">
                            <Globe size={20} />
                          </div>
                          <div className="text-left">
                            <div className="flex items-center gap-2">
                              <span className="font-bold text-foreground">my-awesome-site-{i}</span>
                              <span className="px-2 py-0.5 bg-success/20 text-success text-[10px] font-bold rounded-full uppercase tracking-wider">Live</span>
                            </div>
                            <div className="text-xs text-muted-foreground font-mono">https://site-{i}.{process.env.NEXT_PUBLIC_PROJECT_URL}</div>
                          </div>
                        </div>
                        <div className="flex items-center gap-8 text-xs text-muted-foreground">
                          <div className="hidden sm:flex items-center gap-2">
                            <GitBranch size={14} /> main
                          </div>
                          <div className="hidden sm:block">Jan 12, 2026</div>
                          <MoreVertical size={16} className="cursor-pointer hover:text-foreground transition-colors" />
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Features Grid */}
      <section id="features" className="py-24 relative overflow-hidden">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold mb-4 text-foreground">Everything you need for Static Sites</h2>
            <p className="text-muted-foreground">Built by developers, for developers.</p>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <FeatureCard
              icon={<Rocket />}
              title="One-Click Deploy"
              description="Connect your GitHub and we'll handle the rest. Automatic builds on every push."
            />
            <FeatureCard
              icon={<Globe />}
              title="Global Edge"
              description="Your site is served from 300+ locations worldwide for lightning-fast speeds."
            />
            <FeatureCard
              icon={<Eye />}
              title="Live Previews"
              description="Unique URLs for every pull request to review changes before they go live."
            />
            <FeatureCard
              icon={<BarChart3 />}
              title="Deep Analytics"
              description="Privacy-first analytics to understand your visitors without cookies."
            />
          </div>
        </div>
      </section>

      {/* Workflow Section */}
      <section id="workflow" className="py-24 bg-muted/20 border-y border-border/30">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-16 items-center">
            <div>
              <h2 className="text-3xl md:text-5xl font-bold mb-8 text-foreground">Optimized for your workflow.</h2>
              <div className="space-y-6">
                <WorkflowItem
                  icon={<User />}
                  title="Personal Portfolios"
                  description="Showcase your work on a platform that's as fast as your code."
                />
                <WorkflowItem
                  icon={<BookOpen />}
                  title="Documentation"
                  description="Host your project docs with perfect SEO and instant search."
                />
                <WorkflowItem
                  icon={<Megaphone />}
                  title="Landing Pages"
                  description="Spin up marketing pages in seconds with our high-performance infrastructure."
                />
              </div>
            </div>
            <div className="relative">
              <div className="absolute inset-0 bg-primary/10 blur-3xl rounded-full"></div>
              <div className="relative rounded-2xl shadow-2xl border border-border/30 overflow-hidden group/img transition-all duration-500 hover:shadow-primary/5">
                <img
                  alt="Workflow interface"
                  src="https://images.unsplash.com/photo-1618477247222-acbdb0e159b3?auto=format&fit=crop&q=80&w=1000"
                  className="w-full h-auto grayscale group-hover/img:grayscale-0 transition-all duration-700"
                />
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Final CTA */}
      <section className="py-32 relative">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center glass-panel py-16 rounded-[40px] border border-primary/20 relative overflow-hidden shadow-2xl dark:shadow-black/50">
          <div className="absolute -top-12 left-1/2 -translate-x-1/2 w-24 h-24 bg-primary/10 rounded-full blur-xl" />
          <h2 className="text-4xl md:text-6xl font-extrabold mb-6 text-foreground">From Code to Live <br />in Seconds.</h2>
          <p className="text-lg text-muted-foreground mb-10 max-w-xl mx-auto">
            Join thousands of developers who trust Statify for their hosting needs.
            Start for free, upgrade as you grow.
          </p>
          <Button asChild size="lg" className="bg-primary text-primary-foreground hover:bg-primary/90 px-10 py-7 rounded-2xl text-xl font-bold hover:scale-105 transition-transform shadow-neon">
            <Link href="/register">Deploy Your First Site</Link>
          </Button>
          <p className="mt-6 text-sm text-muted-foreground">No credit card required to get started.</p>
        </div>
      </section>
    </div>
  );
}

function FeatureCard({ icon, title, description }: { icon: React.ReactNode, title: string, description: string }) {
  return (
    <div className="glass-panel p-8 rounded-2xl border border-border/50 hover:border-primary/50 transition-all group cursor-pointer hover:shadow-lg dark:hover:shadow-primary/5">
      <div className="w-12 h-12 bg-primary/10 rounded-xl flex items-center justify-center text-primary mb-6 group-hover:bg-primary group-hover:text-primary-foreground transition-colors shadow-[0_0_15px_var(--neon-brand-soft)]">
        {icon}
      </div>
      <h3 className="text-xl font-bold mb-2 text-foreground">{title}</h3>
      <p className="text-sm text-muted-foreground leading-relaxed">{description}</p>
    </div>
  );
}

function WorkflowItem({ icon, title, description }: { icon: React.ReactNode, title: string, description: string }) {
  return (
    <div className="flex gap-4 group cursor-default">
      <div className="shrink-0 w-12 h-12 bg-primary/10 rounded-full flex items-center justify-center text-primary border border-primary/20 group-hover:bg-primary group-hover:text-primary-foreground transition-all duration-300 shadow-[0_0_15px_var(--neon-brand-soft)]">
        {icon}
      </div>
      <div>
        <h4 className="font-bold text-lg text-foreground">{title}</h4>
        <p className="text-muted-foreground text-sm leading-relaxed">{description}</p>
      </div>
    </div>
  );
}
