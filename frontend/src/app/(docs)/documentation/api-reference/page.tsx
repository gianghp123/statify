import Link from "next/link";
import { Copy, Lightbulb } from "lucide-react";
import { cn } from "@/lib/utils";

export default function ApiReferencePage() {
  return (
    <div className="flex flex-col xl:flex-row gap-10">
      <div className="flex-1 min-w-0">
        <div className="mb-6 flex items-center gap-2 text-sm text-muted-foreground">
          <Link href="/documentation" className="hover:text-foreground transition-colors">
            Docs
          </Link>
          <span className="text-muted-foreground/50">/</span>
          <span className="font-medium text-foreground">API Reference</span>
        </div>

        <article className="prose prose-slate dark:prose-invert max-w-none">
          <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl mb-4">
            Create Deployment
          </h1>
          <p className="leading-7 not-first:mt-6 text-xl text-muted-foreground mb-8">
            Create a new deployment for a specific project. This endpoint handles file uploads and build triggers.
          </p>

          <section className="mb-10 not-prose">
            <h2 className="text-sm font-bold uppercase tracking-widest text-muted-foreground mb-4 border-b pb-2">
              Endpoint
            </h2>
            <div className="flex items-center gap-2 font-mono text-sm bg-muted/50 p-3 rounded-lg border border-border">
              <span className="text-muted-foreground">https://api.statify.online/v1</span>
              <span className="text-primary font-bold">/deployments</span>
            </div>
          </section>

          <section className="mb-10 not-prose">
            <h2 className="text-sm font-bold uppercase tracking-widest text-muted-foreground mb-4 border-b pb-2">
              Body Parameters
            </h2>
            <div className="space-y-6">
              <div className="flex flex-col gap-1">
                <div className="flex items-center gap-2">
                  <span className="font-mono text-sm font-semibold">name</span>
                  <span className="text-xs px-1.5 py-0.5 bg-muted rounded text-muted-foreground uppercase">string</span>
                  <span className="text-xs text-orange-500 font-medium">Required</span>
                </div>
                <p className="text-sm text-muted-foreground">
                  A custom name for this deployment. Typically used for internal tracking.
                </p>
              </div>
              <div className="flex flex-col gap-1">
                <div className="flex items-center gap-2">
                  <span className="font-mono text-sm font-semibold">project_id</span>
                  <span className="text-xs px-1.5 py-0.5 bg-muted rounded text-muted-foreground uppercase">string</span>
                  <span className="text-xs text-orange-500 font-medium">Required</span>
                </div>
                <p className="text-sm text-muted-foreground">
                  The unique identifier of the project you want to deploy to.
                </p>
              </div>
              <div className="flex flex-col gap-1">
                <div className="flex items-center gap-2">
                  <span className="font-mono text-sm font-semibold">files</span>
                  <span className="text-xs px-1.5 py-0.5 bg-muted rounded text-muted-foreground uppercase">object[]</span>
                </div>
                <p className="text-sm text-muted-foreground">
                  An array of file objects mapping paths to content.
                </p>
              </div>
            </div>
          </section>

          <section className="mb-10 not-prose">
            <h2 className="text-sm font-bold uppercase tracking-widest text-muted-foreground mb-4 border-b pb-2">
              Responses
            </h2>
            <div className="bg-card rounded-xl border overflow-hidden shadow-sm">
              <table className="w-full text-left text-sm">
                <thead className="bg-muted/50 border-b">
                  <tr>
                    <th className="px-6 py-3 font-semibold text-muted-foreground">Code</th>
                    <th className="px-6 py-3 font-semibold text-muted-foreground">Description</th>
                  </tr>
                </thead>
                <tbody className="divide-y">
                  <tr>
                    <td className="px-6 py-4 font-mono font-medium text-emerald-500">201 Created</td>
                    <td className="px-6 py-4 text-muted-foreground">
                      Deployment initiated successfully. Response includes the deployment ID.
                    </td>
                  </tr>
                  <tr>
                    <td className="px-6 py-4 font-mono font-medium text-amber-500">400 Bad Request</td>
                    <td className="px-6 py-4 text-muted-foreground">
                      The request body is malformed or missing required fields.
                    </td>
                  </tr>
                  <tr>
                    <td className="px-6 py-4 font-mono font-medium text-red-500">401 Unauthorized</td>
                    <td className="px-6 py-4 text-muted-foreground">
                      Invalid or missing API token in the Authorization header.
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </section>
        </article>
      </div>

      <div className="lg:col-span-5 space-y-6 xl:w-[450px]">
        <div className="sticky top-24">
          <div className="bg-slate-950 dark:bg-black rounded-2xl overflow-hidden shadow-2xl border border-slate-800">
            <div className="flex items-center justify-between px-4 py-3 bg-slate-900 border-b border-slate-800">
              <div className="flex items-center gap-2">
                <div className="flex gap-1.5">
                  <div className="w-2.5 h-2.5 rounded-full bg-slate-700"></div>
                  <div className="w-2.5 h-2.5 rounded-full bg-slate-700"></div>
                  <div className="w-2.5 h-2.5 rounded-full bg-slate-700"></div>
                </div>
                <span className="text-xs font-medium text-slate-400 ml-2">curl --request POST</span>
              </div>
              <button className="flex items-center gap-1 text-[10px] font-bold text-slate-400 hover:text-white transition-colors">
                <Copy className="w-3 h-3" />
                COPY
              </button>
            </div>
            <div className="p-6 overflow-x-auto">
              <pre className="font-mono text-[13px] leading-relaxed text-slate-300">
                <code>
                  <span className="text-primary">curl</span> -X POST https://api.statify.online/v1/deployments \{"\n"}
                  {"  "}-H <span className="text-emerald-400">"Authorization: Bearer <span className="text-slate-500">ST_TOKEN_XXXX"</span></span> \{"\n"}
                  {"  "}-H <span className="text-emerald-400">"Content-Type: application/json"</span> \{"\n"}
                  {"  "}-d '{"{"}\n
                  {"    "}<span className="text-primary">"name"</span>: <span className="text-emerald-400">"production-launch"</span>,\n
                  {"    "}<span className="text-primary">"project_id"</span>: <span className="text-emerald-400">"prj_8f2k9l1a"</span>,\n
                  {"    "}<span className="text-primary">"files"</span>: [\n
                  {"      "}{"{"}\n
                  {"        "}<span className="text-primary">"file"</span>: <span className="text-emerald-400">"index.html"</span>,\n
                  {"        "}<span className="text-primary">"data"</span>: <span className="text-emerald-400">"&lt;html&gt;Hello World&lt;/html&gt;"</span>\n
                  {"      "}{"}"}\n
                  {"    "}]\n
                  {"  "}'{'}'}
                </code>
              </pre>
            </div>
          </div>

          <div className="mt-6 bg-slate-950 dark:bg-black rounded-2xl overflow-hidden shadow-2xl border border-slate-800">
            <div className="px-4 py-3 bg-slate-900 border-b border-slate-800 flex justify-between items-center">
              <span className="text-xs font-bold text-emerald-400 uppercase tracking-widest">Example Response</span>
              <span className="text-[10px] text-slate-500 font-mono">201 CREATED</span>
            </div>
            <div className="p-6">
              <pre className="font-mono text-[13px] leading-relaxed text-slate-300">
                <code>
                  {"{"}\n
                  {"  "}<span className="text-primary">"id"</span>: <span className="text-emerald-400">"dep_v7k2m9p1"</span>,\n
                  {"  "}<span className="text-primary">"status"</span>: <span className="text-emerald-400">"initializing"</span>,\n
                  {"  "}<span className="text-primary">"url"</span>: <span className="text-emerald-400">"https://vite.statify.online"</span>,\n
                  {"  "}<span className="text-primary">"created_at"</span>: <span className="text-emerald-400">"2024-03-24T14:22:01Z"</span>\n
                  {"}"}
                </code>
              </pre>
            </div>
          </div>

          <div className="mt-8 bg-primary/5 border border-primary/20 rounded-2xl p-5">
            <div className="flex items-center gap-3 mb-2 text-primary">
              <Lightbulb className="w-5 h-5" />
              <h4 className="font-bold text-sm">Pro Tip</h4>
            </div>
            <p className="text-sm text-muted-foreground leading-relaxed">
              Use the Statify CLI for large deployments. It handles file hashing and concurrent uploads automatically, providing a 10x faster experience for sites with &gt;1000 assets.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
