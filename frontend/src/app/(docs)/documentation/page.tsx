import { ArrowRight, Bolt, Copy, ExternalLink, Globe, Info } from "lucide-react";
import Link from "next/link";

export default function DocsPage() {
  return (
    <>
      <main className="flex-1 min-w-0 pt-8 pb-16">
        <nav aria-label="Breadcrumb" className="flex items-center gap-2 text-sm text-muted-foreground mb-6">
          <Link href="/documentation" className="hover:text-foreground transition-colors">
            Docs
          </Link>
          <span className="text-muted-foreground/50">/</span>
          <span className="font-medium text-foreground">Introduction</span>
        </nav>

        <article className="max-w-none">
          <h1 className="text-4xl font-extrabold tracking-tight text-foreground mb-4">
            Introduction
          </h1>
          <p className="text-xl text-muted-foreground leading-relaxed mb-8">
            Welcome to the Statify documentation. Statify is the fastest way to deploy your static websites, providing a global edge network, automatic CI/CD, and deep analytics built-in.
          </p>

          <div className="p-6 bg-primary/5 border border-primary/20 rounded-xl mb-8 flex gap-4 not-prose">
            <Info className="w-6 h-6 text-primary shrink-0" />
            <div>
              <p className="text-sm text-primary font-medium m-0">Statify is currently in Public Beta.</p>
              <p className="text-sm text-muted-foreground mt-1 m-0">
                Some features are still being refined. Join our Discord to share your feedback or report bugs.
              </p>
            </div>
          </div>

          <h2 className="text-2xl font-bold text-foreground mt-12 mb-4">
            Core Principles
          </h2>
          <p className="text-muted-foreground">
            Statify was built for developers who care about speed, reliability, and simplicity. We handle the infrastructure so you can focus on building great products.
          </p>

          <ul className="grid grid-cols-1 md:grid-cols-2 gap-6 list-none pl-0 mt-8 not-prose">
            <li className="bg-card p-6 rounded-xl border shadow-sm hover:shadow-md transition-shadow">
              <div className="w-10 h-10 bg-primary/10 rounded-lg flex items-center justify-center text-primary mb-4">
                <Bolt className="w-6 h-6" />
              </div>
              <h4 className="text-lg font-bold m-0">Blazing Fast Deploys</h4>
              <p className="text-sm text-muted-foreground mt-2">
                Deploy in seconds with our optimized build pipeline and global delivery network.
              </p>
            </li>
            <li className="bg-card p-6 rounded-xl border shadow-sm hover:shadow-md transition-shadow">
              <div className="w-10 h-10 bg-primary/10 rounded-lg flex items-center justify-center text-primary mb-4">
                <Globe className="w-6 h-6" />
              </div>
              <h4 className="text-lg font-bold m-0">Edge First</h4>
              <p className="text-sm text-muted-foreground mt-2">
                Your content is replicated across 300+ edge locations for sub-50ms latency globally.
              </p>
            </li>
          </ul>

          <h2 className="text-2xl font-bold text-foreground mt-12 mb-4" id="installation">
            Installation
          </h2>
          <p className="text-muted-foreground">
            Start by installing the Statify CLI globally on your machine. This tool allows you to deploy projects, manage environment variables, and more.
          </p>

          <div className="relative group mt-6 not-prose">
            <div className="absolute -top-3 left-4 bg-foreground text-background text-[10px] px-2 py-1 rounded font-mono uppercase tracking-widest z-10 font-bold">
              Terminal
            </div>
            <div className="bg-slate-950 dark:bg-black rounded-xl p-5 overflow-x-auto border border-slate-800 shadow-2xl">
              <pre className="flex items-center gap-4 m-0">
                <span className="text-slate-500 select-none">$</span>
                <code className="text-primary font-mono text-sm">npm install -g statify-cli</code>
              </pre>
            </div>
            <button className="absolute top-4 right-4 text-slate-400 hover:text-white transition-colors" title="Copy code">
              <Copy className="w-4 h-4" />
            </button>
          </div>

          <div className="mt-12 p-8 bg-linear-to-br from-primary/5 to-transparent border border-primary/20 rounded-2xl not-prose">
            <h3 className="text-xl font-bold mb-2">Ready to deploy?</h3>
            <p className="text-muted-foreground mb-6">Our quickstart guide will get you live in under 2 minutes.</p>
            <Link
              href="/documentation#quickstart"
              className="inline-flex items-center gap-2 px-6 py-3 bg-primary hover:bg-primary/90 text-primary-foreground font-semibold rounded-lg transition-all shadow-lg shadow-primary/20"
            >
              Get Started
              <ArrowRight className="w-4 h-4" />
            </Link>
          </div>
        </article>

      </main>

      <aside className="hidden xl:block w-64 pt-8 sticky top-16 h-[calc(100vh-4rem)] overflow-y-auto self-start">
        <h5 className="text-sm font-semibold mb-4 uppercase tracking-wider text-foreground">On this page</h5>
        <ul className="space-y-3 text-sm border-l border-border">
          <li>
            <Link href="#" className="block pl-4 -ml-px border-l-2 border-primary text-primary font-medium">
              Introduction
            </Link>
          </li>
          <li>
            <Link href="#" className="block pl-4 -ml-px border-l-2 border-transparent hover:border-muted-foreground/50 text-muted-foreground hover:text-foreground transition-colors">
              Core Principles
            </Link>
          </li>
          <li>
            <Link href="#installation" className="block pl-4 -ml-px border-l-2 border-transparent hover:border-muted-foreground/50 text-muted-foreground hover:text-foreground transition-colors">
              Installation
            </Link>
          </li>
          <li>
            <Link href="#" className="block pl-4 -ml-px border-l-2 border-transparent hover:border-muted-foreground/50 text-muted-foreground hover:text-foreground transition-colors">
              Next Steps
            </Link>
          </li>
        </ul>

        <div className="mt-8 p-4 bg-muted/50 rounded-xl border">
          <p className="text-xs font-semibold mb-2 uppercase">Need help?</p>
          <p className="text-xs text-muted-foreground leading-relaxed mb-4">
            Can't find what you're looking for? Reach out to our engineering team.
          </p>
          <Link href="#" className="text-xs font-medium text-primary hover:underline flex items-center gap-1">
            Contact Support
            <ExternalLink className="w-3 h-3" />
          </Link>
        </div>
      </aside>
    </>
  );
}
