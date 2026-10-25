const fs = require("fs");
const vm = require("vm");
const code = fs.readFileSync(__dirname + "/source.js", "utf8");
const ctx = { window: {}, document: { cookie: "" } };
vm.createContext(ctx);
vm.runInContext(code, ctx);
const { parseSource, clipboardValue } = ctx.window.LdySource;

const cases = [
  ["", ""],
  ["  DY01  ", "DY01"],
  ["channel://DY01", "DY01"],
  ["share://abc", ""],
  ["official", ""],
  ["invite", ""],
  ["a", ""],
];
for (const [input, want] of cases) {
  const got = parseSource(input);
  if (got !== want) {
    console.error("parseSource(%j)=%j want %j", input, got, want);
    process.exit(1);
  }
}
if (clipboardValue("DY01") !== "channel://DY01") {
  console.error("clipboard mismatch");
  process.exit(1);
}
console.log("ok");
