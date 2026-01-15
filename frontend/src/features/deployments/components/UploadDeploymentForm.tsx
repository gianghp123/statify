"use client";

import { useState } from "react";
import { Upload, X, FileArchive, Info, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { createDeployment } from "../services/deployment.actions";
import { toast } from "sonner";
import { useRouter } from "next/navigation";
import { MAX_FILE_SIZE } from "@/lib/constants";

interface UploadDeploymentFormProps {
  projectId: number;
  onSuccess?: () => void;
  onCancel?: () => void;
}

export function UploadDeploymentForm({ projectId, onSuccess, onCancel }: UploadDeploymentFormProps) {
  const router = useRouter();
  const [file, setFile] = useState<File | null>(null);
  const [isPending, setIsPending] = useState(false);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      const selectedFile = e.target.files[0];

      if (selectedFile.name.endsWith(".zip") || selectedFile.type !== 'application/zip') {
        if (selectedFile.size > MAX_FILE_SIZE) {
          toast.error("File is too large. Maximum size is 50MB.");
          return;
        }
        setFile(selectedFile);
      } else {
        toast.error("Please upload a .zip file");
      }
    }

  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!file) return;

    setIsPending(true);
    try {
      const res = await createDeployment(projectId, file);

      if (res.success) {
        toast.success("Deployment triggered successfully!");
        setFile(null);
        if (onSuccess) {
          onSuccess();
        }
        router.refresh();
      } else {
        toast.error(res.message || "Failed to trigger deployment");
      }
    } catch (error) {
      console.error(error);
      toast.error("An unexpected error occurred");
    } finally {
      setIsPending(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="space-y-4">
        <div className="space-y-2">
          <Label className="text-foreground font-bold text-lg">Upload build files</Label>
          <div
            className={`relative border-2 border-dashed rounded-xl p-10 transition-all flex flex-col items-center justify-center gap-4 ${file ? 'border-primary/50 bg-primary/5 shadow-inner' : 'border-border hover:border-primary/40 bg-muted/20 hover:bg-muted/30'
              }`}
          >
            <input
              type="file"
              accept=".zip"
              onChange={handleFileChange}
              className="absolute inset-0 opacity-0 cursor-pointer"
              required={!file}
              disabled={isPending}
            />

            {file ? (
              <>
                <div className="size-16 rounded-2xl bg-primary/20 flex items-center justify-center text-primary shadow-[0_0_15px_var(--neon-brand-glow)]">
                  <FileArchive className="w-8 h-8" />
                </div>
                <div className="text-center">
                  <p className="text-foreground font-extrabold">{file.name}</p>
                  <p className="text-xs text-muted-foreground font-medium">{(file.size / 1024 / 1024).toFixed(2)} MB</p>
                </div>
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  onClick={() => setFile(null)}
                  className="text-destructive hover:text-destructive hover:bg-destructive/10 font-bold transition-all"
                  disabled={isPending}
                >
                  <X className="w-4 h-4 mr-2" />
                  Remove file
                </Button>
              </>
            ) : (
              <>
                <div className="size-16 rounded-2xl bg-muted/40 border border-border/50 flex items-center justify-center text-muted-foreground shadow-sm">
                  <Upload className="w-8 h-8" />
                </div>
                <div className="text-center">
                  <p className="text-foreground font-extrabold text-xl">Click or drag to upload</p>
                  <p className="text-sm text-muted-foreground mt-1 font-medium">
                    Upload a ZIP folder containing your build files
                  </p>
                </div>
              </>
            )}
          </div>
          <div className="flex items-start gap-3 p-4 rounded-xl bg-accent/10 border border-border/50 text-xs text-foreground/80 leading-relaxed shadow-sm">
            <Info className="w-5 h-5 text-primary shrink-0" />
            <p className="font-medium">
              Make sure your ZIP archive contains an <span className="text-primary font-mono font-bold bg-muted/50 px-1 rounded">index.html</span> file
              at the root level or your project's build output directory.
            </p>
          </div>
        </div>
      </div>

      <div className="flex items-center gap-4 pt-4">
        {onCancel && (
          <Button
            type="button"
            variant="ghost"
            onClick={onCancel}
            className="text-muted-foreground hover:text-foreground font-bold"
            disabled={isPending}
          >
            Cancel
          </Button>
        )}
        <Button
          type="submit"
          disabled={!file || isPending}
          className="bg-primary text-primary-foreground font-black px-10 h-12 rounded-xl shadow-neon-brand hover:scale-[1.02] transition-all flex-1"
        >
          {isPending ? (
            <>
              <Loader2 className="w-5 h-5 mr-2 animate-spin" />
              Deploying...
            </>
          ) : (
            "Deploy Now"
          )}
        </Button>
      </div>
    </form>
  );
}
