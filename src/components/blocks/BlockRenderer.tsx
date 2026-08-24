import type { ContentBlock, ModuleContent } from "@/types/content";
import { TheoryView } from "./TheoryView";
import { CalloutView } from "./CalloutView";
import { ExampleView } from "./ExampleView";
import { InfographicView } from "./InfographicView";
import { DiagramView } from "./DiagramView";
import { AccordionView } from "./AccordionView";
import { CompareView } from "./CompareView";
import { StepsView } from "./StepsView";
import { QuizView } from "./QuizView";
import { PracticeView } from "./PracticeView";
import { ChecklistView } from "./ChecklistView";
import { TemplateView } from "./TemplateView";
import { CaseView } from "./CaseView";
import { SoftView } from "./SoftView";
import { SortView } from "./SortView";
import { MatchView } from "./MatchView";
import { SceneView } from "./SceneView";
import { OrderView } from "./OrderView";
import { SpotView } from "./SpotView";
import { BpmnView } from "./BpmnView";

export function BlockRenderer({
  block,
  module,
}: {
  block: ContentBlock;
  module: ModuleContent;
}) {
  switch (block.kind) {
    case "theory":
      return <TheoryView block={block} />;
    case "callout":
      return <CalloutView block={block} />;
    case "example":
      return <ExampleView block={block} />;
    case "infographic":
      return <InfographicView block={block} />;
    case "diagram":
      return <DiagramView block={block} />;
    case "accordion":
      return <AccordionView block={block} />;
    case "compare":
      return <CompareView block={block} />;
    case "steps":
      return <StepsView block={block} />;
    case "quiz":
      return <QuizView block={block} moduleId={module.id} />;
    case "practice":
      return <PracticeView block={block} moduleId={module.id} />;
    case "checklist":
      return <ChecklistView block={block} moduleId={module.id} />;
    case "template":
      return <TemplateView block={block} />;
    case "case":
      return <CaseView block={block} moduleId={module.id} />;
    case "soft":
      return <SoftView block={block} />;
    case "sort":
      return <SortView block={block} moduleId={module.id} />;
    case "match":
      return <MatchView block={block} moduleId={module.id} />;
    case "scene":
      return <SceneView block={block} moduleId={module.id} />;
    case "order":
      return <OrderView block={block} moduleId={module.id} />;
    case "spot":
      return <SpotView block={block} moduleId={module.id} />;
    case "bpmn":
      return <BpmnView block={block} />;
    default:
      return null;
  }
}
