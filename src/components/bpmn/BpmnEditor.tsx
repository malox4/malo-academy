import { useCallback, useEffect, useId, useRef, useState } from "react";
import BpmnModeler from "bpmn-js/lib/Modeler";
import { Download, FolderOpen, RotateCcw, Undo2, Redo2, Maximize2, AlertCircle } from "lucide-react";
import { cn } from "@/lib/cn";
import "bpmn-js/dist/assets/diagram-js.css";
import "bpmn-js/dist/assets/bpmn-js.css";
import "bpmn-js/dist/assets/bpmn-font/css/bpmn-embedded.css";

type CanvasService = { zoom: (level: string | number) => void };
type CommandStackService = {
  canUndo: () => boolean;
  canRedo: () => boolean;
  undo: () => void;
  redo: () => void;
};
type EventBusService = { on: (event: string, fn: () => void) => void };

function downloadFile(name: string, body: string, mime: string) {
  const blob = new Blob([body], { type: mime });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = name;
  a.click();
  URL.revokeObjectURL(url);
}

export function BpmnEditor({
  src,
  storageKey = "malo-bpmn-atm",
  heightClass = "h-[68vh] min-h-[520px]",
}: {
  src: string;
  storageKey?: string;
  heightClass?: string;
}) {
  const hostRef = useRef<HTMLDivElement>(null);
  const modelerRef = useRef<BpmnModeler | null>(null);
  const originalRef = useRef("");
  const saveTimer = useRef<number>(0);
  const fileRef = useRef<HTMLInputElement>(null);
  const reactId = useId();
  const [status, setStatus] = useState<"loading" | "ready" | "error">("loading");
  const [error, setError] = useState("");
  const [dirty, setDirty] = useState(false);
  const [canUndo, setCanUndo] = useState(false);
  const [canRedo, setCanRedo] = useState(false);

  const syncStack = useCallback(() => {
    const stack = modelerRef.current?.get("commandStack") as CommandStackService | undefined;
    if (!stack) return;
    setCanUndo(stack.canUndo());
    setCanRedo(stack.canRedo());
  }, []);

  const persist = useCallback(() => {
    const modeler = modelerRef.current;
    if (!modeler) return;
    window.clearTimeout(saveTimer.current);
    saveTimer.current = window.setTimeout(async () => {
      try {
        const { xml } = await modeler.saveXML({ format: true });
        if (xml) localStorage.setItem(storageKey, xml);
        setDirty(xml !== originalRef.current);
      } catch {
        /* ignore autosave failures */
      }
    }, 400);
  }, [storageKey]);

  const loadXml = useCallback(
    async (xml: string, markDirty = false) => {
      const modeler = modelerRef.current;
      if (!modeler) return;
      await modeler.importXML(xml);
      const canvas = modeler.get("canvas") as CanvasService;
      requestAnimationFrame(() => canvas.zoom("fit-viewport"));
      setDirty(markDirty);
      setStatus("ready");
      setError("");
      syncStack();
    },
    [syncStack],
  );

  useEffect(() => {
    const el = hostRef.current;
    if (!el) return;
    const modeler = new BpmnModeler({
      container: el,
    });
    modelerRef.current = modeler;
    const bus = modeler.get("eventBus") as EventBusService;
    bus.on("commandStack.changed", () => {
      syncStack();
      persist();
    });

    let cancelled = false;
    (async () => {
      try {
        const res = await fetch(src);
        if (!res.ok) throw new Error(`Не удалось загрузить ${src}`);
        const factory = await res.text();
        if (cancelled) return;
        originalRef.current = factory;
        const saved = localStorage.getItem(storageKey);
        await loadXml(saved && saved !== factory ? saved : factory, Boolean(saved && saved !== factory));
      } catch (e) {
        if (cancelled) return;
        setStatus("error");
        setError(e instanceof Error ? e.message : "Не удалось открыть BPMN");
      }
    })();

    return () => {
      cancelled = true;
      window.clearTimeout(saveTimer.current);
      modeler.destroy();
      modelerRef.current = null;
      el.replaceChildren();
    };
  }, [src, storageKey, loadXml, persist, syncStack]);

  const onDownloadBpmn = async () => {
    const { xml } = await modelerRef.current!.saveXML({ format: true });
    if (xml) downloadFile("atm-process.bpmn", xml, "application/bpmn20-xml");
  };

  const onDownloadSvg = async () => {
    const { svg } = await modelerRef.current!.saveSVG();
    if (svg) downloadFile("atm-process.svg", svg, "image/svg+xml");
  };

  const onReset = async () => {
    localStorage.removeItem(storageKey);
    await loadXml(originalRef.current, false);
  };

  const onOpenFile = async (file: File) => {
    const xml = await file.text();
    await loadXml(xml, true);
    localStorage.setItem(storageKey, xml);
  };

  return (
    <div className="space-y-3">
      <div className="flex flex-wrap items-center gap-2">
        <button
          type="button"
          disabled={!canUndo}
          onClick={() => (modelerRef.current?.get("commandStack") as CommandStackService).undo()}
          className="inline-flex items-center gap-1.5 rounded-full border border-white/12 px-3 py-1.5 text-xs text-muted disabled:opacity-30"
        >
          <Undo2 size={12} /> Отменить
        </button>
        <button
          type="button"
          disabled={!canRedo}
          onClick={() => (modelerRef.current?.get("commandStack") as CommandStackService).redo()}
          className="inline-flex items-center gap-1.5 rounded-full border border-white/12 px-3 py-1.5 text-xs text-muted disabled:opacity-30"
        >
          <Redo2 size={12} /> Повторить
        </button>
        <button
          type="button"
          onClick={() => (modelerRef.current?.get("canvas") as CanvasService).zoom("fit-viewport")}
          className="inline-flex items-center gap-1.5 rounded-full border border-white/12 px-3 py-1.5 text-xs text-muted"
        >
          <Maximize2 size={12} /> Вписать
        </button>
        <span className="hidden h-4 w-px bg-white/10 sm:block" />
        <button
          type="button"
          onClick={onDownloadBpmn}
          disabled={status !== "ready"}
          className="inline-flex items-center gap-1.5 rounded-full bg-gold px-3 py-1.5 text-xs font-medium text-ink disabled:opacity-40"
        >
          <Download size={12} /> Скачать .bpmn
        </button>
        <button
          type="button"
          onClick={onDownloadSvg}
          disabled={status !== "ready"}
          className="inline-flex items-center gap-1.5 rounded-full border border-white/12 px-3 py-1.5 text-xs text-muted disabled:opacity-40"
        >
          <Download size={12} /> SVG
        </button>
        <button
          type="button"
          onClick={() => fileRef.current?.click()}
          className="inline-flex items-center gap-1.5 rounded-full border border-white/12 px-3 py-1.5 text-xs text-muted"
        >
          <FolderOpen size={12} /> Открыть файл
        </button>
        <button
          type="button"
          onClick={onReset}
          disabled={status !== "ready"}
          className="inline-flex items-center gap-1.5 rounded-full border border-white/12 px-3 py-1.5 text-xs text-muted disabled:opacity-40"
        >
          <RotateCcw size={12} /> Сбросить
        </button>
        <input
          ref={fileRef}
          type="file"
          accept=".bpmn,.xml,application/xml,text/xml"
          className="hidden"
          onChange={(e) => {
            const file = e.target.files?.[0];
            if (file) void onOpenFile(file);
            e.target.value = "";
          }}
        />
        {dirty && <span className="text-[11px] text-gold">черновик в браузере</span>}
      </div>

      <div className={cn("bpmn-host relative overflow-hidden rounded-3xl border border-white/10", heightClass)}>
        <div ref={hostRef} id={reactId} className="absolute inset-0" />
        {status === "loading" && (
          <div className="absolute inset-0 grid place-items-center bg-[#f4f1ea] text-sm text-[#5c6578]">Загрузка схемы…</div>
        )}
        {status === "error" && (
          <div className="absolute inset-0 grid place-items-center bg-[#f4f1ea] px-6 text-center text-sm text-[#8a3b3b]">
            <span className="inline-flex items-center gap-2">
              <AlertCircle size={16} /> {error}
            </span>
          </div>
        )}
      </div>
    </div>
  );
}
