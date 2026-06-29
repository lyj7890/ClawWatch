import { GithubLogoIcon, BookOpenTextIcon, FileTextIcon } from "@phosphor-icons/react";
import Logo from "./Logo";

const cols = [
  {
    title: "产品",
    links: [
      { label: "功能", href: "#features" },
      { label: "架构", href: "#architecture" },
      { label: "接入指南", href: "#how-it-works" },
    ],
  },
  {
    title: "资源",
    links: [
      { label: "文档", href: "#docs" },
      { label: "API 参考", href: "#docs" },
      { label: "更新日志", href: "#docs" },
    ],
  },
];

export default function Footer() {
  return (
    <footer id="docs" className="scroll-mt-20 border-t border-ink-800/60 bg-ink-950">
      <div className="mx-auto max-w-6xl px-5 py-14 sm:px-8">
        <div className="grid gap-10 sm:grid-cols-2 lg:grid-cols-4">
          <div className="lg:col-span-2">
            <a href="#top" className="flex items-center gap-2.5">
              <Logo className="h-7 w-7" />
              <span className="font-mono text-[15px] font-medium tracking-tight text-white">
                claw<span className="text-brand-400">watch</span>
              </span>
            </a>
            <p className="mt-4 max-w-xs text-sm leading-relaxed text-ink-500">
              实时 AI Agent 监控平台，为工程团队而生。可自托管、
              开源、MIT 许可。
            </p>
            <div className="mt-5 flex items-center gap-3">
              <a
                href="https://github.com/lyj7890/ClawWatch"
                target="_blank"
                rel="noreferrer"
                className="flex items-center gap-2 rounded-lg border border-ink-800 px-3 py-2 text-sm text-ink-300 transition-colors hover:border-ink-700 hover:text-white"
              >
                <GithubLogoIcon size={16} weight="fill" />
                GitHub
              </a>
              <a
                href="#docs"
                className="flex items-center gap-2 rounded-lg border border-ink-800 px-3 py-2 text-sm text-ink-300 transition-colors hover:border-ink-700 hover:text-white"
              >
                <BookOpenTextIcon size={16} weight="duotone" />
                文档
              </a>
            </div>
          </div>

          {cols.map((col) => (
            <div key={col.title}>
              <p className="font-mono text-xs uppercase tracking-[0.18em] text-ink-500">
                {col.title}
              </p>
              <ul className="mt-4 space-y-3">
                {col.links.map((l) => (
                  <li key={l.label}>
                    <a
                      href={l.href}
                      className="text-sm text-ink-400 transition-colors hover:text-ink-100"
                    >
                      {l.label}
                    </a>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        <div className="mt-12 flex flex-col items-start justify-between gap-4 border-t border-ink-800/60 pt-6 sm:flex-row sm:items-center">
          <p className="font-mono text-xs text-ink-600">
            © {new Date().getFullYear()} ClawWatch
          </p>
          <a
            href="https://opensource.org/license/mit"
            target="_blank"
            rel="noreferrer"
            className="flex items-center gap-1.5 font-mono text-xs text-ink-600 transition-colors hover:text-ink-400"
          >
            <FileTextIcon size={13} />
            MIT License
          </a>
        </div>
      </div>
    </footer>
  );
}
