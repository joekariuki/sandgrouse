# AI tools shouldn't assume unlimited bandwidth.

Every major AI tool built today (Claude Code, Cursor, Codex, Copilot, Gemini) is designed with one invisible assumption: that the person using it has fast, cheap, unlimited internet.

That assumption is wrong for most of the world.

## The problem

As a software engineer currently based in Nairobi, Kenya, I use Claude Code to pair program, Claude to think through problems, and ChatGPT for research. These tools have transformed how I work.

Power outages are quite common in Nairobi. When the electricity cuts and the Wi-Fi goes down with it, I switch to my phone's 4G hotspot and keep working.

Here's what I've noticed: a 1.2GB data bundle disappears in one to two Claude Code requests. I end up buying 20 KES top-ups (about $0.15 each) every few minutes just to keep my session alive. In one recent coding session, I went through 8 of them. That's $1.20 in mobile data for a couple hours of AI-assisted coding, something developers on office Wi-Fi never think about.

This isn't a minor inconvenience. In Kenya, 1GB of mobile data costs real money. Across Africa, South Asia, Southeast Asia, and Latin America, where billions of people live and work, mobile data is metered, expensive, and often the most accessible means of having a reliable connection.

In Sub-Saharan Africa alone, 1GB of data costs 2.4% of average monthly income, and for the poorest 40%, it's closer to 5%. More than 30 African countries experience regular power outages, and the pattern repeats across South Asia and Latin America. When the power cuts, Wi-Fi dies, and mobile data becomes the only way to stay online.

The developers building AI tools don't experience this. They're in San Francisco, New York, London - cities where gigabit fiber is cheap and Wi-Fi is everywhere. They've never had to think about how many megabytes a single API call consumes, because for them, the answer doesn't matter.

Maybe sth like its not something they have to consider in their environment.

For the rest of us, it matters a lot.

## Why it happens

LLM-based tools are bandwidth hogs by design. Every time Claude Code reads a file, runs a command, or thinks about your code, it sends the entire conversation history - including every file it has read - back to the cloud as a single HTTP request. A mid-session request can easily be 200-300KB of JSON text.

A single prompt like "fix this bug" might trigger 10-20 of these requests in a chain: read the file, think, edit, read again, run the tests, read the output, edit again. Each one carries the growing context. Each one goes over your mobile connection.

And here's the thing that makes it worse: almost none of this traffic is compressed. JSON text compresses at 70-90% ratios, but most LLM API clients don't negotiate compression properly. You're sending raw text over the wire when a fraction of the bytes would carry the same information.

On top of that, the same content gets re-sent constantly. The system prompt is identical on every request. Files that haven't changed get re-transmitted in full. Conversation history that was sent 30 seconds ago gets sent again with each new message. There's enormous redundancy that nobody has bothered to eliminate because, again, bandwidth is free in San Francisco.

## What we're building

Sandgrouse is a local proxy that sits between your AI tools and the cloud. It intercepts LLM API traffic and optimizes it before it leaves your device:

- **Compression**: Enforces gzip/brotli on all traffic (70-80% reduction immediately)
- **Deduplication**: Identifies and eliminates repeated content like system prompts, unchanged files, and conversation history
- **Delta encoding**: For files that AI coding tools re-read, sends only what changed
- **Smart routing**: Routes simple tasks to a local model, saving the cloud API and your bandwidth for complex work
- **Measurement**: Shows you exactly how much bandwidth you're saving, request by request

The proxy is transparent. Your AI tools don't know it's there. You install it, start it, and your data consumption drops by 80-90%.

## Why "sandgrouse"

The sandgrouse is a bird found across Africa, Asia, and Southern Europe. The male sandgrouse has a remarkable adaptation: its breast feathers can absorb and hold water like a sponge. Every morning, it flies up to 30 kilometers to a waterhole, soaks its feathers, and flies back across the desert to its chicks, who drink the water from its plumage.

It's the most efficient water transport system in nature. A small body carrying maximum payload across harsh, resource-scarce terrain.

That's what this project does. It carries your AI data efficiently across bandwidth-scarce networks, delivering full capability with minimal waste.

## Who this is for

If you've ever:

- Watched a data bundle vanish during a Claude Code session
- Rationed your AI tool usage based on how much data you have left this week
- Switched between Wi-Fi and mobile hotspot and felt the cost difference
- Wished you could use AI tools freely without worrying about bandwidth
- Built a product whose users are in regions with metered internet

Then this is for you.

## The bigger picture

This project starts as a CLI tool for developers. But the mission is larger: **make AI usable everywhere in the world.**

Today, access to AI is quietly gated by bandwidth, not by the cost of the AI service itself (most have free tiers) but by the invisible cost of the internet connection needed to use them. A student in Lagos, a developer in Dhaka, an entrepreneur in Lima: they all have the same AI tools available as someone in New York. But the cost of using those tools is fundamentally different.

I'm building sandgrouse because nobody else is solving this. The companies building AI tools don't experience the problem. The infrastructure companies building networks are focused on speed, not efficiency. The gap between "AI is available" and "AI is usable" is a bandwidth gap, and it's widening as AI tools get more powerful and more data-hungry.

We're going to close it.

---

_If this resonates with you, star the repo, try the tool, and tell me what you think. If you're a developer who experiences this problem, I especially want to hear from you._
