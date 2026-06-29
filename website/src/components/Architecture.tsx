import { RobotIcon, BroadcastIcon, MonitorIcon, ArrowRightIcon } from "@phosphor-icons/react";
import type { Icon } from "@phosphor-icons/react";
import Reveal from "./Reveal";

type Node = {
  icon: Icon;
  label: string;
  sub: string;
  detail: string;
};

const nodes: Node[] = [
  { icon: RobotIcon, label: "Agent", sub: "你的运行时", detail: "clawwatch-agent" },
  { icon: BroadcastIcon, label: "Hub", sub: "接入 + 分发", detail: ":4848 / ws" },
  { icon: MonitorIcon, label: "Console", sub: "实时面板", detail: "浏览器" },
];

export default function Architecture() {
  return (
    <section id="architecture" className="scroll-mt-20 border-y border-ink-800/60 bg-ink-925/40">
      <div className="mx-auto max-w-6xl px-5 py-20 sm:px-8 sm:py-28">
        <Reveal>
          <p className="font-mono text-xs uppercase tracking-[0.2em] text-brand-400">
            架构
          </p>
          <h2 className="mt-3 max-w-xl text-3xl font-semibold tracking-tight text-white sm:text-4xl">
            一个 Hub，所有 Agent，实时流入你的控制台。
          </h2>
        </Reveal>

        <Reveal delay={0.08} className="mt-12">
          <div className="overflow-hidden rounded-xl border border-ink-800 bg-ink-950">
            <div className="flex items-center gap-2 border-b border-ink-800 bg-ink-900/50 px-4 py-3 font-mono text-xs text-ink-500">
              <span className="h-2.5 w-2.5 rounded-full bg-ink-700" />
              data flow
            </div>

            <div className="grid items-stretch gap-4 p-6 sm:p-10 md:grid-cols-[1fr_auto_1fr_auto_1fr]">
              {nodes.map((n, i) => (
                <div key={n.label} className="contents">
                  <div className="flex flex-col items-center rounded-lg border border-ink-800 bg-ink-900/40 px-6 py-7 text-center">
                    <div className="flex h-12 w-12 items-center justify-center rounded-lg border border-ink-800 bg-ink-950 text-brand-400">
                      <n.icon size={24} weight="duotone" />
                    </div>
                    <p className="mt-4 text-base font-medium text-white">{n.label}</p>
                    <p className="mt-1 text-sm text-ink-500">{n.sub}</p>
                    <p className="mt-3 font-mono text-xs text-brand-400/80">{n.detail}</p>
                  </div>

                  {i < nodes.length - 1 && (
                    <div className="flex items-center justify-center text-ink-600">
                      <ArrowRightIcon
                        size={22}
                        weight="bold"
                        className="hidden rotate-0 md:block"
                      />
                      <ArrowRightIcon
                        size={22}
                        weight="bold"
                        className="rotate-90 md:hidden"
                      />
                    </div>
                  )}
                </div>
              ))}
            </div>

            <div className="border-t border-ink-800 bg-ink-900/30 px-6 py-4 font-mono text-xs text-ink-500 sm:px-10">
              <span className="text-brand-400">→</span> 事件为只追加模式，按主机隔离存储，支持 Session 历史回溯。
            </div>
          </div>
        </Reveal>
      </div>
    </section>
  );
}
