import { lazy, Suspense } from "react";
import type { BpmnBlock } from "@/types/content";

const BpmnEditor = lazy(() => import("@/components/bpmn/BpmnEditor").then((m) => ({ default: m.BpmnEditor })));

export function BpmnView({ block }: { block: BpmnBlock }) {
  return (
    <section className="glass rounded-3xl p-6 md:p-8">
      <h3 className="font-display text-2xl">{block.title}</h3>
      <p className="mt-2 text-sm text-muted">{block.hint}</p>
      <div className="mt-5">
        <Suspense fallback={<div className="h-[560px] animate-pulse rounded-3xl bg-white/4" />}>
          <BpmnEditor src={block.src} storageKey={block.storageKey} heightClass="h-[560px] min-h-[480px]" />
        </Suspense>
      </div>
    </section>
  );
}
