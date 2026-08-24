export type GradeId = "intern" | "junior" | "middle" | "senior";
export type TrackMode = "linear" | "free";
export type SkillKind = "hard" | "soft";
export type BlockKind =
  | "theory"
  | "callout"
  | "example"
  | "infographic"
  | "diagram"
  | "accordion"
  | "compare"
  | "steps"
  | "quiz"
  | "practice"
  | "checklist"
  | "template"
  | "case"
  | "soft"
  | "sort"
  | "match"
  | "scene"
  | "order"
  | "spot"
  | "bpmn";

export type Grade = {
  id: GradeId;
  title: string;
  roman: string;
  tagline: string;
  color: string;
  glow: string;
  levels: LevelMeta[];
};

export type LevelMeta = {
  id: string;
  gradeId: GradeId;
  rank: 1 | 2 | 3;
  title: string;
  subtitle: string;
  hours: string;
  moduleIds: string[];
};

export type ModuleMeta = {
  id: string;
  levelId: string;
  gradeId: GradeId;
  order: number;
  title: string;
  teaser: string;
  minutes: number;
  skill: SkillKind;
  xp: number;
  tags: string[];
};

export type TheoryBlock = {
  kind: "theory";
  title: string;
  lead?: string;
  paragraphs: string[];
  bullets?: string[];
};

export type CalloutBlock = {
  kind: "callout";
  tone: "gold" | "mint" | "violet" | "rose";
  title: string;
  text: string;
};

export type ExampleBlock = {
  kind: "example";
  title: string;
  context: string;
  bad?: string;
  good?: string;
  note?: string;
};

export type InfographicItem = {
  label: string;
  hint: string;
  detail: string;
};

export type InfographicBlock = {
  kind: "infographic";
  title: string;
  caption?: string;
  items: InfographicItem[];
};

export type DiagramNode = {
  id: string;
  label: string;
  sub?: string;
  detail: string;
};

export type DiagramBlock = {
  kind: "diagram";
  title: string;
  hint: string;
  nodes: DiagramNode[];
};

export type BpmnBlock = {
  kind: "bpmn";
  title: string;
  hint: string;
  src: string;
  storageKey?: string;
};

export type AccordionBlock = {
  kind: "accordion";
  title: string;
  items: { q: string; a: string }[];
};

export type CompareBlock = {
  kind: "compare";
  title: string;
  leftTitle: string;
  rightTitle: string;
  rows: { label: string; left: string; right: string }[];
};

export type StepsBlock = {
  kind: "steps";
  title: string;
  items: { n: string; title: string; text: string }[];
};

export type QuizQuestion = {
  id: string;
  q: string;
  options: string[];
  answer: number;
  why: string;
};

export type QuizBlock = {
  kind: "quiz";
  title: string;
  passScore: number;
  questions: QuizQuestion[];
};

export type PracticeBlock = {
  kind: "practice";
  title: string;
  brief: string;
  prompt: string;
  placeholder: string;
  minChars: number;
  rubric: string[];
  sample: string;
};

export type ChecklistBlock = {
  kind: "checklist";
  title: string;
  items: string[];
};

export type TemplateBlock = {
  kind: "template";
  title: string;
  description: string;
  body: string;
};

export type CaseBlock = {
  kind: "case";
  title: string;
  situation: string;
  question: string;
  options: { text: string; good: boolean; why: string }[];
  debrief: string;
};

export type SoftBlock = {
  kind: "soft";
  title: string;
  scene: string;
  doThis: string[];
  dont: string[];
  phrase: string;
};

export type SortBlock = {
  kind: "sort";
  title: string;
  prompt: string;
  buckets: { id: string; title: string }[];
  items: { id: string; text: string; bucket: string; why: string }[];
};

export type MatchBlock = {
  kind: "match";
  title: string;
  prompt: string;
  pairs: { left: string; right: string }[];
};

export type SceneStep = {
  from: string;
  line: string;
  options: { text: string; good: boolean; why: string }[];
};

export type SceneBlock = {
  kind: "scene";
  title: string;
  setting: string;
  steps: SceneStep[];
};

export type OrderBlock = {
  kind: "order";
  title: string;
  prompt: string;
  items: { id: string; text: string; pos: number }[];
};

export type SpotBlock = {
  kind: "spot";
  title: string;
  prompt: string;
  lines: { id: string; text: string; bad: boolean; why: string }[];
};

export type ContentBlock =
  | TheoryBlock
  | CalloutBlock
  | ExampleBlock
  | InfographicBlock
  | DiagramBlock
  | AccordionBlock
  | CompareBlock
  | StepsBlock
  | QuizBlock
  | PracticeBlock
  | ChecklistBlock
  | TemplateBlock
  | CaseBlock
  | SoftBlock
  | SortBlock
  | MatchBlock
  | SceneBlock
  | OrderBlock
  | SpotBlock
  | BpmnBlock;

export const PLAY_KINDS: BlockKind[] = [
  "sort",
  "match",
  "scene",
  "order",
  "spot",
  "case",
  "quiz",
  "practice",
  "bpmn",
];

export function isPlayBlock(kind: BlockKind) {
  return PLAY_KINDS.includes(kind);
}

export type LabMission = {
  id: string;
  gradeId: GradeId;
  title: string;
  teaser: string;
  minutes: number;
  xp: number;
  setting: string;
  blocks: ContentBlock[];
};

export type ModuleContent = ModuleMeta & {
  goals: string[];
  blocks: ContentBlock[];
};

export type InterviewItem = {
  id: string;
  topic: string;
  question: string;
  short: string;
  strong: string;
  traps: string[];
};

export type ExamQuestion = {
  id: string;
  scene: string;
  q: string;
  options: string[];
  answer: number;
  why: string;
  weight: number;
};

export type Achievement = {
  id: string;
  title: string;
  description: string;
  icon: string;
  xp: number;
};

export type PracticeRecord = {
  moduleId: string;
  blockTitle: string;
  text: string;
  updatedAt: number;
};
