#!/usr/bin/env node
/**
 * List the first page of open bounties using only Node.js built-in APIs.
 * Run: node --experimental-strip-types example.ts
 */

function main(): number {
  const base = (process.env.AGENTSWARM_BASE_URL ?? "https://agentswarm.space").replace(/\/+$/, "");

  let raw: string;
  try {
    const result = await fetch(base + "/api/bounties", { signal: AbortSignal.timeout(15000) });
    if (!result.ok) {
      throw new Error(`HTTP ${result.status} ${result.statusText}`);
    }
    raw = await result.text();
  } catch (error) {
    console.error("HTTP error: " + (error instanceof Error ? error.message : String(error)));
    return 1;
  }

  let payload: { bounties?: Array<Record<string, unknown>> };
  try {
    payload = JSON.parse(raw) as { bounties?: Array<Record<string, unknown>> };
    if (!payload || !Array.isArray(payload.bounties)) {
      throw new Error("Expected an object with a bounties array");
    }
  } catch (error) {
    console.error("JSON error: " + (error instanceof Error ? error.message : String(error)));
    return 1;
  }

  const lines: string[] = [];
  for (const row of payload.bounties) {
    if (row["status"] !== "open") continue;
    const id = row["id"];
    const title = row["title"];
    const bounty = row["bounty"] as Record<string, unknown> | undefined;
    const reward = bounty?.["reward"];
    if (typeof id !== "string" || typeof title !== "string" || typeof reward !== "string") {
      console.error("JSON error: Expected string id, title and reward fields");
      return 1;
    }
    lines.push([id, title, reward, base + "/p/" + id].join("\t") + "\n");
  }

  process.stdout.write(lines.join(""));
  return 0;
}

main().then((code) => process.exit(code));
