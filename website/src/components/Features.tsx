import {
  LightningIcon,
  StackIcon,
  ShieldCheckIcon,
  PathIcon,
} from "@phosphor-icons/react";
import type { Icon } from "@phosphor-icons/react";
import Reveal from "./Reveal";

type Feature = {
  icon: Icon;
  title: string;
  body: string;
  meta: string;
  className: string;
};

const features: Feature[] = [
  {
    icon: LightningIcon,
    title: "实时流式传输",
    body: "每一步操作、工具调用和 Token 消耗通过持久连接实时推送。无轮询、无延迟 — Agent 思考的同时，你就能看到。",
    meta: "WebSocket · < 5ms 开销",
    className: "md:col-span-2",
  },
  {
    icon: StackIcon,
    title: "多 Agent 支持",
    body: "在一个控制台监控数百个并发 Agent。按模型、Session 或状态筛选。",
    meta: "无限接入",
    className: "md:col-span-1",
  },
  {
    icon: ShieldCheckIcon,
    title: "Token 鉴权",
    body: "每个 Agent 和查看者都有独立的 Token 权限。支持轮换、撤销和审计，无需重新部署。",
    meta: "按 Agent 隔离",
    className: "md:col-span-1",
  },
  {
    icon: PathIcon,
    title: "轨迹追踪",
    body: "每次运行的完整可回放时间线。逐步检查输入、输出和决策 — Session 结束后仍可回溯。",
    meta: "持久化 · 可回放",
    className: "md:col-span-2",
  },
];

export default function Features() {
  return (
    <section id="features" className="mx-auto max-w-6xl scroll-mt-20 px-5 py-20 sm:px-8 sm:py-28">
      <Reveal>
        <p className="font-mono text-xs uppercase tracking-[0.2em] text-brand-400">
          核心能力
        </p>
        <h2 className="mt-3 max-w-xl text-3xl font-semibold tracking-tight text-white sm:text-4xl">
          为生产瘕的 Agent 团队打造
        </h2>
      </Reveal>

      <div className="mt-12 grid gap-4 md:grid-cols-3">
        {features.map((f, i) => (
          <Reveal key={f.title} delay={i * 0.06} className={f.className}>
            <article className="group relative h-full overflow-hidden rounded-xl border border-ink-800 bg-ink-900/40 p-6 transition-colors hover:border-ink-700 sm:p-7">
              <div
                aria-hidden
                className="pointer-events-none absolute -right-16 -top-16 h-40 w-40 rounded-full bg-brand-500/0 blur-3xl transition-colors duration-500 group-hover:bg-brand-500/10"
              />
              <div className="flex h-10 w-10 items-center justify-center rounded-lg border border-ink-800 bg-ink-950 text-brand-400">
                <f.icon size={20} weight="duotone" />
              </div>
              <h3 className="mt-5 text-lg font-medium text-white">{f.title}</h3>
              <p className="mt-2 text-[15px] leading-relaxed text-ink-400">
                {f.body}
              </p>
              <p className="mt-5 font-mono text-xs text-ink-500">{f.meta}</p>
            </article>
          </Reveal>
        ))}
      </div>
    </section>
  );
}
