import { GithubLogo } from "@phosphor-icons/react";
import Logo from "./Logo";

const links = [
  { label: "功能", href: "#features" },
  { label: "架构", href: "#architecture" },
  { label: "接入", href: "#how-it-works" },
  { label: "文档", href: "#docs" },
];

export default function Nav() {
  return (
    <header className="sticky top-0 z-50 border-b border-ink-800/60 bg-ink-950/70 backdrop-blur-xl">
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-5 sm:px-8">
        <a href="#top" className="flex items-center gap-2.5">
          <Logo className="h-7 w-7" />
          <span className="font-mono text-[15px] font-medium tracking-tight text-ink-100 text-white">
            claw<span className="text-brand-400">watch</span>
          </span>
        </a>

        <nav className="hidden items-center gap-7 md:flex">
          {links.map((l) => (
            <a
              key={l.label}
              href={l.href}
              className="text-sm text-ink-400 transition-colors hover:text-ink-100"
            >
              {l.label}
            </a>
          ))}
        </nav>

        <div className="flex items-center gap-3">
          <a
            href="https://github.com/lyj7890/ClawWatch"
            target="_blank"
            rel="noreferrer"
            className="flex h-9 w-9 items-center justify-center rounded-lg border border-ink-800 text-ink-400 transition-colors hover:border-ink-700 hover:text-ink-100"
            aria-label="GitHub repository"
          >
            <GithubLogo size={18} weight="fill" />
          </a>
          <a
            href="#how-it-works"
            className="hidden rounded-lg bg-brand-400 px-3.5 py-2 text-sm font-medium text-ink-950 transition-colors hover:bg-brand-300 sm:inline-block"
          >
            快速开始
          </a>
        </div>
      </div>
    </header>
  );
}
