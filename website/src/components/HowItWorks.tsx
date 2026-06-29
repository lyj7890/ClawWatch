import { DownloadSimpleIcon, PlugsConnectedIcon, ChartLineUpIcon } from "@phosphor-icons/react";
import type { Icon } from "@phosphor-icons/react";
import Reveal from "./Reveal";

type Step = {
  icon: Icon;
  n: string;
  title: string;
  body: string;
  code: string;
};

const steps: Step[] = [
  {
    icon: DownloadSimpleIcon,
    n: "01",
    title: "安装 Agent",
    body: "轻量级 Agent 二进制，零配置，一个命令启动。",
    code: "./clawwatch-agent --hub wss://hub.your.domain",
  },
  {
    icon: PlugsConnectedIcon,
    n: "02",
    title: "连接 Hub",
    body: "指定 Hub 地址和 Token，自动注册并开始流式上报。",
    code: "--hub wss://hub.acme.dev --token ***",
  },
  {
    icon: ChartLineUpIcon,
    n: "03",
    title: "查看面板",
    body: "打开浏览器，实时查看所有 Agent 的轨迹、Token 和工具调用。",
    code: "open https://hub.acme.dev",
  },
];

export default function HowItWorks() {
  return (
    <section id="how-it-works" className="mx-auto max-w-6xl scroll-mt-20 px-5 py-20 sm:px-8 sm:py-28">
      <Reveal>
        <p className="font-mono text-xs uppercase tracking-[0.2em] text-brand-400">
          快速开始
        </p>
        <h2 className="mt-3 max-w-xl text-3xl font-semibold tracking-tight text-white sm:text-4xl">
          三步即可上线
        </h2>
      </Reveal>

      <div className="mt-12 grid gap-4 md:grid-cols-3">
        {steps.map((s, i) => (
          <Reveal key={s.n} delay={i * 0.08}>
            <div className="relative h-full rounded-xl border border-ink-800 bg-ink-900/40 p-6 sm:p-7">
              <div className="flex items-center justify-between">
                <div className="flex h-10 w-10 items-center justify-center rounded-lg border border-ink-800 bg-ink-950 text-brand-400">
                  <s.icon size={20} weight="duotone" />
                </div>
                <span className="font-mono text-2xl font-medium text-ink-700">
                  {s.n}
                </span>
              </div>
              <h3 className="mt-5 text-lg font-medium text-white">{s.title}</h3>
              <p className="mt-2 text-[15px] leading-relaxed text-ink-400">
                {s.body}
              </p>
              <div className="mt-5 rounded-md border border-ink-800 bg-ink-950 px-3 py-2 font-mono text-xs text-brand-400/90">
                <span className="text-ink-600">$ </span>
                {s.code}
              </div>
            </div>
          </Reveal>
        ))}
      </div>
    </section>
  );
}
