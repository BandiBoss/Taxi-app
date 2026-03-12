import js from "@eslint/js";
import globals from "globals";
import pluginReact from "eslint-plugin-react";
import importPlugin from "eslint-plugin-import";

export default [
  js.configs.recommended,
  pluginReact.configs.flat.recommended,
  {
    files: ["**/*.{js,jsx,mjs,cjs}"],
    languageOptions: { globals: globals.browser },
    rules: {
      "react/prop-types": "off",
    },
    
  },
  {
    ignores: ["**/*.test.{js,jsx,ts,tsx}", "src/setupTests.js"],
  },
  {
    plugins: { import: importPlugin },
    rules: {
      "import/prefer-default-export": "off",
    },
  },
];
