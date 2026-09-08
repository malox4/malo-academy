/** Articles: core rewritten as teaching; catalog and gap keep ids. */
import { stampTerms } from "./pedagogy.js";
import { CORE_TOPICS } from "./topicsCore.js";
import { EXTRA_TOPICS } from "./topicsCatalog.js";
import { GAP_TOPICS_A } from "./topicsGap.js";
import { GAP_TOPICS_B } from "./topicsGapMore.js";

export const TOPICS = [...CORE_TOPICS, ...EXTRA_TOPICS, ...GAP_TOPICS_A, ...GAP_TOPICS_B].map(stampTerms);

export const TOPIC_BY_ID = Object.fromEntries(TOPICS.map((t) => [t.id, t]));

export function topicById(id) {
  return TOPIC_BY_ID[id] || null;
}

export const LESSON_TOPIC = {
  "intern-1-profession": "car-analyst",
  "intern-1-sdlc": "req-levels",
  "intern-1-team": "req-stakeholders",
  "intern-1-questions": "car-intro",
  "intern-2-requirement": "req-need",
  "intern-2-levels": "req-levels",
  "intern-2-stakeholders": "req-stakeholders",
  "intern-2-quality": "req-quality",
  "intern-3-elicitation": "req-need",
  "intern-3-notes": "req-jira",
  "intern-3-story": "req-story",
  "intern-3-meeting": "car-intro",
  "junior-1-stories": "req-story",
  "junior-1-usecase": "req-story",
  "junior-1-ac": "req-ac",
  "junior-2-brd": "req-brd",
  "junior-2-srs": "req-quality",
  "junior-2-trace": "req-trace",
  "junior-2-write": "req-jira",
  "junior-3-uml": "arch-uml",
  "junior-3-bpmn": "arch-bpmn",
  "junior-3-gap": "arch-bpmn",
  "middle-2-api": "rest-http",
  "middle-2-sql": "db-sql",
  "middle-2-integration": "int-idem",
  "middle-3-nfr": "req-nfr",
  "middle-3-arch": "arch-c4",
  "middle-3-errors": "int-timeout",
};

export function topicForLesson(lessonId) {
  return topicById(LESSON_TOPIC[lessonId]);
}

export const LAB_TOPIC = {
  "lab-req-integration": "int-task",
  "lab-sql-pet": "db-sql",
  "lab-rest-desk": "rest-http",
  "lab-bpmn-hold": "arch-bpmn",
  "lab-bank-ledger": "arch-hold",
  "lab-intern-1": "bank-aml-kyc",
  "lab-intern-2": "int-idem",
  "lab-intern-3": "int-task",
  "lab-junior-1": "api-openapi",
  "lab-junior-2": "db-keys",
  "lab-junior-3": "bank-recon",
  "lab-middle-1": "arch-styles",
  "lab-middle-2": "db-sql",
  "lab-middle-3": "int-map",
  "lab-senior-1": "bank-aml-kyc",
  "lab-senior-2": "soft-workshop",
  "lab-senior-3": "req-estimate",
  "lab-c4-contour": "arch-c4",
  "lab-req-ac": "req-ac",
  "lab-req-brd": "req-brd",
  "lab-req-nfr": "req-nfr",
  "lab-req-trace": "req-trace",
  "lab-req-jira": "req-jira",
  "lab-sql-join": "db-join",
  "lab-er-class": "db-er",
  "lab-sql-m2m": "db-m2m",
  "lab-sql-index": "db-index",
  "lab-sql-json": "db-json",
  "lab-rest-codes": "rest-codes",
  "lab-rest-resource": "rest-resource",
  "lab-rest-version": "rest-version",
  "lab-webhook": "int-webhook",
  "lab-queue": "int-queue",
  "lab-soap": "int-soap",
  "lab-oauth": "int-oauth",
  "lab-mapping-dlq": "int-map",
  "lab-uml-seq": "arch-seq",
  "lab-iban": "arch-iban",
  "lab-p2p": "arch-p2p",
  "lab-bpmn-mistakes": "arch-bpmn-mistakes",
  "lab-req-validate": "req-validate",
  "lab-req-lifecycle": "req-lifecycle",
  "lab-req-test-analysis": "req-test-analysis",
  "lab-req-estimate": "req-estimate",
  "lab-req-gost": "req-gost-34602",
  "lab-db-layers": "db-layers",
  "lab-db-window": "db-window",
  "lab-db-cte": "db-cte",
  "lab-db-migrate": "db-migrate",
  "lab-db-dwh": "db-dwh",
  "lab-api-openapi": "api-openapi",
  "lab-api-postman": "api-postman",
  "lab-api-asyncapi": "api-asyncapi",
  "lab-api-graphql": "api-graphql",
  "lab-api-schema": "api-schema",
  "lab-api-pagination": "api-pagination",
  "lab-arch-styles": "arch-styles",
  "lab-arch-eip": "arch-eip",
  "lab-arch-saga": "arch-saga",
  "lab-arch-eventstorm": "arch-eventstorm",
  "lab-arch-arc42": "arch-arc42",
  "lab-proc-bpmn-collab": "proc-bpmn-collab",
  "lab-proc-uml-hold": "proc-uml-hold",
  "lab-proc-figma": "proc-figma",
  "lab-proc-mobile-web": "proc-mobile-web",
  "lab-nfr-slo": "nfr-slo",
  "lab-nfr-security": "nfr-security",
  "lab-nfr-pci": "nfr-pci-log",
  "lab-nfr-wcag": "nfr-wcag",
  "lab-bank-iso20022": "bank-iso20022",
  "lab-bank-cp-cnp": "bank-cp-cnp",
  "lab-bank-aml-kyc": "bank-aml-kyc",
  "lab-bank-ledger-deep": "bank-ledger-deep",
  "lab-bank-recon": "bank-recon",
  "lab-car-portfolio": "car-portfolio",
  "lab-car-cert": "car-cert",
  "lab-soft-raci": "soft-raci",
  "lab-soft-workshop": "soft-workshop",
};

export function topicForLab(labId) {
  return topicById(LAB_TOPIC[labId]);
}

export const TOPIC_LAB = Object.fromEntries(
  Object.entries(LAB_TOPIC).map(([labId, topicId]) => [topicId, labId]),
);

const TOPIC_LAB_ALIAS = {
  "req-gost-34602": "lab-req-gost",
  "nfr-pci-log": "lab-nfr-pci",
};

export function labIdForTopic(topicId) {
  if (!topicId) return "";
  if (TOPIC_LAB_ALIAS[topicId]) return TOPIC_LAB_ALIAS[topicId];
  const direct = "lab-" + topicId;
  if (LAB_TOPIC[direct]) return direct;
  return TOPIC_LAB[topicId] || direct;
}
