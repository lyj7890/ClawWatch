import Nav from "./components/Nav";
import Hero from "./components/Hero";
import Features from "./components/Features";
import Architecture from "./components/Architecture";
import HowItWorks from "./components/HowItWorks";
import CodeExample from "./components/CodeExample";
import Footer from "./components/Footer";

export default function App() {
  return (
    <div className="relative min-h-screen overflow-x-hidden bg-ink-950 text-ink-300">
      {/* ambient top glow */}
      <div
        aria-hidden
        className="pointer-events-none fixed inset-x-0 top-0 z-0 h-[420px] bg-[radial-gradient(60%_100%_at_50%_0%,rgba(16,185,129,0.10),transparent_70%)]"
      />
      <div aria-hidden className="bg-grain pointer-events-none fixed inset-0 z-0 opacity-60" />

      <div className="relative z-10">
        <Nav />
        <main>
          <Hero />
          <Features />
          <Architecture />
          <HowItWorks />
          <CodeExample />
        </main>
        <Footer />
      </div>
    </div>
  );
}
