import { useState } from "react";
import { CopyIcon, CheckIcon, TerminalWindowIcon } from "@phosphor-icons/react";
import Reveal from "./Reveal";

type Line = { tokens: { t: string; c?: string }[] };

const CODE: Line[] = [
  { tokens: [{ t: "帮我接入 ClawWatch 监控：", c: "text-ink-300" }] },
  { tokens: [] },
  { tokens: [{ t: "1. 安装 Agent", c: "text-ink-400" }] },
  { tokens: [{ t: "   curl -fsSL https://raw.githubusercontent.com/", c: "text-brand-400" }] },
  { tokens: [{ t: "   lyj7890/ClawWatch/main/install.sh | sh", c: "text-brand-400" }] },
  { tokens: [] },
  { tokens: [{ t: "2. 启动连接", c: "text-ink-400" }] },
  { tokens: [{ t: "   clawwatch-agent --hub wss://clawatch.intra.mlamp.cn", c: "text-brand-400" }] },
  { tokens: [] },
  { tokens: [{ t: "3. 告诉我日志里输出的 Console Token", c: "text-ink-400" }] },
];

const PLAIN = `帮我接入 ClawWatch 监控：

1. 安装 Agent
   curl -fsSL https://raw.githubusercontent.com/lyj7890/ClawWatch/main/install.sh | sh

2. 启动连接
   clawwatch-agent --hub wss://clawatch.intra.mlamp.cn

3. 告诉我日志里输出的 Console Token`;

export default function CodeExample() {
  const [copied, setCopied] = useState(false);

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(PLAIN);
      setCopied(true);
      setTimeout(() => setCopied(false), 1600);
    } catch {
      /* clipboard unavailable */
    }
  };

  return (
    <section className="border-t border-ink-800/60 bg-ink-925/40">
      <div className="mx-auto grid max-w-6xl items-center gap-10 px-5 py-20 sm:px-8 sm:py-28 lg:grid-cols-[0.85fr_1.15fr]">
        <Reveal>
          <p className="font-mono text-xs uppercase tracking-[0.2em] text-brand-400">
            即刻接入
          </p>
          <h2 className="mt-3 text-3xl font-semibold tracking-tight text-white sm:text-4xl">
            复制这段，发给你的龙虾
          </h2>
          <p className="mt-4 max-w-md text-[15px] leading-relaxed text-ink-400">
            一次复制，龙虾自动完成安装、连接和配置。你只需要打开 Console 查看。
          </p>
        </Reveal>

        <Reveal delay={0.08}>
          <div className="overflow-hidden rounded-xl border border-ink-800 bg-ink-950 shadow-xl shadow-black/30">
            <div className="flex items-center gap-2 border-b border-ink-800 bg-ink-900/50 px-4 py-3">
              <TerminalWindowIcon size={15} className="text-ink-500" />
              <span className="font-mono text-xs text-ink-500">发给龙虾的指令</span>
              <button
                onClick={copy}
                className="ml-auto flex items-center gap-1.5 rounded-md border border-ink-800 px-2.5 py-1 font-mono text-xs text-ink-400 transition-colors hover:border-ink-700 hover:text-ink-100"
                aria-label="Copy code"
              >
                {copied ? (
                  <>
                    <CheckIcon size={13} weight="bold" className="text-brand-400" />
                    copied
                  </>
                ) : (
                  <>
                    <CopyIcon size={13} />
                    copy
                  </>
                )}
              </button>
            </div>

            <pre className="term-scroll overflow-x-auto px-5 py-5 font-mono text-[13px] leading-relaxed">
              <code className="text-ink-300">
                {CODE.map((line, i) => (
                  <div key={i} className="flex">
                    <span className="mr-4 w-5 shrink-0 select-none text-right text-ink-700">
                      {i + 1}
                    </span>
                    <span className="whitespace-pre">
                      {line.tokens.length === 0
                        ? "\u00A0"
                        : line.tokens.map((tok, j) => (
                            <span key={j} className={tok.c}>
                              {tok.t}
                            </span>
                          ))}
                    </span>
                  </div>
                ))}
              </code>
            </pre>
          </div>
        </Reveal>
      </div>
    </section>
  );
}
