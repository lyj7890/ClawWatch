import { useEffect, useRef, useState } from "react";
import { motion, useReducedMotion } from "motion/react";
import { ArrowRight, GithubLogo, PulseIcon } from "@phosphor-icons/react";

type LogLine = {
  t: string;
  agent: string;
  msg: string;
  tone: "info" | "ok" | "tool" | "muted";
};

const FEED: LogLine[] = [
  { t: "12:04:01", agent: "agent-07", msg: "session.start  model=claude-opus", tone: "info" },
  { t: "12:04:01", agent: "agent-07", msg: "tool.call  read(src/index.ts)", tone: "tool" },
  { t: "12:04:02", agent: "agent-12", msg: "session.start  model=gpt-4.1", tone: "info" },
  { t: "12:04:03", agent: "agent-07", msg: "tool.result  482 tokens", tone: "muted" },
  { t: "12:04:04", agent: "agent-07", msg: "step.complete  ✓ patch applied", tone: "ok" },
  { t: "12:04:05", agent: "agent-12", msg: "tool.call  exec(npm test)", tone: "tool" },
  { t: "12:04:07", agent: "agent-12", msg: "step.complete  ✓ 38 passed", tone: "ok" },
  { t: "12:04:08", agent: "agent-07", msg: "session.end  4.2k tokens · 7 steps", tone: "muted" },
];

const toneColor: Record<LogLine["tone"], string> = {
  info: "text-ink-300",
  ok: "text-brand-400",
  tool: "text-sky-300/90",
  muted: "text-ink-500",
};

function useFeed(reduced: boolean) {
  const [count, setCount] = useState(reduced ? FEED.length : 0);
  useEffect(() => {
    if (reduced) return;
    if (count >= FEED.length) {
      const reset = setTimeout(() => setCount(0), 2600);
      return () => clearTimeout(reset);
    }
    const id = setTimeout(() => setCount((c) => c + 1), 620);
    return () => clearTimeout(id);
  }, [count, reduced]);
  return FEED.slice(0, count);
}

export default function Hero() {
  const reduced = useReducedMotion() ?? false;
  const visible = useFeed(reduced);
  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const el = scrollRef.current;
    if (el) el.scrollTop = el.scrollHeight;
  }, [visible.length]);

  return (
    <section id="top" className="mx-auto max-w-6xl px-5 pb-20 pt-16 sm:px-8 sm:pt-24">
      <div className="grid items-center gap-12 lg:grid-cols-[1.05fr_1fr] lg:gap-10">
        {/* Left — copy, left-aligned */}
        <div>
          <motion.div
            initial={reduced ? false : { opacity: 0, y: 12 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
            className="mb-6 inline-flex items-center gap-2 rounded-full border border-ink-800 bg-ink-900/60 px-3 py-1 text-xs text-ink-400"
          >
            <PulseIcon size={14} weight="bold" className="text-brand-400" />
            Live agent observability
          </motion.div>

          <motion.h1
            initial={reduced ? false : { opacity: 0, y: 16 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.55, delay: 0.05 }}
            className="text-balance text-5xl font-semibold leading-[1.02] tracking-tight text-white sm:text-6xl lg:text-7xl"
          >
            Claw<span className="text-brand-400">Watch</span>
          </motion.h1>

          <motion.p
            initial={reduced ? false : { opacity: 0, y: 16 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.55, delay: 0.12 }}
            className="mt-5 max-w-md text-lg leading-relaxed text-ink-400 sm:text-xl"
          >
            实时 AI Agent 监控平台，为工程团队而生。流式传输每一步
            操作轨迹、工具调用和 Token 消耗 — 跨所有 Agent，实时呈现。
          </motion.p>

          <motion.div
            initial={reduced ? false : { opacity: 0, y: 16 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.55, delay: 0.19 }}
            className="mt-8 flex flex-wrap items-center gap-3"
          >
            <a
              href="#how-it-works"
              className="group inline-flex items-center gap-2 rounded-lg bg-brand-400 px-5 py-3 text-sm font-medium text-ink-950 transition-colors hover:bg-brand-300"
            >
              快速开始
              <ArrowRight
                size={16}
                weight="bold"
                className="transition-transform group-hover:translate-x-0.5"
              />
            </a>
            <a
              href="https://github.com/lyj7890/ClawWatch"
              target="_blank"
              rel="noreferrer"
              className="inline-flex items-center gap-2 rounded-lg border border-ink-800 px-5 py-3 text-sm font-medium text-ink-200 transition-colors hover:border-ink-700 hover:text-white"
            >
              <GithubLogo size={16} weight="fill" />
              GitHub
            </a>
          </motion.div>

          <motion.div
            initial={reduced ? false : { opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.6, delay: 0.28 }}
            className="mt-8 flex items-center gap-5 font-mono text-xs text-ink-500"
          >
            <span>MIT 开源</span>
            <span className="h-1 w-1 rounded-full bg-ink-700" />
            <span>可自托管</span>
            <span className="h-1 w-1 rounded-full bg-ink-700" />
            <span>&lt; 5ms 开销</span>
          </motion.div>
        </div>

        {/* Right — terminal stream */}
        <motion.div
          initial={reduced ? false : { opacity: 0, y: 24, scale: 0.98 }}
          animate={{ opacity: 1, y: 0, scale: 1 }}
          transition={{ duration: 0.6, delay: 0.18 }}
          className="relative"
        >
          <div
            aria-hidden
            className="absolute -inset-4 -z-10 rounded-3xl bg-[radial-gradient(50%_50%_at_70%_30%,rgba(16,185,129,0.14),transparent_70%)]"
          />
          <div className="overflow-hidden rounded-xl border border-ink-800 bg-ink-925 shadow-2xl shadow-black/40">
            {/* title bar */}
            <div className="flex items-center gap-2 border-b border-ink-800 bg-ink-900/60 px-4 py-3">
              <span className="h-3 w-3 rounded-full bg-ink-700" />
              <span className="h-3 w-3 rounded-full bg-ink-700" />
              <span className="h-3 w-3 rounded-full bg-ink-700" />
              <span className="ml-2 font-mono text-xs text-ink-500">
                clawwatch · live feed
              </span>
              <span className="ml-auto flex items-center gap-1.5 font-mono text-xs text-brand-400">
                <span className="relative flex h-2 w-2">
                  <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-brand-400/70" />
                  <span className="relative inline-flex h-2 w-2 rounded-full bg-brand-400" />
                </span>
                LIVE
              </span>
            </div>

            {/* feed */}
            <div
              ref={scrollRef}
              className="term-scroll h-[300px] space-y-1 overflow-y-auto px-4 py-3 font-mono text-[12.5px] leading-relaxed sm:h-[340px]"
            >
              {visible.map((line, i) => (
                <motion.div
                  key={`${i}-${line.t}-${line.msg}`}
                  initial={reduced ? false : { opacity: 0, x: -6 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ duration: 0.25 }}
                  className="flex gap-2.5"
                >
                  <span className="shrink-0 text-ink-600">{line.t}</span>
                  <span className="shrink-0 text-ink-300">{line.agent}</span>
                  <span className={toneColor[line.tone]}>{line.msg}</span>
                </motion.div>
              ))}
              <div className="flex gap-2.5 text-ink-500">
                <span className="text-brand-400">›</span>
                <span className="cursor-blink inline-block h-4 w-2 bg-ink-500" />
              </div>
            </div>
          </div>
        </motion.div>
      </div>
    </section>
  );
}
