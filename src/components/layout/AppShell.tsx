import { useEffect, type ReactNode } from "react";
import { useLocation } from "react-router-dom";
import { Sidebar } from "./Sidebar";
import { Ambient } from "./Ambient";
import { MobileBar } from "./MobileBar";
import { XpToast } from "@/components/ui/XpToast";
import { useProgress } from "@/stores/progressStore";
import { cn } from "@/lib/cn";

export function AppShell({ children }: { children: ReactNode }) {
  const touch = useProgress((s) => s.touchStreak);
  const location = useLocation();
  const wide = location.pathname === "/bpmn" || location.pathname === "/module/junior-3-bpmn";
  useEffect(() => {
    touch();
  }, [touch]);

  return (
    <div className="relative min-h-screen text-paper">
      <Ambient />
      <Sidebar />
      <main className="relative ml-0 min-h-screen pb-20 md:ml-[260px] md:pb-0">
        <div className={cn("mx-auto px-5 py-8 md:px-10 md:py-10", wide ? "max-w-[1600px]" : "max-w-6xl")}>{children}</div>
      </main>
      <MobileBar />
      <XpToast />
    </div>
  );
}
