import esbuild from "esbuild";
import esbuildPluginTsc from "esbuild-plugin-tsc";
import { NodeResolvePlugin as NodeResolve } from "@esbuild-plugins/node-resolve";

/**
 * @param {import "esbuild".BuildOptions} options
 * @returns {import "esbuild".BuildOptions}
 */
function createBuildSettings(options) {
	return {
		entryPoints: ["api/index.ts"],
		outfile: "api/dist/index.js",
		plugins: [
			esbuildPluginTsc({
				force: true,
				tsconfigPath: "./tsconfig.json",
			}),
			NodeResolve({
				extensions: [".ts", ".js"],
				onResolved: (resolved) => {
					if (resolved.includes("node_modules")) {
						return {
							external: true,
						};
					}
					return resolved;
				},
			}),
		],
		...options,
	};
}

const options = createBuildSettings({
	bundle: true,
	minify: true,
	minifySyntax: true,
	format: "esm",
	treeShaking: true,
	platform: "node",
	sourcemap: true,
});

esbuild.build(options);
