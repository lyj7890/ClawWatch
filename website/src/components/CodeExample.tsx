import { useState } from "react";
import { CopyIcon, CheckIcon, TerminalWindowIcon } from "@phosphor-icons/react";
import Reveal from "./Reveal";

type Line = { tokens: { t: string; c?: string }[] };

const CODE: Line[] = [
  { tokens: [{ t: "# 安装 ClawWatch Agent", c: "text-ink-600" }] },
  { tokens: [{ t: "curl -fsSL", c: "text-fuchsia-300/90" }, { t: " https://raw.githubusercontent.com/" }] },
  { tokens: [{ t: "  lyj7890/ClawWatch/main/install.sh", c: "text-brand-400" }, { t: " | ", c: "text-ink-400" }, { t: "sh", c: "text-fuchsia-300/90" }] },
  { tokens: [] },
  { tokens: [{ t: "# 启动并连接到 Hub", c: "text-ink-600" }] },
  { tokens: [{ t: "clawwatch-agent", c: "text-sky-300" }, { t: " \\" }] },
  { tokens: [{ t: "  --hub ", c: "text-ink-400" }, { t: "wss://clawatch.intra.mlamp.cn", c: "text-brand-400" }] },
  { tokens: [] },
  { tokens: [{ t: "# ✅ Agent 已连接，打开浏览器查看", c: "text-ink-600" }] },
  { tokens: [{ t: "open", c: "text-fuchsia-300/90" }, { t: " https://clawatch.intra.mlamp.cn/console", c: "text-brand-400" }] },
];

const PLAIN = `# 安装 ClawWatch Agent
curl -fsSL https://raw.githubusercontent.com/lyj7890/ClawWatch/main/install.sh | sh

# 启动并连接到 Hub
clawwatch-agent --hub wss://clawatch.intra.mlamp.cn

# ✅ Agent 已连接，打开浏览器查看
open https://clawatch.intra.mlamp.cn/console`;

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
            代码示例
          </p>
          <h2 className="mt-3 text-3xl font-semibold tracking-tight text-white sm:text-4xl">
            几行代码，即刻接入。
          </h2>
          <p className="mt-4 max-w-md text-[15px] leading-relaxed text-ink-400">
            Agent 启动时指定 Hub 地址和 Token 即可开始上报。数据异步推送，
            不影响 Agent 主流程性能。
          </p>
        </Reveal>

        <Reveal delay={0.08}>
          <div className="overflow-hidden rounded-xl border border-ink-800 bg-ink-950 shadow-xl shadow-black/30">
            <div className="flex items-center gap-2 border-b border-ink-800 bg-ink-900/50 px-4 py-3">
              <TerminalWindowIcon size={15} className="text-ink-500" />
              <span className="font-mono text-xs text-ink-500">terminal</span>
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
