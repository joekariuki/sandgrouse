import type { Metadata } from "next";
import Image from "next/image";
import Link from "next/link";
import { CopyButton } from "@/components/copy-button";

export const metadata: Metadata = {
  title: "Sandgrouse - Stop burning data bundles on AI tools",
  description:
    "Open source LLM traffic optimization proxy. Cuts Claude Code bandwidth in half and shows you where the rest goes. Built for developers on metered connections.",
};

export default function HomePage() {
  return (
    <main className="flex flex-col items-center px-6 py-16 gap-24">
      {/* Section 1: Hero */}
      <section className="text-center max-w-2xl">
        <h1 className="text-4xl sm:text-5xl font-bold tracking-tight mb-6">
          Stop burning data bundles on AI tools.
        </h1>
        <p className="text-fd-muted-foreground text-lg mb-8">
          Sandgrouse optimizes LLM API traffic so developers on metered
          connections get full AI power at a fraction of the data cost.
        </p>
        <div className="relative inline-block w-full max-w-md">
          <pre className="bg-fd-card border border-fd-border rounded-lg px-6 pr-10 py-3 text-sm font-mono text-left">
            npx sandgrouse
          </pre>
          <CopyButton text="npx sandgrouse" />
        </div>
        <p className="text-fd-muted-foreground text-sm mt-4">
          Works with Claude Code, Cursor, ChatGPT, and any OpenAI-compatible
          tool.
        </p>
      </section>

      {/* Section 2: The problem */}
      <section className="max-w-2xl">
        <h2 className="text-2xl font-bold mb-4">
          A 1.2GB data bundle shouldn&apos;t disappear after two Claude Code
          prompts.
        </h2>
        <p className="text-fd-muted-foreground leading-relaxed">
          LLM coding tools send your entire conversation history, every file
          they&apos;ve read, and an identical system prompt on every API request
          — and they send each request twice. A single prompt can trigger 50+
          requests in a chain, each one larger than the last as context
          accumulates. Requests are 99% of your bandwidth, none of it is
          compressed, and the APIs won&apos;t accept it compressed either.
        </p>
      </section>

      {/* Section 3: How it works */}
      <section className="max-w-2xl w-full">
        <h2 className="text-2xl font-bold mb-6 text-center">
          One proxy. Half the data. Full visibility.
        </h2>
        <pre className="bg-fd-card border border-fd-border rounded-lg p-6 text-sm font-mono text-fd-muted-foreground overflow-x-auto">
          {`Your AI tool (Claude Code, Cursor, etc.)
         |
         |  Sends every request twice (~150KB each, growing)
         v
   Sandgrouse proxy (localhost:8080)
         |
         |  Coalesces duplicates, compresses responses
         |  Shows you exactly where every byte goes
         v
   Cloud API (Anthropic, OpenAI)`}
        </pre>
        <p className="text-fd-muted-foreground text-sm text-center mt-4">
          Everything runs locally. No data is sent anywhere except the original
          API destination.
        </p>
      </section>

      {/* Section 4: Quick start */}
      <section className="max-w-2xl w-full">
        <h2 className="text-2xl font-bold mb-6 text-center">
          Get started in 30 seconds.
        </h2>
        <div className="flex flex-col gap-6">
          <div>
            <p className="text-sm font-medium mb-2 text-fd-muted-foreground">
              Step 1: Install
            </p>
            <div className="relative">
              <pre className="bg-fd-card border border-fd-border rounded-lg px-4 pr-12 py-3 text-sm font-mono">
                npm install -g sandgrouse
              </pre>
              <CopyButton text="npm install -g sandgrouse" />
            </div>
          </div>
          <div>
            <p className="text-sm font-medium mb-2 text-fd-muted-foreground">
              Step 2: Start the proxy
            </p>
            <div className="relative">
              <pre className="bg-fd-card border border-fd-border rounded-lg px-4 pr-12 py-3 text-sm font-mono">
                sg start
              </pre>
              <CopyButton text="sg start" />
            </div>
          </div>
          <div>
            <p className="text-sm font-medium mb-2 text-fd-muted-foreground">
              Step 3: Point your AI tools at it
            </p>
            <div className="relative">
              <pre className="bg-fd-card border border-fd-border rounded-lg px-4 pr-12 py-3 text-sm font-mono">
                export ANTHROPIC_BASE_URL=http://localhost:8080
              </pre>
              <CopyButton text="export ANTHROPIC_BASE_URL=http://localhost:8080" />
            </div>
          </div>
        </div>
        <p className="text-fd-muted-foreground text-sm text-center mt-6">
          Also available via Homebrew (
          <code className="text-fd-foreground">brew install sandgrouse</code>),
          direct download, or{" "}
          <code className="text-fd-foreground">npx sandgrouse</code>.
        </p>
      </section>

      {/* Dashboard preview */}
      <section className="max-w-3xl w-full">
        <div className="rounded-lg border border-fd-border overflow-hidden shadow-lg">
          <Image
            src="/dashboard.png"
            alt="Sandgrouse bandwidth dashboard showing real-time request tracking and compression savings"
            width={1200}
            height={675}
            className="w-full h-auto"
          />
        </div>
        <p className="text-fd-muted-foreground text-sm text-center mt-4">
          Real-time bandwidth dashboard at localhost:8585
        </p>
      </section>

      {/* Section 5: What you get */}
      <section className="max-w-3xl w-full">
        <div className="grid sm:grid-cols-3 gap-6">
          <div className="bg-fd-card border border-fd-border rounded-lg p-6">
            <div className="text-2xl mb-3">⚡</div>
            <h3 className="font-bold mb-2">Request coalescing</h3>
            <p className="text-fd-muted-foreground text-sm leading-relaxed">
              Claude Code sends every request twice. Sandgrouse catches the
              duplicate and drops it. ~50% bandwidth reduction, instantly.
            </p>
          </div>
          <div className="bg-fd-card border border-fd-border rounded-lg p-6">
            <div className="text-2xl mb-3">📊</div>
            <h3 className="font-bold mb-2">Bandwidth dashboard</h3>
            <p className="text-fd-muted-foreground text-sm leading-relaxed">
              See exactly how much data each session, each request, each tool
              consumes. Real-time at localhost:8585.
            </p>
          </div>
          <div className="bg-fd-card border border-fd-border rounded-lg p-6">
            <div className="text-2xl mb-3">🗜️</div>
            <h3 className="font-bold mb-2">Response compression</h3>
            <p className="text-fd-muted-foreground text-sm leading-relaxed">
              gzip/brotli on API responses. Plus: context deduplication and
              delta encoding coming in v0.2.
            </p>
          </div>
        </div>
      </section>

      {/* Section 6: Why "sandgrouse" */}
      <section className="max-w-2xl">
        <div className="bg-fd-card border border-fd-border rounded-lg p-8 text-center">
          <h2 className="text-xl font-bold mb-4">
            Why &ldquo;sandgrouse&rdquo;
          </h2>
          <p className="text-fd-muted-foreground leading-relaxed italic">
            The sandgrouse is an African bird whose breast feathers absorb water
            like a sponge. Every morning, it flies up to 30km across the desert
            to a waterhole, soaks its feathers, and flies back to its chicks.
            The most efficient water transport system in nature. Efficient data
            transport across bandwidth-scarce environments. That&apos;s what
            this does.
          </p>
        </div>
      </section>

      {/* Section 7: Footer */}
      <footer className="text-center text-sm text-fd-muted-foreground pb-8">
        <div className="flex flex-wrap justify-center gap-6 mb-4">
          <Link
            href="https://github.com/joekariuki/sandgrouse"
            className="hover:text-fd-foreground transition-colors"
          >
            GitHub — Star the repo
          </Link>
          <Link
            href="https://github.com/joekariuki/sandgrouse/blob/main/MANIFESTO.md"
            className="hover:text-fd-foreground transition-colors"
          >
            Manifesto — Read the full story
          </Link>
        </div>
        <p>
          MIT Licensed. Built from Nairobi by{" "}
          <Link
            href="https://github.com/joekariuki"
            className="text-fd-foreground hover:underline"
          >
            Joe Kariuki
          </Link>
          .
        </p>
      </footer>
    </main>
  );
}
