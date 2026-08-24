/**
 * Builds public/bpmn/atm-process.bpmn — BPMN 2.0 + DI for Camunda / bpmn.io.
 * Run: node scripts/build-atm-bpmn.mjs
 */
import { writeFileSync, mkdirSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const TASK = { w: 120, h: 80 };
const GW = 50;
const EV = 36;

const color = {
  event: { fill: "#C8E6C9", stroke: "#2E7D32" },
  eventEnd: { fill: "#A5D6A7", stroke: "#1B5E20" },
  eventError: { fill: "#EF9A9A", stroke: "#B71C1C" },
  user: { fill: "#FFE0B2", stroke: "#E65100" },
  service: { fill: "#E1BEE7", stroke: "#6A1B9A" },
  error: { fill: "#FFCDD2", stroke: "#C62828" },
  gateway: { fill: "#FFFDE7", stroke: "#F9A825" },
};

const esc = (s) =>
  String(s)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;");

const box = (id, x, y, w, h) => ({ id, x, y, w, h, cx: x + w / 2, cy: y + h / 2 });
const taskBox = (id, x, y) => box(id, x, y, TASK.w, TASK.h);
const gwBox = (id, x, y) => box(id, x, y, GW, GW);
const evBox = (id, x, y) => box(id, x, y, EV, EV);

const edge = {
  right: (s) => ({ x: s.x + s.w, y: s.cy }),
  left: (s) => ({ x: s.x, y: s.cy }),
  top: (s) => ({ x: s.cx, y: s.y }),
  bottom: (s) => ({ x: s.cx, y: s.y + s.h }),
};

/* --- layout: one pool, four lanes, left → right --- */
const shapes = {
  Start_Begin: evBox("Start_Begin", 230, 137),
  Task_InsertCard: taskBox("Task_InsertCard", 310, 115),
  Task_CheckCard: taskBox("Task_CheckCard", 480, 390),
  Gw_CardSupported: gwBox("Gw_CardSupported", 660, 405),
  Task_CardUnsupported: taskBox("Task_CardUnsupported", 780, 250),
  Task_ReturnCardBad: taskBox("Task_ReturnCardBad", 960, 250),
  Task_RequestPin: taskBox("Task_RequestPin", 800, 390),
  Task_CheckPin: taskBox("Task_CheckPin", 980, 390),
  Gw_PinCorrect: gwBox("Gw_PinCorrect", 1160, 405),
  Gw_PinAttempts: gwBox("Gw_PinAttempts", 1160, 530),
  Task_RetainCard: taskBox("Task_RetainCard", 1320, 635),
  End_Retained: evBox("End_Retained", 1520, 657),
  Task_OpenSession: taskBox("Task_OpenSession", 1320, 390),
  Task_SelectOp: taskBox("Task_SelectOp", 1500, 390),
  Gw_SplitOps: gwBox("Gw_SplitOps", 1680, 405),
  Task_ViewBalance: taskBox("Task_ViewBalance", 1820, 250),
  Task_ShowBalance: taskBox("Task_ShowBalance", 2000, 115),
  Task_CheckFunds: taskBox("Task_CheckFunds", 1820, 510),
  Gw_EnoughFunds: gwBox("Gw_EnoughFunds", 2000, 525),
  Task_Insufficient: taskBox("Task_Insufficient", 2160, 510),
  Task_EnterAmount: taskBox("Task_EnterAmount", 2160, 635),
  Task_Dispense: taskBox("Task_Dispense", 2340, 780),
  Task_PrintReceipt: taskBox("Task_PrintReceipt", 2520, 780),
  Gw_MergeWithdraw: gwBox("Gw_MergeWithdraw", 2700, 525),
  Gw_JoinOps: gwBox("Gw_JoinOps", 2700, 405),
  Task_CompleteOp: taskBox("Task_CompleteOp", 2840, 390),
  Task_ReturnCard: taskBox("Task_ReturnCard", 3020, 390),
  End_Session: evBox("End_Session", 3200, 407),
};

const POOL = { x: 160, y: 80, w: 3140, h: 860 };
const LANES = [
  { id: "Lane_Client", name: "Клиент", x: 190, y: 80, w: 3110, h: 150 },
  { id: "Lane_Atm", name: "Банкомат", x: 190, y: 230, w: 3110, h: 360 },
  { id: "Lane_Processing", name: "Процессинг", x: 190, y: 590, w: 3110, h: 140 },
  { id: "Lane_Issuing", name: "Система выдачи", x: 190, y: 730, w: 3110, h: 210 },
];

const laneNodes = {
  Lane_Client: ["Start_Begin", "Task_InsertCard", "Task_ShowBalance"],
  Lane_Atm: [
    "Task_CheckCard",
    "Gw_CardSupported",
    "Task_CardUnsupported",
    "Task_ReturnCardBad",
    "Task_RequestPin",
    "Task_CheckPin",
    "Gw_PinCorrect",
    "Gw_PinAttempts",
    "Task_OpenSession",
    "Task_SelectOp",
    "Gw_SplitOps",
    "Task_ViewBalance",
    "Task_CheckFunds",
    "Gw_EnoughFunds",
    "Task_Insufficient",
    "Gw_MergeWithdraw",
    "Gw_JoinOps",
    "Task_CompleteOp",
    "Task_ReturnCard",
    "End_Session",
  ],
  Lane_Processing: ["Task_RetainCard", "End_Retained", "Task_EnterAmount"],
  Lane_Issuing: ["Task_Dispense", "Task_PrintReceipt"],
};

const flows = [
  { id: "Flow_Start_Insert", src: "Start_Begin", tgt: "Task_InsertCard", from: "right", to: "left" },
  { id: "Flow_Insert_Check", src: "Task_InsertCard", tgt: "Task_CheckCard", from: "bottom", to: "top" },
  { id: "Flow_Check_GwCard", src: "Task_CheckCard", tgt: "Gw_CardSupported", from: "right", to: "left" },
  {
    id: "Flow_CardNo",
    src: "Gw_CardSupported",
    tgt: "Task_CardUnsupported",
    name: "Нет",
    from: "top",
    to: "left",
    mid: [{ x: 685, y: 290 }],
  },
  { id: "Flow_Unsupported_Return", src: "Task_CardUnsupported", tgt: "Task_ReturnCardBad", from: "right", to: "left" },
  {
    id: "Flow_ReturnBad_End",
    src: "Task_ReturnCardBad",
    tgt: "End_Session",
    from: "right",
    to: "top",
    mid: [
      { x: 1120, y: 290 },
      { x: 3218, y: 290 },
    ],
  },
  {
    id: "Flow_CardYes",
    src: "Gw_CardSupported",
    tgt: "Task_RequestPin",
    name: "Да (U2/Visa/Master)",
    from: "right",
    to: "left",
  },
  { id: "Flow_Request_CheckPin", src: "Task_RequestPin", tgt: "Task_CheckPin", from: "right", to: "left" },
  { id: "Flow_CheckPin_Gw", src: "Task_CheckPin", tgt: "Gw_PinCorrect", from: "right", to: "left" },
  {
    id: "Flow_PinNo",
    src: "Gw_PinCorrect",
    tgt: "Gw_PinAttempts",
    name: "Нет",
    from: "bottom",
    to: "top",
  },
  {
    id: "Flow_AttemptsNo",
    src: "Gw_PinAttempts",
    tgt: "Task_CheckPin",
    name: "Нет",
    from: "left",
    to: "bottom",
    mid: [
      { x: 1040, y: 555 },
      { x: 1040, y: 470 },
    ],
  },
  {
    id: "Flow_AttemptsYes",
    src: "Gw_PinAttempts",
    tgt: "Task_RetainCard",
    name: "Да",
    from: "right",
    to: "left",
    mid: [{ x: 1280, y: 555 }],
  },
  { id: "Flow_Retain_End", src: "Task_RetainCard", tgt: "End_Retained", from: "right", to: "left" },
  {
    id: "Flow_PinYes",
    src: "Gw_PinCorrect",
    tgt: "Task_OpenSession",
    name: "Да",
    from: "right",
    to: "left",
  },
  { id: "Flow_Open_Select", src: "Task_OpenSession", tgt: "Task_SelectOp", from: "right", to: "left" },
  { id: "Flow_Select_Split", src: "Task_SelectOp", tgt: "Gw_SplitOps", from: "right", to: "left" },
  {
    id: "Flow_Split_Balance",
    src: "Gw_SplitOps",
    tgt: "Task_ViewBalance",
    from: "top",
    to: "left",
    mid: [{ x: 1705, y: 290 }],
  },
  { id: "Flow_View_Show", src: "Task_ViewBalance", tgt: "Task_ShowBalance", from: "right", to: "left" },
  {
    id: "Flow_Show_Join",
    src: "Task_ShowBalance",
    tgt: "Gw_JoinOps",
    from: "right",
    to: "top",
    mid: [
      { x: 2160, y: 155 },
      { x: 2725, y: 155 },
    ],
  },
  {
    id: "Flow_Split_Funds",
    src: "Gw_SplitOps",
    tgt: "Task_CheckFunds",
    from: "bottom",
    to: "left",
    mid: [{ x: 1705, y: 550 }],
  },
  { id: "Flow_Funds_Gw", src: "Task_CheckFunds", tgt: "Gw_EnoughFunds", from: "right", to: "left" },
  {
    id: "Flow_FundsNo",
    src: "Gw_EnoughFunds",
    tgt: "Task_Insufficient",
    name: "Нет",
    from: "right",
    to: "left",
  },
  { id: "Flow_Insufficient_Merge", src: "Task_Insufficient", tgt: "Gw_MergeWithdraw", from: "right", to: "left" },
  {
    id: "Flow_FundsYes",
    src: "Gw_EnoughFunds",
    tgt: "Task_EnterAmount",
    name: "Да",
    from: "bottom",
    to: "left",
    mid: [{ x: 2025, y: 675 }],
  },
  { id: "Flow_Amount_Dispense", src: "Task_EnterAmount", tgt: "Task_Dispense", from: "right", to: "left" },
  { id: "Flow_Dispense_Print", src: "Task_Dispense", tgt: "Task_PrintReceipt", from: "right", to: "left" },
  {
    id: "Flow_Print_Merge",
    src: "Task_PrintReceipt",
    tgt: "Gw_MergeWithdraw",
    from: "right",
    to: "bottom",
    mid: [{ x: 2725, y: 820 }],
  },
  { id: "Flow_Merge_Join", src: "Gw_MergeWithdraw", tgt: "Gw_JoinOps", from: "top", to: "bottom" },
  { id: "Flow_Join_Complete", src: "Gw_JoinOps", tgt: "Task_CompleteOp", from: "right", to: "left" },
  { id: "Flow_Complete_Return", src: "Task_CompleteOp", tgt: "Task_ReturnCard", from: "right", to: "left" },
  { id: "Flow_Return_End", src: "Task_ReturnCard", tgt: "End_Session", from: "right", to: "left" },
];

const elements = [
  { id: "Start_Begin", kind: "startEvent", name: "Начало", color: color.event },
  { id: "Task_InsertCard", kind: "userTask", name: "Вставить карту", color: color.user },
  { id: "Task_CheckCard", kind: "serviceTask", name: "Проверить карту", color: color.service },
  { id: "Gw_CardSupported", kind: "exclusiveGateway", name: "Карта поддерживается?", color: color.gateway, labelDy: -28 },
  { id: "Task_CardUnsupported", kind: "serviceTask", name: "Карта не поддерживается", color: color.error },
  { id: "Task_ReturnCardBad", kind: "serviceTask", name: "Вернуть карту", color: color.error },
  { id: "Task_RequestPin", kind: "serviceTask", name: "Запросить ПИН-код", color: color.service },
  { id: "Task_CheckPin", kind: "serviceTask", name: "Проверить ПИН-код", color: color.service },
  { id: "Gw_PinCorrect", kind: "exclusiveGateway", name: "ПИН верный?", color: color.gateway, labelDy: -28 },
  { id: "Gw_PinAttempts", kind: "exclusiveGateway", name: "Попытка больше 3-х раз?", color: color.gateway, labelDy: 58, labelW: 160 },
  { id: "Task_RetainCard", kind: "serviceTask", name: "Изъять карту", color: color.error },
  { id: "End_Retained", kind: "endEvent", name: "Сессия завершена", color: color.eventError },
  { id: "Task_OpenSession", kind: "serviceTask", name: "Открыть сессию", color: color.service },
  { id: "Task_SelectOp", kind: "serviceTask", name: "Выбрать операцию", color: color.service },
  { id: "Gw_SplitOps", kind: "parallelGateway", name: "", color: color.gateway },
  { id: "Task_ViewBalance", kind: "serviceTask", name: "Просмотр баланса", color: color.service },
  { id: "Task_ShowBalance", kind: "userTask", name: "Показать баланс на экране", color: color.user },
  { id: "Task_CheckFunds", kind: "serviceTask", name: "Проверить доступность средств", color: color.service },
  { id: "Gw_EnoughFunds", kind: "exclusiveGateway", name: "Средств достаточно?", color: color.gateway, labelDy: -28, labelW: 150 },
  { id: "Task_Insufficient", kind: "serviceTask", name: "Недостаточно средств", color: color.error },
  { id: "Task_EnterAmount", kind: "userTask", name: "Ввести сумму", color: color.user },
  { id: "Task_Dispense", kind: "serviceTask", name: "Выдать наличные", color: color.service },
  { id: "Task_PrintReceipt", kind: "serviceTask", name: "Напечатать чек", color: color.service },
  { id: "Gw_MergeWithdraw", kind: "exclusiveGateway", name: "", color: color.gateway },
  { id: "Gw_JoinOps", kind: "parallelGateway", name: "", color: color.gateway },
  { id: "Task_CompleteOp", kind: "serviceTask", name: "Завершить операцию", color: color.service },
  { id: "Task_ReturnCard", kind: "serviceTask", name: "Вернуть карту", color: color.service },
  { id: "End_Session", kind: "endEvent", name: "Завершение сессии", color: color.eventEnd },
];

function incomingOutgoing() {
  const map = Object.fromEntries(elements.map((e) => [e.id, { in: [], out: [] }]));
  for (const f of flows) {
    map[f.src].out.push(f.id);
    map[f.tgt].in.push(f.id);
  }
  return map;
}

function elXml(el, io) {
  const name = el.name ? ` name="${esc(el.name)}"` : "";
  const ins = io[el.id].in.map((id) => `    <bpmn:incoming>${id}</bpmn:incoming>`).join("\n");
  const outs = io[el.id].out.map((id) => `    <bpmn:outgoing>${id}</bpmn:outgoing>`).join("\n");
  const body = [ins, outs].filter(Boolean).join("\n");
  if (el.kind === "startEvent" || el.kind === "endEvent" || el.kind.endsWith("Gateway")) {
    return `  <bpmn:${el.kind} id="${el.id}"${name}>\n${body}\n  </bpmn:${el.kind}>`;
  }
  return `  <bpmn:${el.kind} id="${el.id}"${name}>\n${body}\n  </bpmn:${el.kind}>`;
}

function flowXml(f) {
  const name = f.name ? ` name="${esc(f.name)}"` : "";
  return `  <bpmn:sequenceFlow id="${f.id}"${name} sourceRef="${f.src}" targetRef="${f.tgt}" />`;
}

function colorAttrs(c) {
  if (!c) return "";
  return ` bioc:stroke="${c.stroke}" bioc:fill="${c.fill}" color:background-color="${c.fill}" color:border-color="${c.stroke}"`;
}

function shapeDi(el) {
  const s = shapes[el.id];
  const label =
    el.name && (el.kind.endsWith("Gateway") || el.kind.endsWith("Event"))
      ? (() => {
          const w = el.labelW ?? Math.min(180, Math.max(80, el.name.length * 7));
          const dy = el.labelDy ?? (el.kind.endsWith("Event") ? 40 : -24);
          return `
      <bpmndi:BPMNLabel>
        <dc:Bounds x="${Math.round(s.cx - w / 2)}" y="${Math.round(s.y + dy)}" width="${w}" height="28" />
      </bpmndi:BPMNLabel>`;
        })()
      : "";
  return `    <bpmndi:BPMNShape id="${el.id}_di" bpmnElement="${el.id}"${colorAttrs(el.color)}>
      <dc:Bounds x="${s.x}" y="${s.y}" width="${s.w}" height="${s.h}" />${label}
    </bpmndi:BPMNShape>`;
}

function edgeDi(f) {
  const src = shapes[f.src];
  const tgt = shapes[f.tgt];
  const a = edge[f.from](src);
  const b = edge[f.to](tgt);
  const pts = [a, ...(f.mid ?? []), b];
  const waypoints = pts.map((p) => `      <di:waypoint x="${Math.round(p.x)}" y="${Math.round(p.y)}" />`).join("\n");
  let label = "";
  if (f.name) {
    const mid = pts[Math.floor(pts.length / 2)];
    const w = Math.min(170, Math.max(40, f.name.length * 7));
    label = `
      <bpmndi:BPMNLabel>
        <dc:Bounds x="${Math.round(mid.x - w / 2)}" y="${Math.round(mid.y - 22)}" width="${w}" height="18" />
      </bpmndi:BPMNLabel>`;
  }
  return `    <bpmndi:BPMNEdge id="${f.id}_di" bpmnElement="${f.id}">
${waypoints}${label}
    </bpmndi:BPMNEdge>`;
}

const io = incomingOutgoing();

const xml = `<?xml version="1.0" encoding="UTF-8"?>
<bpmn:definitions
  xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xmlns:bpmn="http://www.omg.org/spec/BPMN/20100524/MODEL"
  xmlns:bpmndi="http://www.omg.org/spec/BPMN/20100524/DI"
  xmlns:dc="http://www.omg.org/spec/DD/20100524/DC"
  xmlns:di="http://www.omg.org/spec/DD/20100524/DI"
  xmlns:bioc="http://bpmn.io/schema/bpmn/biocolor/1.0"
  xmlns:color="http://www.omg.org/spec/BPMN/non-normative/color/1.0"
  id="Definitions_AtmProcess"
  name="Процесс банкомата"
  targetNamespace="https://malo.academy/bpmn"
  exporter="Malo Academy"
  exporterVersion="1.1.0">
  <bpmn:collaboration id="Collaboration_Atm">
    <bpmn:participant id="Participant_Process" name="Процесс" processRef="Process_Atm" />
  </bpmn:collaboration>
  <bpmn:process id="Process_Atm" name="Снятие наличных в банкомате" isExecutable="false">
    <bpmn:documentation>As-is процесс банкомата: клиент, банкомат, процессинг, система выдачи. XOR карты и ПИН, цикл трёх попыток, изъятие карты, параллель баланс + снятие. Откройте в Camunda Modeler, bpmn.io или в зале Malo (/bpmn).</bpmn:documentation>
    <bpmn:laneSet id="LaneSet_Atm">
${LANES.map(
  (l) => `      <bpmn:lane id="${l.id}" name="${esc(l.name)}">
${laneNodes[l.id].map((n) => `        <bpmn:flowNodeRef>${n}</bpmn:flowNodeRef>`).join("\n")}
      </bpmn:lane>`,
).join("\n")}
    </bpmn:laneSet>
${elements.map((e) => elXml(e, io)).join("\n")}
${flows.map(flowXml).join("\n")}
  </bpmn:process>
  <bpmndi:BPMNDiagram id="BPMNDiagram_Atm">
    <bpmndi:BPMNPlane id="BPMNPlane_Atm" bpmnElement="Collaboration_Atm">
      <bpmndi:BPMNShape id="Participant_Process_di" bpmnElement="Participant_Process" isHorizontal="true">
        <dc:Bounds x="${POOL.x}" y="${POOL.y}" width="${POOL.w}" height="${POOL.h}" />
      </bpmndi:BPMNShape>
${LANES.map(
  (l) => `      <bpmndi:BPMNShape id="${l.id}_di" bpmnElement="${l.id}" isHorizontal="true">
        <dc:Bounds x="${l.x}" y="${l.y}" width="${l.w}" height="${l.h}" />
      </bpmndi:BPMNShape>`,
).join("\n")}
${elements.map(shapeDi).join("\n")}
${flows.map(edgeDi).join("\n")}
    </bpmndi:BPMNPlane>
  </bpmndi:BPMNDiagram>
</bpmn:definitions>
`;

const out = join(dirname(fileURLToPath(import.meta.url)), "..", "public", "bpmn", "atm-process.bpmn");
mkdirSync(dirname(out), { recursive: true });
writeFileSync(out, xml);
console.log("wrote", out, xml.length, "bytes");
