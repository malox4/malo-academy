/** Map a desk/chat move onto board lights, broken nodes, extra DR, ticket heat. */
export function applyPlay(board, ev) {
  const nodes = board?.nodes || [];
  const debit = board?.debit || [];
  const lines = board?.lines || [];
  const mid = nodes[Math.min(1, Math.max(nodes.length - 1, 0))]?.id;
  const hit = ev.hit || ev.lineId || mid || nodes[0]?.id || debit[0]?.id || lines[0]?.id || "";
  const lights = { ...(ev.lights || {}) };
  let broken = ev.good ? [] : [hit].filter(Boolean);
  if (ev.broken) broken = ev.broken;

  if (ev.good) {
    if (hit) lights[hit] = lights[hit] || "ok";
    for (const n of nodes) if (!lights[n.id]) lights[n.id] = "ok";
  } else if (hit && !lights[hit]) {
    lights[hit] = "bad";
  }

  let extraDebit = ev.extraDebit || [];
  if (!extraDebit.length && !ev.good && board?.kind === "taccount") {
    extraDebit = [{ id: "fx-ghost", text: "DR лишняя нога · ход без правила", ok: false }];
  }

  let tickets = ev.tickets || {};
  if (!Object.keys(tickets).length && board?.kind === "ticket" && (ev.lineId || hit)) {
    const id = ev.lineId || hit;
    tickets = { [id]: ev.good ? "ok" : "bad" };
  }

  return {
    broken,
    lights,
    extraDebit,
    tickets,
    pulse: hit,
    caption: ev.why || "",
    flash: ev.good ? "ok" : "bad",
  };
}
