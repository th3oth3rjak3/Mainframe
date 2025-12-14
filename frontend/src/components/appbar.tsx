import { ChevronRight } from "lucide-react";
import { SidebarTrigger } from "@/components/ui/sidebar";

export function AppBar() {
  return (
    <div className="h-14 border-b bg-background px-4 flex items-center justify-between gap-4">
      {/* Left side with trigger and breadcrumbs */}
      <div className="flex items-center gap-4">
        <SidebarTrigger />
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <span className="text-foreground font-medium">Dashboard</span>
          <ChevronRight className="h-4 w-4" />
          <span>Recipes</span>
        </div>
      </div>

      {/* Right side actions */}
      <div className="flex items-center gap-2">{/* Add search or other icons here later */}</div>
    </div>
  );
}
