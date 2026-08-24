import { lazy, Suspense } from "react";
import { Link } from "react-router-dom";
import { Workflow } from "lucide-react";
import { PageMotion } from "@/components/ui/PageMotion";
import { Pill } from "@/components/ui/Pill";

const BpmnEditor = lazy(() => import("@/components/bpmn/BpmnEditor").then((m) => ({ default: m.BpmnEditor })));

const LEGEND = [
  { label: "Событие", className: "bg-[#C8E6C9] text-[#1B5E20]" },
  { label: "Пользовательская задача", className: "bg-[#FFE0B2] text-[#E65100]" },
  { label: "Сервисная задача", className: "bg-[#E1BEE7] text-[#6A1B9A]" },
  { label: "Исключение", className: "bg-[#FFCDD2] text-[#C62828]" },
];

export function BpmnPage() {
  return (
    <PageMotion>
      <Pill tone="violet">BPMN 2.0 · редактор</Pill>
      <div className="mt-4 flex flex-wrap items-end justify-between gap-4">
        <div className="max-w-2xl">
          <h1 className="font-display flex items-center gap-3 text-4xl md:text-5xl">
            <Workflow className="text-violet" size={32} />
            Процесс банкомата
          </h1>
          <p className="mt-3 text-muted">
            Живая схема, не скрин. Двойной клик меняет подпись, палитра слева добавляет элементы, колесо —
            масштаб. Скачайте <code className="text-gold">.bpmn</code> и откройте в Camunda Modeler или{" "}
            <a className="text-gold underline-offset-2 hover:underline" href="https://demo.bpmn.io" target="_blank" rel="noreferrer">
              demo.bpmn.io
            </a>
            .
          </p>
        </div>
        <Link to="/module/junior-3-bpmn" className="text-sm text-gold hover:text-gold-2">
          Модуль Junior · BPMN →
        </Link>
      </div>

      <div className="mt-5 flex flex-wrap gap-2">
        {LEGEND.map((item) => (
          <span key={item.label} className={`rounded-full px-3 py-1 text-[11px] font-medium ${item.className}`}>
            {item.label}
          </span>
        ))}
      </div>

      <div className="mt-6">
        <Suspense fallback={<div className="h-[68vh] min-h-[520px] animate-pulse rounded-3xl bg-white/4" />}>
          <BpmnEditor src="/bpmn/atm-process.bpmn" storageKey="malo-bpmn-atm" />
        </Suspense>
      </div>

      <p className="mt-4 max-w-3xl text-xs leading-5 text-muted">
        Дорожки: Клиент · Банкомат · Процессинг · Система выдачи. XOR на карте и ПИН, цикл трёх попыток, изъятие,
        параллель баланс + снятие. AND после «Выбрать операцию» запускает оба пути — на практике выбор операции
        чаще XOR. Поправьте схему прямо здесь.
      </p>
    </PageMotion>
  );
}
