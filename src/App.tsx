import { Navigate, Route, Routes, useLocation } from "react-router-dom";
import { AnimatePresence } from "framer-motion";
import { AppShell } from "@/components/layout/AppShell";
import { HomePage } from "@/pages/HomePage";
import { GradePage } from "@/pages/GradePage";
import { LevelPage } from "@/pages/LevelPage";
import { ModulePage } from "@/pages/ModulePage";
import { InterviewPage } from "@/pages/InterviewPage";
import { ExamPage } from "@/pages/ExamPage";
import { AchievementsPage } from "@/pages/AchievementsPage";
import { ProfilePage } from "@/pages/ProfilePage";
import { LabListPage, LabMissionPage } from "@/pages/LabPage";
import { QuestPage } from "@/pages/QuestPage";
import { PetPage } from "@/pages/PetPage";
import { HandbookPage } from "@/pages/HandbookPage";
import { BpmnPage } from "@/pages/BpmnPage";
import { useEffect, useState } from "react";
import { useProgress } from "@/stores/progressStore";
import { BootScreen } from "@/components/ui/BootScreen";

export default function App() {
  const location = useLocation();
  const evaluateBadges = useProgress((s) => s.evaluateBadges);
  const xp = useProgress((s) => s.xp);
  const modules = useProgress((s) => s.modules);
  const interviewSeen = useProgress((s) => s.interviewSeen);
  const examBest = useProgress((s) => s.examBest);
  const labs = useProgress((s) => s.labs);
  const quest = useProgress((s) => s.quest);
  const streak = useProgress((s) => s.streak);
  const combo = useProgress((s) => s.combo);
  const [ready, setReady] = useState(() => useProgress.persist.hasHydrated());

  useEffect(() => {
    const unsub = useProgress.persist.onFinishHydration(() => setReady(true));
    if (useProgress.persist.hasHydrated()) setReady(true);
    return unsub;
  }, []);

  useEffect(() => {
    if (ready) evaluateBadges();
  }, [evaluateBadges, xp, modules, interviewSeen, examBest, labs, quest, streak, combo, ready]);

  if (!ready) return <BootScreen />;

  return (
    <AppShell>
      <AnimatePresence mode="wait">
        <Routes location={location} key={location.pathname}>
          <Route path="/" element={<HomePage />} />
          <Route path="/grade/:gradeId" element={<GradePage />} />
          <Route path="/level/:levelId" element={<LevelPage />} />
          <Route path="/module/:moduleId" element={<ModulePage />} />
          <Route path="/lab" element={<LabListPage />} />
          <Route path="/lab/:labId" element={<LabMissionPage />} />
          <Route path="/quest" element={<QuestPage />} />
          <Route path="/quest/:questId" element={<QuestPage />} />
          <Route path="/pet" element={<PetPage />} />
          <Route path="/book" element={<HandbookPage />} />
          <Route path="/bpmn" element={<BpmnPage />} />
          <Route path="/interview" element={<InterviewPage />} />
          <Route path="/exam" element={<ExamPage />} />
          <Route path="/achievements" element={<AchievementsPage />} />
          <Route path="/profile" element={<ProfilePage />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </AnimatePresence>
    </AppShell>
  );
}
