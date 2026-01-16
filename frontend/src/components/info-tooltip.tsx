
import { Info } from "lucide-react"
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip"


export function InfoTooltip({ content }: { content: string }) {
  return <Tooltip>
    <TooltipTrigger>
      <Info className="h-3.5 w-3.5 text-muted-foreground/70 hover:text-primary transition-colors cursor-help" />
    </TooltipTrigger>
    <TooltipContent className="max-w-[250px] font-normal tracking-normal normal-case">
      <p>{content}</p>
    </TooltipContent>
  </Tooltip>
}