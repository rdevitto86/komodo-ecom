/**
 * Generate TypeScript types from OpenAPI specs.
 * Run with: bun run gen:types
 *
 * Two outputs per service (when spec exists and has schemas):
 *
 *   src/lib/types/generated/<service>.ts
 *     Always overwritten. Contains raw openapi-typescript output (paths, operations,
 *     components) plus named type exports for every schema so consumers can import by name.
 *
 *   src/lib/types/<service>.ts
 *     Created ONCE (never overwritten). Re-exports everything from generated via
 *     `export type *` and leaves room below for local BFF types, error shapes, etc.
 *     Edit this file freely — the script will not touch it again.
 *
 * Services without a spec (empty file) are skipped.
 * Add entries here when new specs are written.
 */

import openapiTS, { astToString } from "openapi-typescript";
import { parse as parseYaml } from "yaml";
import { resolve, dirname } from "path";
import {
  writeFileSync,
  readFileSync,
  mkdirSync,
  statSync,
  existsSync,
} from "fs";
import { fileURLToPath } from "url";

const __dirname = dirname(fileURLToPath(import.meta.url));
const root = resolve(__dirname, "../..");
const generatedDir = resolve(__dirname, "../src/lib/types/generated");
const typesDir = resolve(__dirname, "../src/lib/types");

const services: { name: string; spec: string }[] = [
  { name: "address",        spec: "apis/komodo-address-api/openapi.yaml" },
  { name: "auth",           spec: "apis/komodo-auth-api/openapi.yaml" },
  { name: "cart",           spec: "apis/komodo-cart-api/openapi.yaml" },
  { name: "event-bus",      spec: "apis/komodo-event-bus-api/openapi.yaml" },
  { name: "order-returns",  spec: "apis/komodo-order-returns-api/openapi.yaml" },
  { name: "search",         spec: "apis/komodo-search-api/openapi.yaml" },
  { name: "shop-inventory", spec: "apis/komodo-shop-inventory-api/openapi.yaml" },
  { name: "shop-items",     spec: "apis/komodo-shop-items-api/openapi.yaml" },
  { name: "user",           spec: "apis/komodo-user-api/openapi.yaml" },
];

mkdirSync(generatedDir, { recursive: true });

let passed = 0;
let failed = 0;

for (const { name, spec } of services) {
  const input = resolve(root, spec);
  const generatedOutput = resolve(generatedDir, `${name}.ts`);
  const surfaceOutput = resolve(typesDir, `${name}.ts`);

  try {
    const size = statSync(input).size;
    if (size === 0) {
      console.log(`  - ${name} (spec is empty, skipped)`);
      continue;
    }

    // Parse spec for schema names
    const specText = readFileSync(input, "utf-8");
    const parsed = parseYaml(specText) as {
      components?: { schemas?: Record<string, unknown> };
    };
    const schemas = Object.keys(parsed?.components?.schemas ?? {});

    // Generate raw types + append named exports so generated file is self-contained
    const ast = await openapiTS(new URL(`file://${input}`));
    const namedExports = schemas.length > 0
      ? [
          ``,
          `// Named schema exports — import directly without referencing components['schemas']`,
          ...schemas.map((s) => `export type ${s} = components['schemas']['${s}'];`),
        ].join("\n")
      : "";

    const generatedContent = [
      `// Auto-generated from ${spec}`,
      `// Do not edit — run \`bun run gen:types\` to regenerate.`,
      ``,
      astToString(ast),
      namedExports,
    ].join("\n");
    writeFileSync(generatedOutput, generatedContent);

    // Surface file: create once, never overwrite — developer owns it
    if (!existsSync(surfaceOutput)) {
      const surfaceContent = [
        `// Types for ${name} — edit freely, this file is never overwritten by gen:types.`,
        `//`,
        `// Re-exports all generated schema types from the OpenAPI spec.`,
        `// Add local BFF types, error shapes, view models, etc. below.`,
        `export type * from './generated/${name}';`,
        ``,
        `// Local types (add below):`,
        ``,
      ].join("\n");
      writeFileSync(surfaceOutput, surfaceContent);
      console.log(`  ✓ ${name} (${schemas.length} types, surface file created)`);
    } else {
      console.log(`  ✓ ${name} (${schemas.length} types, surface file unchanged)`);
    }

    passed++;
  } catch (err) {
    console.error(`  ✗ ${name}: ${(err as Error).message}`);
    failed++;
  }
}

console.log(`\n${passed} generated, ${failed} failed.`);
if (failed > 0) process.exit(1);
