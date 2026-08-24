import { NavLink } from "react-router-dom";
import { Award, BookOpen, Cable, Compass, Map, Mic2, Radio, ScrollText, Sparkles, Swords, UserRound, Workflow } from "lucide-react";
import { learnerRank, overallProgress, useProgress } from "@/stores/progressStore";
import { cn } from "@/lib/cn";

const links = [
  { to: "/", label: "Путь", icon: Map, end: true },
  { to: "/pet", label: "API", icon: Cable },
  { to: "/book", label: "Справочник", icon: BookOpen },
  { to: "/bpmn", label: "BPMN", icon: Workflow },
  { to: "/lab", label: "Практика", icon: Swords },
  { to: "/quest", label: "Квест", icon: Radio },
  { to: "/interview", label: "Собеседование", icon: Mic2 },
  { to: "/exam", label: "Симуляция", icon: ScrollText },
  { to: "/achievements", label: "Достижения", icon: Award },
  { to: "/profile", label: "Профиль", icon: UserRound },
];

export function Sidebar() {
  const xp = useProgress((s) => s.xp);
  const name = useProgress((s) => s.learnerName);
  const streak = useProgress((s) => s.streak) ?? 0;
  const combo = useProgress((s) => s.combo) ?? 0;
  const progress = overallProgress();
  const rank = learnerRank(xp);

  return (
    <aside className="fixed inset-y-0 left-0 z-30 hidden w-[260px] flex-col border-r border-white/5 bg-[#0d1322]/80 px-5 py-6 backdrop-blur-xl md:flex">
      <div className="mb-8 flex items-center gap-3">
        <div className="grid h-11 w-11 place-items-center rounded-2xl border border-gold/30 bg-gold/10 text-gold glow-gold">
          <Sparkles size={18} />
        </div>
        <div>
          <div className="font-display text-lg leading-none text-gold-2">Malo</div>
          <div className="mt-1 text-[11px] uppercase tracking-[0.22em] text-muted">зал · пет живой</div>
        </div>
      </div>

      <nav className="flex flex-1 flex-col gap-1">
        {links.map((l) => (
          <NavLink
            key={l.to}
            to={l.to}
            end={l.end}
            className={({ isActive }) =>
              cn(
                "flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm transition",
                isActive
                  ? "bg-white/6 text-gold-2 shadow-[inset_0_0_0_1px_rgba(212,165,116,0.22)]"
                  : "text-muted hover:bg-white/4 hover:text-paper",
              )
            }
          >
            <l.icon size={16} />
            {l.label}
          </NavLink>
        ))}
      </nav>

      <div className="glass rounded-2xl p-4">
        <div className="flex items-center justify-between text-xs text-muted">
          <span className="inline-flex items-center gap-1">
            <Compass size={12} /> {name || "Путник"}
          </span>
          <span className="text-gold">{xp} XP</span>
        </div>
        <div className="mt-1 text-[11px] text-gold-2">{rank.title}</div>
        <div className="mt-3 h-1.5 overflow-hidden rounded-full bg-white/8">
          <div
            className="h-full rounded-full bg-gradient-to-r from-mint to-gold transition-all duration-700"
            style={{ width: `${progress}%` }}
          />
        </div>
        <div className="mt-2 flex justify-between text-[11px] text-muted">
          <span>Путь {progress}%</span>
          <span>
            {streak}д{combo >= 3 ? ` · ×${combo}` : ""}
          </span>
        </div>
      </div>
    </aside>
  );
}
