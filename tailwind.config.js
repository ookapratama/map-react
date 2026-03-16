/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: "#4f46e5",
          hover: "#4338ca",
          light: "#eef2ff",
        },
        secondary: "#0ea5e9",
        accent: "#f59e0b",
        success: "#10b981",
        danger: "#ef4444",
        warning: "#f59e0b",
      },
      fontFamily: {
        sans: ["Inter", "sans-serif"],
        heading: ["Outfit", "sans-serif"],
      },
      borderRadius: {
        xl: "12px",
      },
    },
  },
  plugins: [],
};
