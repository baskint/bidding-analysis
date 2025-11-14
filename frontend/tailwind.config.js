/** @type {import('tailwindcss').Config} */
module.exports = {
  // CRITICAL: Change darkMode from 'media' (default) to 'class'
  darkMode: 'class', 
  content: [
    "./src/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
