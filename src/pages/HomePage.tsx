import { Link } from "react-router-dom";
import { motion } from "framer-motion";
import { ArrowRight, BookOpen, Cable, Compass, Lock, Radio, Sparkles, Swords, Unlock, Workflow } from "lucide-react";
import { CURRICULUM } from "@/content/curriculum";
import { PageMotion } from "@/components/ui/PageMotion";
import { Pill } from "@/components/ui/Pill";
import { ProgressBar } from "@/components/ui/ProgressBar";
import { gradeProgress, learnerRank, overallProgress, useProgress } from "@/stores/progressStore";
import { PathMap } from "@/components/academy/PathMap";

export function HomePage() {
  const name = useProgress((s) => s.learnerName);
  const track = useProgress((s) => s.track);
  const setTrack = useProgress((s) => s.setTrack);
  const setName = useProgress((s) => s.setName);
  const xp = useProgress((s) => s.xp);
  const overall = overallProgress();
  const streak = useProgress((s) => s.streak) ?? 0;
  const daily = useProgress((s) => s.daily) ?? { quiz: false, practice: false, drill: false, lab: false, quest: false, day: "" };
  const rank = learnerRank(xp);

  return (
    <PageMotion>
      <div className="flex flex-wrap items-end justify-between gap-6">
        <div className="max-w-2xl">
          <Pill>BA · SA · живой контур, не слайд</Pill>
          <h1 className="font-display mt-4 text-4xl leading-[1.1] text-paper md:text-6xl">
            Malo Academy
          </h1>
          <p className="mt-4 max-w-xl text-base leading-relaxed text-muted md:text-lg">
            Стреляйте в API, пишите ноги журнала, закрывайте смены. Холд, capture, IBAN, FX — контур помнит, что вы сделали.
          </p>
        </div>
        <div className="glass w-full max-w-sm rounded-3xl p-5">
          <label className="text-[11px] uppercase tracking-[0.18em] text-muted">Ваше имя в зале</label>
          <input
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Как к вам обращаться"
            className="mt-2 w-full border-b border-white/10 bg-transparent py-2 text-lg outline-none placeholder:text-muted/50"
          />
          <div className="mt-4 flex items-center justify-between text-sm">
            <span className="text-muted">{xp} XP · {rank.title}</span>
            <span className="text-gold">{overall}% пути</span>
          </div>
          <ProgressBar value={overall} className="mt-2" />
        </div>
      </div>

      <div className="mt-8 flex flex-wrap items-center gap-3">
        <button
          onClick={() => setTrack("linear")}
          className={`inline-flex items-center gap-2 rounded-full border px-4 py-2 text-sm transition ${
            track === "linear"
              ? "border-gold/40 bg-gold/15 text-gold-2"
              : "border-white/10 text-muted hover:text-paper"
          }`}
        >
          <Lock size={14} /> Линейный путь
        </button>
        <button
          onClick={() => setTrack("free")}
          className={`inline-flex items-center gap-2 rounded-full border px-4 py-2 text-sm transition ${
            track === "free"
              ? "border-mint/40 bg-mint/15 text-mint"
              : "border-white/10 text-muted hover:text-paper"
          }`}
        >
          <Unlock size={14} /> Свободное изучение
        </button>
        <span className="text-xs text-muted">
          {track === "linear" ? "Открывается следующий модуль после закрытия текущего." : "Любой модуль доступен сразу."}
        </span>
      </div>

      <div className="mt-6 grid gap-3 sm:grid-cols-2 lg:grid-cols-5">
        {[
          { ok: daily.quiz, label: "Квиз дня" },
          { ok: daily.practice, label: "Практика" },
          { ok: daily.drill, label: "Мини-игра" },
          { ok: daily.lab, label: "Смена в зале" },
          { ok: daily.quest, label: "Квест" },
        ].map((q) => (
          <div key={q.label} className={`glass rounded-2xl px-4 py-3 text-sm ${q.ok ? "border-mint/30 text-mint" : "text-muted"}`}>
            {q.ok ? "✓" : "○"} {q.label}
          </div>
        ))}
      </div>
      <p className="mt-2 text-xs text-muted">Серия: {streak} дн. Откройте модуль — сразу вкладка «Играть».</p>

      <PathMap />

      <div className="mt-10 grid gap-4 md:grid-cols-2">
        {CURRICULUM.map((g, i) => {
          const p = gradeProgress(i);
          return (
            <motion.div
              key={g.id}
              initial={{ opacity: 0, y: 16 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.08 * i }}
            >
              <Link
                to={`/grade/${g.id}`}
                className="glass group block rounded-3xl p-6 transition hover:border-gold/30"
                style={{ boxShadow: `0 0 48px ${g.glow}` }}
              >
                <div className="flex items-start justify-between">
                  <div>
                    <div className="text-[11px] uppercase tracking-[0.22em] text-muted">Грейд {g.roman}</div>
                    <h2 className="font-display mt-1 text-3xl" style={{ color: g.color }}>
                      {g.title}
                    </h2>
                  </div>
                  <ArrowRight className="text-muted transition group-hover:translate-x-1 group-hover:text-gold" size={18} />
                </div>
                <p className="mt-3 text-sm leading-relaxed text-muted">{g.tagline}</p>
                <div className="mt-5 flex items-center justify-between text-xs text-muted">
                  <span>{g.levels.length} уровня · {g.levels.flatMap((l) => l.moduleIds).length} модулей</span>
                  <span style={{ color: g.color }}>{p}%</span>
                </div>
                <ProgressBar value={p} className="mt-2" />
              </Link>
            </motion.div>
          );
        })}
      </div>

      <div className="mt-8 grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Link to="/book" className="glass rounded-3xl p-5 transition hover:border-gold/30">
          <BookOpen className="text-gold" size={18} />
          <div className="font-display mt-3 text-xl">Справочник</div>
          <p className="mt-2 text-sm text-muted">
            Дебет, кредит, холд, nostro, FX, trailer — с примером ноги и API. Открывается рядом с квестом.
          </p>
        </Link>
        <Link to="/bpmn" className="glass rounded-3xl p-5 transition hover:border-gold/30">
          <Workflow className="text-violet" size={18} />
          <div className="font-display mt-3 text-xl">BPMN · банкомат</div>
          <p className="mt-2 text-sm text-muted">
            Редактируемая схема снятия: дорожки, XOR/AND, изъятие карты. Скачайте .bpmn в Camunda.
          </p>
        </Link>
        <Link to="/pet" className="glass rounded-3xl p-5 transition hover:border-gold/30">
          <Cable className="text-rose" size={18} />
          <div className="font-display mt-3 text-xl">API · живой пет</div>
          <p className="mt-2 text-sm text-muted">
            Postman на /api/v1: журнал, T-счета, холд→capture, IBAN/MT103, FX, зарплата, клиринг, сверка Orient.
          </p>
        </Link>
        <Link to="/quest" className="glass rounded-3xl p-5 transition hover:border-gold/30">
          <Radio className="text-rose" size={18} />
          <div className="font-display mt-3 text-xl">Квесты · живые продукты</div>
          <p className="mt-2 text-sm text-muted">
            Wallet, Core ledger, ShopLine, MedQueue, CityPark. Пишете ноги журнала, AC, 409, SMS. Контур помнит текст.
          </p>
        </Link>
        <Link to="/lab" className="glass rounded-3xl p-5 transition hover:border-gold/30">
          <Swords className="text-rose" size={18} />
          <div className="font-display mt-3 text-xl">Зал практики</div>
          <p className="mt-2 text-sm text-muted">Миссии: KYC, журнал, OpenAPI, SQL, файл банка, go/no-go. Не только Wallet.</p>
        </Link>
        <Link to="/interview" className="glass rounded-3xl p-5 transition hover:border-gold/30">
          <Compass className="text-gold" size={18} />
          <div className="font-display mt-3 text-xl">Собеседование</div>
          <p className="mt-2 text-sm text-muted">Типичные вопросы и сильные ответы — без зубрёжки определений.</p>
        </Link>
        <Link to="/exam" className="glass rounded-3xl p-5 transition hover:border-gold/30">
          <Sparkles className="text-violet" size={18} />
          <div className="font-display mt-3 text-xl">Финальная симуляция</div>
          <p className="mt-2 text-sm text-muted">Большой ситуационный экзамен: думайте, а не угадывайте термины.</p>
        </Link>
        <Link to="/achievements" className="glass rounded-3xl p-5 transition hover:border-gold/30">
          <Sparkles className="text-mint" size={18} />
          <div className="font-display mt-3 text-xl">Печати и бейджи</div>
          <p className="mt-2 text-sm text-muted">Геймификация без шума: печати грейдов, XP и редкие знаки.</p>
        </Link>
      </div>
    </PageMotion>
  );
}
