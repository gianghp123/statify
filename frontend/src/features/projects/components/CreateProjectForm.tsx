"use client";

import { useState } from "react";
import { Upload, X, FileArchive, Globe, Info, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useRouter } from "next/navigation";

export function CreateProjectForm() {
  const router = useRouter();
  const [projectName, setProjectName] = useState("");
  const [file, setFile] = useState<File | null>(null);
  const [isUploading, setIsUploading] = useState(false);

  const subdomain = projectName
    .toLowerCase()
    .replace(/[^a-z0-9]/g, "-")
    .replace(/-+/g, "-")
    .replace(/^-|-$/g, "") || "your-project";

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const selectedFile = e.target.files[0];
      if (selectedFile.name.endsWith(".zip")) {
        setFile(selectedFile);
      } else {
        alert("Please upload a .zip file");
      }
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!projectName || !file) return;

    setIsUploading(true);
    // Simulate upload delay
    setTimeout(() => {
      setIsUploading(false);
      router.push("/");
    }, 2000);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-8">
      <div className="space-y-6">
        {/* Project Name */}
        <div className="space-y-2">
          <Label htmlFor="name" className="text-white font-semibold">Project Name</Label>
          <Input
            id="name"
            placeholder="my-awesome-site"
            value={projectName}
            onChange={(e) => setProjectName(e.target.value)}
            className="bg-card border-border text-white focus:ring-primary h-12 text-lg"
            required
          />
          <p className="text-sm text-muted-foreground flex items-center gap-2">
            <Globe className="w-4 h-4 text-primary" />
            Your site will be available at: 
            <span className="text-white font-mono bg-white/5 px-2 py-0.5 rounded border border-white/10">
              {subdomain}.statify.app
            </span>
          </p>
        </div>

        {/* File Upload */}
        <div className="space-y-2">
          <Label className="text-white font-semibold">Deployment Source</Label>
          <div 
            className={`relative border-2 border-dashed rounded-xl p-12 transition-all flex flex-col items-center justify-center gap-4 ${
              file ? 'border-primary/50 bg-primary/5' : 'border-border hover:border-primary/30 bg-card/50'
            }`}
          >
            <input
              type="file"
              accept=".zip"
              onChange={handleFileChange}
              className="absolute inset-0 opacity-0 cursor-pointer"
              required={!file}
            />
            
            {file ? (
              <>
                <div className="size-16 rounded-2xl bg-primary/20 flex items-center justify-center text-primary shadow-neon">
                  <FileArchive className="w-8 h-8" />
                </div>
                <div className="text-center">
                  <p className="text-white font-bold">{file.name}</p>
                  <p className="text-xs text-muted-foreground">{(file.size / 1024 / 1024).toFixed(2)} MB</p>
                </div>
                <Button 
                  type="button" 
                  variant="ghost" 
                  size="sm" 
                  onClick={() => setFile(null)}
                  className="text-red-400 hover:text-red-300 hover:bg-red-400/10"
                >
                  <X className="w-4 h-4 mr-2" />
                  Remove file
                </Button>
              </>
            ) : (
              <>
                <div className="size-16 rounded-2xl bg-white/5 border border-white/10 flex items-center justify-center text-muted-foreground">
                  <Upload className="w-8 h-8" />
                </div>
                <div className="text-center">
                  <p className="text-white font-bold text-lg">Click or drag to upload</p>
                  <p className="text-sm text-muted-foreground mt-1">
                    Upload a ZIP folder containing your build files
                  </p>
                </div>
              </>
            )}
          </div>
          <div className="flex items-start gap-2 p-4 rounded-lg bg-indigo-500/10 border border-indigo-500/20 text-xs text-indigo-300 leading-relaxed">
            <Info className="w-4 h-4 shrink-0 mt-0.5" />
            <p>
              Make sure your ZIP archive contains an <span className="text-white font-mono">index.html</span> file 
              at the root level or your project's build output directory.
            </p>
          </div>
        </div>
      </div>

      <div className="flex items-center gap-4 pt-4">
        <Button 
          type="button" 
          variant="ghost" 
          onClick={() => router.back()}
          className="text-muted-foreground hover:text-white"
        >
          Cancel
        </Button>
        <Button 
          type="submit" 
          disabled={!projectName || !file || isUploading}
          className="bg-primary text-primary-foreground font-bold px-8 shadow-neon hover:shadow-neon-strong transition-all"
        >
          {isUploading ? (
            <>
              <Loader2 className="w-4 h-4 mr-2 animate-spin" />
              Creating Project...
            </>
          ) : (
            "Deploy Project"
          )}
        </Button>
      </div>
    </form>
  );
}
