import { ArrowRight, Bolt, Globe, Info, Upload, Layout, PieChart } from "lucide-react";
import Link from "next/link";
import { BreadScrum } from "@/components/bread-scrum";
import { DocsTOC } from "@/components/docs/docs-toc";

export default function DocsPage() {
  return (
    <>
      <main className="flex-1 min-w-0 pt-8 pb-16">
        <BreadScrum items={[
          { name: "Docs", href: "/documentation" },
          { name: "Introduction", isCurrent: true }
        ]} />

        <article className="max-w-none">
          <h1 className="text-4xl font-extrabold tracking-tight text-foreground mb-4" id="introduction">
            Introduction
          </h1>
          <p className="text-xl text-muted-foreground leading-relaxed mb-8">
            Welcome to the Statify documentation. Statify is the fastest way to deploy your static websites, providing a global edge network and deep analytics built-in.
          </p>

          <div className="p-6 bg-primary/5 border border-primary/20 rounded-xl mb-8 flex gap-4 not-prose">
            <Info className="w-6 h-6 text-primary shrink-0" />
            <div>
              <p className="text-sm text-primary font-medium m-0">Statify is currently in Public Beta.</p>
              <p className="text-sm text-muted-foreground mt-1 m-0">
                We are currently focusing on UI-based deployments. Git integration and CLI tools are coming soon.
              </p>
            </div>
          </div>

          <h2 className="text-2xl font-bold text-foreground mt-12 mb-4" id="core-principles">
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
              <h4 className="text-lg font-bold m-0">Zero Config Deploys</h4>
              <p className="text-sm text-muted-foreground mt-2">
                Deploy in seconds by simply uploading your build folder as a zip file.
              </p>
            </li>
            <li className="bg-card p-6 rounded-xl border shadow-sm hover:shadow-md transition-shadow">
              <div className="w-10 h-10 bg-primary/10 rounded-lg flex items-center justify-center text-primary mb-4">
                <Globe className="w-6 h-6" />
              </div>
              <h4 className="text-lg font-bold m-0">Global Edge Network</h4>
              <p className="text-sm text-muted-foreground mt-2">
                Your content is served from our worldwide edge network for minimal latency.
              </p>
            </li>
          </ul>

          <h2 className="text-2xl font-bold text-foreground mt-12 mb-4" id="project-creation">
            Project Creation
          </h2>
          <p className="text-muted-foreground">
            To start deploying, you first need to create a project. A project acts as a container for your deployments and configuration.
          </p>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-6 not-prose">
            <div className="p-4 bg-muted/30 rounded-lg border">
              <Layout className="w-5 h-5 mb-2 text-primary" />
              <h4 className="font-bold text-sm">1. Name Your Project</h4>
              <p className="text-xs text-muted-foreground mt-1">Choose a descriptive name for your site.</p>
            </div>
            <div className="p-4 bg-muted/30 rounded-lg border">
              <Globe className="w-5 h-5 mb-2 text-primary" />
              <h4 className="font-bold text-sm">2. Pick a Subdomain</h4>
              <p className="text-xs text-muted-foreground mt-1">Your site will be live at `name.statify.app`.</p>
            </div>
            <div className="p-4 bg-muted/30 rounded-lg border">
              <Upload className="w-5 h-5 mb-2 text-primary" />
              <h4 className="font-bold text-sm">3. Upload Assets</h4>
              <p className="text-xs text-muted-foreground mt-1">Select your production build folder.</p>
            </div>
          </div>

          <h2 className="text-2xl font-bold text-foreground mt-12 mb-4" id="deployment">
            Deploying via UI
          </h2>
          <p className="text-muted-foreground">
            Statify supports manual deployments via the web interface. This is the recommended way for static sites and small applications.
          </p>
          <div className="bg-card p-6 rounded-xl border mt-6 not-prose">
            <h4 className="font-bold mb-3 flex items-center gap-2">
              <Info className="w-4 h-4 text-primary" />
              How to deploy:
            </h4>
            <ol className="space-y-4 text-sm text-muted-foreground">
              <li className="flex gap-3">
                <span className="flex items-center justify-center w-6 h-6 rounded-full bg-primary/20 text-primary text-xs font-bold shrink-0">1</span>
                <span>Run your build script (e.g., `npm run build` or `vite build`) to generate your static files.</span>
              </li>
              <li className="flex gap-3">
                <span className="flex items-center justify-center w-6 h-6 rounded-full bg-primary/20 text-primary text-xs font-bold shrink-0">2</span>
                <span>Compress the contents of your build output folder (usually `dist`, `build`, or `out`) into a **.zip** file.</span>
              </li>
              <li className="flex gap-3">
                <span className="flex items-center justify-center w-6 h-6 rounded-full bg-primary/20 text-primary text-xs font-bold shrink-0">3</span>
                <span>Navigate to your project dashboard and click **"Deploy New Version"**.</span>
              </li>
              <li className="flex gap-3">
                <span className="flex items-center justify-center w-6 h-6 rounded-full bg-primary/20 text-primary text-xs font-bold shrink-0">4</span>
                <span>Drag and drop your zip file and wait for the upload to complete.</span>
              </li>
            </ol>
          </div>

          <h2 className="text-2xl font-bold text-foreground mt-12 mb-4" id="analytics">
            Real-time Analytics
          </h2>
          <p className="text-muted-foreground">
            Monitor your site's performance with built-in analytics. No tracking scripts required—we analyze traffic at the edge.
          </p>
          <div className="mt-6 flex items-center gap-4 p-4 bg-primary/5 border border-primary/20 rounded-xl not-prose">
            <PieChart className="w-8 h-8 text-primary" />
            <p className="text-sm text-muted-foreground m-0">
              View request counts, bandwidth usage, and latency distribution directly in your dashboard.
            </p>
          </div>

          <h2 className="text-2xl font-bold text-foreground mt-12 mb-4" id="settings">
            Project Settings
          </h2>
          <p className="text-muted-foreground">
            Customize your project configuration, manage domain settings, and configure SPA mode in the settings tab.
          </p>

          <div className="mt-12 p-8 bg-linear-to-br from-primary/5 to-transparent border border-primary/20 rounded-2xl not-prose">
            <h3 className="text-xl font-bold mb-2">Ready to go live?</h3>
            <p className="text-muted-foreground mb-6">Create your first project and join the future of static deployment.</p>
            <Link
              href="/dashboard"
              className="inline-flex items-center gap-2 px-6 py-3 bg-primary hover:bg-primary/90 text-primary-foreground font-semibold rounded-lg transition-all shadow-lg shadow-primary/20"
            >
              Go to Dashboard
              <ArrowRight className="w-4 h-4" />
            </Link>
          </div>
        </article>

      </main>

      <DocsTOC items={[
        { id: "introduction", title: "Introduction" },
        { id: "core-principles", title: "Core Principles" },
        { id: "project-creation", title: "Project Creation" },
        { id: "deployment", title: "Deployment via UI" },
        { id: "analytics", title: "Real-time Analytics" },
        { id: "settings", title: "Project Settings" },
      ]} />
    </>
  );
}
